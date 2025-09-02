package monitor

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/badgerv/monitoring-api/internal/storage"
	"github.com/jackc/pgx/v5"
)

type PostgresRepository struct {
	db *storage.DB
}

func NewPostgresRepository(db *storage.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetAllEndpoints(ctx context.Context) ([]Endpoint, error) {
	query := `SELECT id, service_name, url, server_name, api_method 
	          FROM endpoints`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var endpoints []Endpoint
	for rows.Next() {
		ep := Endpoint{}
		err := rows.Scan(&ep.ID, &ep.ServiceName, &ep.URL, &ep.ServerName, &ep.APIMethod)
		if err != nil {
			return nil, err
		}
		endpoints = append(endpoints, ep)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return endpoints, nil
}

// func (r *PostgresRepository) CreateEndpoint(ctx context.Context) ([]Endpoint, error) {
// 	query := `SELECT id, service_name, url, server_name, api_method
// 	          FROM endpoints`

// 	rows, err := r.db.Pool.Query(ctx, query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var endpoints []Endpoint
// 	for rows.Next() {
// 		ep := Endpoint{}
// 		err := rows.Scan(&ep.ID, &ep.ServiceName, &ep.URL, &ep.ServerName, &ep.APIMethod)
// 		if err != nil {
// 			return nil, err
// 		}
// 		endpoints = append(endpoints, ep)
// 	}

// 	if err = rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return endpoints, nil
// }

// Inserts a new endpoint into the database and returns the created record
func (r *PostgresRepository) CreateEndpoint(ctx context.Context, ep *Endpoint) (*Endpoint, error) {
    // Normalize URL (same rules as JSON loader)
    normalizeURL := func(raw string) string {
        raw = strings.TrimSpace(raw)
        if strings.HasSuffix(raw, "/") {
            raw = strings.TrimSuffix(raw, "/")
        }
        return strings.ToLower(raw)
    }
    ep.URL = normalizeURL(ep.URL)

    // First check if endpoint already exists
    var existingID int
    checkQuery := `SELECT id FROM endpoints WHERE url = $1`
    err := r.db.Pool.QueryRow(ctx, checkQuery, ep.URL).Scan(&existingID)
    if err == nil {
        // Found an existing endpoint → return error
        return nil, fmt.Errorf("endpoint with URL '%s' already exists (id=%d)", ep.URL, existingID)
    }
    if err != nil && err != pgx.ErrNoRows {
        // Real DB error
        return nil, fmt.Errorf("failed checking existing endpoint: %w", err)
    }

    // Start a transaction to ensure atomicity
    tx, err := r.db.Pool.Begin(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to start transaction: %w", err)
    }
    defer tx.Rollback(ctx) // Rollback if not committed

    // Insert into endpoints table
    insertEndpointQuery := `
        INSERT INTO endpoints (service_name, url, server_name, api_method, expected_status_code)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, service_name, url, server_name, api_method, expected_status_code
    `
    row := tx.QueryRow(ctx, insertEndpointQuery, ep.ServiceName, ep.URL, ep.ServerName, ep.APIMethod, ep.ExpectedCode)
    newEp := &Endpoint{}
    err = row.Scan(&newEp.ID, &newEp.ServiceName, &newEp.URL, &newEp.ServerName, &newEp.APIMethod, &newEp.ExpectedCode)
    if err != nil {
        return nil, fmt.Errorf("failed to insert endpoint: %w", err)
    }

    // Check if any endpoint_info fields are provided
    hasInfo := ep.GitlabURL != nil || ep.DockerContainerName != nil || ep.KubernetesPodName != nil ||
        len(ep.Tags) > 0 || ep.Description != nil || ep.LastChangedBy != nil
    if hasInfo {
        insertInfoQuery := `
            INSERT INTO endpoint_info (
                endpoint_id,
                gitlab_url,
                docker_container_name,
                kubernetes_pod_name,
                tags,
                description,
                last_changed_by
            )
            VALUES ($1, $2, $3, $4, $5, $6, $7)
            RETURNING id
        `
        var infoID int
        err = tx.QueryRow(ctx, insertInfoQuery,
            newEp.ID,
            ep.GitlabURL,
            ep.DockerContainerName,
            ep.KubernetesPodName,
            ep.Tags,
            ep.Description,
            ep.LastChangedBy,
        ).Scan(&infoID)
        if err != nil {
            return nil, fmt.Errorf("failed to insert endpoint_info: %w", err)
        }
    }

    // Commit the transaction
    if err := tx.Commit(ctx); err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %w", err)
    }

    return newEp, nil
}

func (r *PostgresRepository) GetEndpointExpectedCode(ctx context.Context, id int) (int, error) {
	query := `
		SELECT expected_status_code
		FROM endpoints
		WHERE id = $1
	`

	var expectedCode int
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&expectedCode)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	return expectedCode, nil
}

