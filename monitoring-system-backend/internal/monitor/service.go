package monitor

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/badgerv/monitoring-api/internal/storage"

	"encoding/json"
	"log"
	"os"
	"time"
)

type Service struct {
	db     *storage.DB
	dbRepo *PostgresRepository
}

func NewService(db *storage.DB, dbRepo *PostgresRepository) *Service {
	return &Service{db: db, dbRepo: dbRepo}
}

func (s *Service) LoadAndSyncEndpoints(path string) ([]Endpoint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Fetch DB endpoints
	dbEndpoints, err := s.dbRepo.GetAllEndpoints(ctx)
	if err != nil {
		return nil, err
	}

	// Try to read JSON file
	file, err := os.ReadFile(path)
	if err != nil {
		// No JSON file, just return what’s in DB
		return dbEndpoints, nil
	}

	var jsonEndpoints []Endpoint
	if err := json.Unmarshal(file, &jsonEndpoints); err != nil {
		return dbEndpoints, err
	}

	// Normalize helper
	normalizeURL := func(raw string) string {
		raw = strings.TrimSpace(raw)
		if strings.HasSuffix(raw, "/") {
			raw = strings.TrimSuffix(raw, "/")
		}
		return strings.ToLower(raw) // also force lowercase
	}

	// Merge JSON entries into DB
	for _, jep := range jsonEndpoints {
		jep.URL = normalizeURL(jep.URL)

		exists := false
		for _, dep := range dbEndpoints {
			if normalizeURL(dep.URL) == jep.URL &&
				dep.APIMethod == jep.APIMethod &&
				dep.ServerName == jep.ServerName &&
				dep.ExpectedCode == jep.ExpectedCode {
				exists = true
				break
			}
		}

		if !exists {
			var newID int
			err = s.db.Pool.QueryRow(ctx, `
				INSERT INTO endpoints (service_name, url, server_name, api_method, expected_status_code)
				VALUES ($1, $2, $3, $4, $5)
				ON CONFLICT (url) DO NOTHING
				RETURNING id
			`,
				jep.ServiceName, jep.URL, jep.ServerName, jep.APIMethod, jep.ExpectedCode,
			).Scan(&newID)

			if err != nil {
				// When ON CONFLICT DO NOTHING triggers, no rows are returned
				if err.Error() == "no rows in result set" {
					continue
				}
				log.Printf("Error inserting endpoint %s: %v", jep.URL, err)
				continue
			}

			jep.ID = newID
			dbEndpoints = append(dbEndpoints, jep)
		}
	}

	return dbEndpoints, nil
}

// func (s *Service) CheckEndpoint(ctx context.Context, endpointID int, url, method, serverName, serviceName string) error {
// 	start := time.Now()
// 	log.Println("Calling the API endpoint for", url)

// 	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
// 	client := &http.Client{Timeout: 5 * time.Second}

// 	resp, err := client.Do(req)

// 	log.Println(resp.StatusCode)

// 	if err != nil {
// 		log.Println("An error occured while trying to call endpoint %v", err)
// 	}
// 	latency := time.Since(start).Milliseconds()

// 	var errMsg string
// 	var statusCode int
// 	var success int64 // will be 1 if success, else 0
// 	var lastrun bool

// 	expectedStatusCode, err := s.dbRepo.GetEndpointExpectedCode(ctx, endpointID)

// 	if err != nil {
// 		errMsg = err.Error()
// 	} else {
// 		statusCode = resp.StatusCode
// 		if statusCode == expectedStatusCode {
// 			success = 1
// 			lastrun = true
// 		}
// 		resp.Body.Close()
// 	}

// 	// Insert into checks log table
// 	_, insertErr := s.db.Pool.Exec(ctx,
// 		`INSERT INTO checks (endpoint_id, status_code, latency_ms, error)
//          VALUES ($1, $2, $3, $4)`,
// 		endpointID, statusCode, latency, errMsg,
// 	)

// 	if insertErr != nil {
// 		return insertErr
// 	}

// 	// Update stats table
// 	_, statsErr := s.db.Pool.Exec(ctx,
// 		`INSERT INTO endpoint_stats (endpoint_id, total_checks, total_latency, successful_checks, last_run)
// 	 VALUES ($1, 1, $2, $3, $4)
// 	 ON CONFLICT (endpoint_id) DO UPDATE
// 	 SET total_checks = endpoint_stats.total_checks + 1,
// 	     total_latency = endpoint_stats.total_latency + EXCLUDED.total_latency,
// 	     successful_checks = endpoint_stats.successful_checks + EXCLUDED.successful_checks,
// 	     last_run = EXCLUDED.last_run`,
// 		endpointID, latency, success, lastrun,
// 	)