func (r *PostgresRepository) GetEndPointByID(ctx context.Context, id int) (*EndpointDetail, error) {
	query := `
		SELECT 
			e.id,
			e.service_name,
			e.url,
			e.server_name,
			e.api_method,
			e.expected_status_code,

			-- Endpoint Stats
			COALESCE(es.endpoint_id, e.id) AS endpoint_id,
			COALESCE(es.total_checks, 0) AS total_checks,
			COALESCE(es.avg_latency, 0) AS avg_latency,
			COALESCE(es.successful_checks, 0) AS successful_checks,
			COALESCE(es.uptime_percentage, 0) AS uptime_percentage,
			COALESCE(es.failure_count, 0) AS failure_count,
			CASE 
				WHEN es.last_run = true THEN 'success'
				ELSE 'failure'
			END AS last_run,

			-- Endpoint Info
			COALESCE(i.description, '-') AS description,
			COALESCE(i.gitlab_url, '-') AS gitlab_url,
			COALESCE(i.docker_container_name, '-') AS docker_container_name,
			COALESCE(i.kubernetes_pod_name, '-') AS kubernetes_pod_name,
			COALESCE(i.tags, ARRAY[]::TEXT[]) AS tags,
			COALESCE(i.has_been_modified, false) AS has_been_modified,
			COALESCE(i.last_changed_by, '-') AS last_changed_by,
			COALESCE(i.created_at, NOW()) AS created_at,
			COALESCE(i.updated_at, NOW()) AS updated_at

		FROM endpoints e
		LEFT JOIN endpoint_info i 
			ON e.id = i.endpoint_id
		LEFT JOIN endpoint_stats es 
			ON e.id = es.endpoint_id
		WHERE e.id = $1;
	`

	var detail EndpointDetail
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&detail.ID,
		&detail.ServiceName,
		&detail.URL,
		&detail.ServerName,
		&detail.APIMethod,
		&detail.ExpectedCode,

		// Stats
		&detail.EndpointID,
		&detail.TotalChecks,
		&detail.AvgLatency,
		&detail.SuccessfulChecks,
		&detail.UptimePercentage,
		&detail.FailureCount,
		&detail.LastRunSucceeded,
		// Info
		&detail.Description,
		&detail.GitlabURL,
		&detail.DockerContainerName,
		&detail.KubernetesPodName,
		&detail.Tags,
		&detail.HasBeenModified,
		&detail.LastChangedBy,
		&detail.CreatedAt,
		&detail.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			log.Printf("endpoint %v not found", id)
			return nil, fmt.Errorf("endpoint %v not found", id)
		default:
			log.Printf("db query failed for endpoint %v: %v", id, err)
			return nil, fmt.Errorf("failed to query endpoint %v: %w", id, err)
		}
	}

	return &detail, nil
}


func (r *PostgresRepository) GetAllEndpointEssentials(ctx context.Context) (*[]EndpointBasicsDTO, error) {
	query := `
SELECT
    e.id,
    e.service_name,
    e.server_name,
    e.url,
    COALESCE(i.total_checks, 0) AS total_checks,
    COALESCE(i.uptime_percentage, 0) AS uptime_percentage,
    COALESCE(i.successful_checks, 0) AS successful_checks,
    COALESCE(i.avg_latency, 0) AS avg_latency,
    COALESCE(i.last_run, false) AS last_run,
    COALESCE(i.failure_count, 0) AS failure_count
FROM endpoints e
LEFT JOIN endpoint_stats i ON e.id = i.endpoint_id;
`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var basics []EndpointBasics
	for rows.Next() {
		var b EndpointBasics
		if err := rows.Scan(
			&b.ID,
			&b.ServiceName,
			&b.ServerName,
			&b.URL,
			&b.TotalChecks,
			&b.UptimePercentage,
			&b.SuccessfulChecks,
			&b.AverageLatency,
			&b.LastRun,
			&b.FailureCount,
		); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		basics = append(basics, b)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows iteration error: %w", rows.Err())
	}

	// Step 2: Transform into []EndpointBasicsDTO with calculations
	dtos := make([]EndpointBasicsDTO, len(basics))
	for i, b := range basics {
		downTimeCount := b.TotalChecks - b.SuccessfulChecks
		dtos[i] = EndpointBasicsDTO{
			ID:               b.ID,
			ServiceName:      b.ServiceName,
			ServerName:       b.ServerName,
			URL:              b.URL,
			TotalChecks:      b.TotalChecks,
			UptimePercentage: b.UptimePercentage,
			DownTimeCount:    downTimeCount,
			SuccessfulChecks: b.SuccessfulChecks,
			AverageLatency:   b.AverageLatency,
			LastRun:          b.LastRun,
			FailureCount:     b.FailureCount,
		}
	}
	return &dtos, nil
}

func (r *PostgresRepository) GetAggregateStats(ctx context.Context) (*AggregateDTO, error) {
	// Reuse existing essentials fetch
	endpoints, err := r.GetAllEndpointEssentials(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch endpoint essentials: %w", err)
	}

	if len(*endpoints) == 0 {
		return &AggregateDTO{}, nil
	}

	var totalChecks, totalSuccess, totalDown, weightedLatencySum int
	var latencyWeight float64

	for _, e := range *endpoints {
		totalChecks += e.TotalChecks
		totalSuccess += e.SuccessfulChecks
		totalDown += e.DownTimeCount

		// Weighted sum of latency (multiply by number of checks to avoid bias)
		weightedLatencySum += int(e.AverageLatency * float64(e.TotalChecks))
		latencyWeight += float64(e.TotalChecks)
	}

	// Compute uptime % safely
	var overallUptime float64
	if totalChecks > 0 {
		overallUptime = (float64(totalSuccess) / float64(totalChecks)) * 100
	}

	// Compute weighted avg latency
	var avgLatency float64
	if latencyWeight > 0 {
		avgLatency = float64(weightedLatencySum) / latencyWeight
	}

	return &AggregateDTO{
		TotalEndpoints:   len(*endpoints),
		TotalChecks:      totalChecks,
		SuccessfulChecks: totalSuccess,
		DownTimeCount:    totalDown,
		OverallUptime:    overallUptime,
		AverageLatency:   avgLatency,
	}, nil
}

// func (r *PostgresRepository) GetEndpointByID(ctx context.Context, id int64) (*Endpoint, error) {
// 	query := `SELECT id, name, url, status, last_checked FROM endpoints WHERE id = $1`
// 	row := r.db.Pool.QueryRow(ctx, query, id)

// 	var ep Endpoint
// 	err := row.Scan(&ep.ID, &ep.Name, &ep.URL, &ep.Status, &ep.LastChecked)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &ep, nil
// }

// func (r *PostgresRepository) SaveEndpoint(ctx context.Context, ep *Endpoint) error {
// 	query := `
// 		INSERT INTO endpoints (name, url, status, last_checked)
// 		VALUES ($1, $2, $3, $4)
// 		RETURNING id
// 	`
// 	return r.db.Pool.QueryRow(ctx, query, ep.Name, ep.URL, ep.Status, ep.LastChecked).Scan(&ep.ID)
// }

// func (r *PostgresRepository) UpdateEndpointStatus(ctx context.Context, id int64, status string, lastChecked time.Time) error {
// 	query := `UPDATE endpoints SET status = $1, last_checked = $2 WHERE id = $3`
// 	_, err := r.db.Pool.Exec(ctx, query, status, lastChecked, id)
// 	return err
// }

// func (r *PostgresRepository) ListEndpoints(ctx context.Context) ([]*Endpoint, error) {
// 	query := `SELECT id, name, url, status, last_checked FROM endpoints`
// 	rows, err := r.db.Pool.Query(ctx, query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var endpoints []*Endpoint
// 	for rows.Next() {
// 		var ep Endpoint
// 		if err := rows.Scan(&ep.ID, &ep.Name, &ep.URL, &ep.Status, &ep.LastChecked); err != nil {
// 			return nil, err
// 		}
// 		endpoints = append(endpoints, &ep)
// 	}

// 	return endpoints, nil
// }