// 	if statsErr != nil {
// 		return statsErr
// 	}

// 	return nil
// }


func (s *Service) CheckEndpoint(ctx context.Context, endpointID int, url, method, serverName, serviceName string) error {
	start := time.Now()
	log.Println("Calling the API endpoint for", url)

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Do(req)

	var statusCode int
	var errMsg string
	var success int64 // will be 1 if success, else 0
	var failure int64 // will be 1 if failure, else 0
	var lastrun bool

	if err != nil {
		log.Printf("An error occurred while trying to call endpoint: %v", err)
		errMsg = err.Error()
	} else {
		log.Println(resp.StatusCode)
		statusCode = resp.StatusCode
	}

	expectedStatusCode, err := s.dbRepo.GetEndpointExpectedCode(ctx, endpointID)

	if err != nil {
		errMsg = err.Error()
		failure = 1
	} else if statusCode == expectedStatusCode && err == nil {
		success = 1
		lastrun = true
		failure = 0
	} else {
		failure = 1
	}

	if resp != nil {
		resp.Body.Close()
	}

	latency := time.Since(start).Milliseconds()

	// Insert into checks log table
	_, insertErr := s.db.Pool.Exec(ctx,
		`INSERT INTO checks (endpoint_id, status_code, latency_ms, error)
         VALUES ($1, $2, $3, $4)`,
		endpointID, statusCode, latency, errMsg,
	)

	if insertErr != nil {
		return insertErr
	}

	// Update stats table
	_, statsErr := s.db.Pool.Exec(ctx,
		`INSERT INTO endpoint_stats (endpoint_id, total_checks, total_latency, successful_checks, failure_count, last_run)
	 VALUES ($1, 1, $2, $3, $4, $5)
	 ON CONFLICT (endpoint_id) DO UPDATE
	 SET total_checks = endpoint_stats.total_checks + 1,
	     total_latency = endpoint_stats.total_latency + EXCLUDED.total_latency,
	     successful_checks = endpoint_stats.successful_checks + EXCLUDED.successful_checks,
	     failure_count = CASE WHEN EXCLUDED.successful_checks = 1 THEN 0 ELSE endpoint_stats.failure_count + 1 END,
	     last_run = EXCLUDED.last_run`,
		endpointID, latency, success, failure, lastrun,
	)

	if statsErr != nil {
		return statsErr
	}

	return nil
}

func (s *Service) CheckEndpointStatus(ctx context.Context, url, method string, expectedStatus int) (time.Duration, error) {
    // Create the HTTP request with context
    req, err := http.NewRequestWithContext(ctx, method, url, nil)
    if err != nil {
        return 0, err
    }

    // HTTP client with timeout
    client := &http.Client{Timeout: 5 * time.Second}

    // Measure start time
    start := time.Now()

    // Send the request
    resp, err := client.Do(req)
    latency := time.Since(start) // calculate latency
    if err != nil {
        return latency, err
    }
    defer resp.Body.Close() // ensure body is closed

    // Check if status code matches expected
    if resp.StatusCode != expectedStatus {
        return latency, fmt.Errorf("unexpected status code: got %d, expected %d", resp.StatusCode, expectedStatus)
    }

    // Success
    return latency, nil
}



func (s *Service) GetEndpointByID(ctx context.Context, endpointID int) (*EndpointDetail, error) {
	return s.dbRepo.GetEndPointByID(ctx, endpointID)
}

func (s *Service) GetAllEndpointEssentials(ctx context.Context) (*[]EndpointBasicsDTO, error) {
	return s.dbRepo.GetAllEndpointEssentials(ctx)
}

func (s *Service) GetAggregateStats(ctx context.Context) (*AggregateDTO, error) {
	return s.dbRepo.GetAggregateStats(ctx)
}

// Exposes repo function to the handler
func (s *Service) CreateEndpoint(ctx context.Context, ep *Endpoint) (*Endpoint, error) {
	return s.dbRepo.CreateEndpoint(ctx, ep)
}

// func (s *Service) GetEndpointEssentials(ctx context.Context) ([]EndpointBasicsDTO, error) {
// }
