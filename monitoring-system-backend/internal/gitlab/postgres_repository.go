package gitlab

import (
	"context"
	"fmt"
	"time"

	"database/sql"

	"github.com/badgerv/monitoring-api/internal/rbac"
	"github.com/badgerv/monitoring-api/internal/storage"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// PostgresRepository implements the Repository interface using PostgreSQL with pgx.
type PostgresRepository struct {
	db     *storage.DB
	logger *zap.Logger
	rbac   *rbac.Service
}

// NewPostgresRepository creates a new PostgresRepository with the given storage.DB and logger.
func NewPostgresRepository(db *storage.DB, logger *zap.Logger, rbac *rbac.Service) (*PostgresRepository, error) {
	repo := &PostgresRepository{db: db, logger: logger, rbac: rbac}
	if err := repo.initTables(); err != nil {
		return nil, err
	}
	return repo, nil
}

// initTables creates the necessary database tables if they don't exist.
func (r *PostgresRepository) initTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS services (
			id TEXT PRIMARY KEY,
			gitlab_repo_id TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			url TEXT NOT NULL,
			type TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS pipeline_units (
			id TEXT PRIMARY KEY,
			macro_service_id TEXT,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			FOREIGN KEY (macro_service_id) REFERENCES services(id)
		)`,
		`CREATE TABLE IF NOT EXISTS pipeline_dependencies (
			pipeline_unit_id TEXT,
			micro_service_id TEXT,
			order_index INTEGER NOT NULL,
			PRIMARY KEY (pipeline_unit_id, micro_service_id),
			FOREIGN KEY (pipeline_unit_id) REFERENCES pipeline_units(id),
			FOREIGN KEY (micro_service_id) REFERENCES services(id)
		)`,
		`CREATE TABLE IF NOT EXISTS pipeline_runs (
			id TEXT PRIMARY KEY,
			pipeline_unit_id TEXT NOT NULL,
			status TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			gitlab_pipeline_id INTEGER,
			approver_id TEXT,
			execution_time BIGINT,
			selected_micro_service_ids TEXT[],  -- New field as TEXT array
			FOREIGN KEY (pipeline_unit_id) REFERENCES pipeline_units(id)
		)`,
		`CREATE TABLE IF NOT EXISTS authorization_requests (
			id TEXT PRIMARY KEY,
			pipeline_run_id TEXT NOT NULL,
			requester_id TEXT NOT NULL,
			approver_id TEXT NOT NULL,
			status TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			comment TEXT,
			FOREIGN KEY (pipeline_run_id) REFERENCES pipeline_runs(id)
		)`,
		`CREATE TABLE IF NOT EXISTS execution_history (
			id TEXT PRIMARY KEY,
			pipeline_run_id TEXT NOT NULL,
			requester_id TEXT NOT NULL,
			approver_id TEXT NOT NULL,
			status TEXT NOT NULL,
			started_at TIMESTAMP NOT NULL,
			completed_at TIMESTAMP,
			execution_time BIGINT,
			error_message TEXT,
			FOREIGN KEY (pipeline_run_id) REFERENCES pipeline_runs(id)
		)`,
	}

	ctx := context.Background()
	for _, query := range queries {
		_, err := r.db.Pool.Exec(ctx, query)
		if err != nil {
			r.logger.Error("Failed to create table", zap.Error(err))
			return err
		}
	}
	return nil
}

// CreateService creates a new service entry.
func (r *PostgresRepository) CreateService(ctx context.Context, service Service) (Service, error) {
	if service.ID == "" {
		service.ID = uuid.New().String()
	}
	service.CreatedAt = time.Now()
	service.UpdatedAt = service.CreatedAt

	query := `INSERT INTO services (id, gitlab_repo_id, name, url, type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, gitlab_repo_id, name, url, type, created_at, updated_at`
	var createdService Service
	err := r.db.Pool.QueryRow(ctx, query, service.ID, service.GitLabRepoID, service.Name, service.URL, service.Type, service.CreatedAt, service.UpdatedAt).
		Scan(&createdService.ID, &createdService.GitLabRepoID, &createdService.Name, &createdService.URL, &createdService.Type, &createdService.CreatedAt, &createdService.UpdatedAt)
	if err != nil {
		r.logger.Error("Failed to create service", zap.Error(err))
		return Service{}, err
	}
	return createdService, nil
}

// GetServiceByID retrieves a service by its ID.
func (r *PostgresRepository) GetServiceByID(ctx context.Context, id string) (Service, error) {
	query := `SELECT id, gitlab_repo_id, name, url, type, created_at, updated_at FROM services WHERE id = $1`

	var service Service

	err := r.db.Pool.QueryRow(ctx, query, id).
		Scan(&service.ID, &service.GitLabRepoID, &service.Name, &service.URL, &service.Type, &service.CreatedAt, &service.UpdatedAt)
	if err == pgx.ErrNoRows {
		return Service{}, nil
	}
	fmt.Println("This is the microservice", service)
	if err != nil {
		r.logger.Error("Failed to get service by ID", zap.String("id", id), zap.Error(err))
		return Service{}, err
	}

	return service, nil
}

// GetServiceByGitLabRepoID retrieves a service by its GitLab repository ID.
func (r *PostgresRepository) GetServiceByGitLabRepoID(ctx context.Context, gitlabRepoID string) (Service, error) {
	query := `SELECT id, gitlab_repo_id, name, url, type, created_at, updated_at FROM services WHERE gitlab_repo_id = $1`
	var service Service
	err := r.db.Pool.QueryRow(ctx, query, gitlabRepoID).
		Scan(&service.ID, &service.GitLabRepoID, &service.Name, &service.URL, &service.Type, &service.CreatedAt, &service.UpdatedAt)
	if err == pgx.ErrNoRows {
		return Service{}, nil
	}
	if err != nil {
		r.logger.Error("Failed to get service by GitLab repo ID", zap.String("gitlab_repo_id", gitlabRepoID), zap.Error(err))
		return Service{}, err
	}
	return service, nil
}

func (r *PostgresRepository) ListServicesByType(ctx context.Context, serviceType ServiceType) ([]Service, error) {
	var query string
	var args []interface{}
	if serviceType == "" {
		query = `SELECT id, gitlab_repo_id, name, url, type, created_at, updated_at FROM services`
	} else {
		query = `SELECT id, gitlab_repo_id, name, url, type, created_at, updated_at FROM services WHERE type = $1`
		args = []interface{}{serviceType}
	}
	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to list services", zap.String("type", string(serviceType)), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var services []Service
	for rows.Next() {
		var s Service
		if err := rows.Scan(&s.ID, &s.GitLabRepoID, &s.Name, &s.URL, &s.Type, &s.CreatedAt, &s.UpdatedAt); err != nil {
			r.logger.Error("Failed to scan service", zap.Error(err))
			return nil, err
		}
		services = append(services, s)
	}
	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating services", zap.Error(err))
		return nil, err
	}
	return services, nil
}

// CreatePipelineUnit creates a new pipeline unit.
func (r *PostgresRepository) CreatePipelineUnit(ctx context.Context, unit PipelineUnit) (PipelineUnit, error) {
	if unit.ID == "" {
		unit.ID = uuid.New().String()
	}
	unit.CreatedAt = time.Now()
	unit.UpdatedAt = unit.CreatedAt

	query := `INSERT INTO pipeline_units (id, macro_service_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, macro_service_id, created_at, updated_at`
	var createdUnit PipelineUnit
	err := r.db.Pool.QueryRow(ctx, query, unit.ID, unit.MacroServiceID, unit.CreatedAt, unit.UpdatedAt).
		Scan(&createdUnit.ID, &createdUnit.MacroServiceID, &createdUnit.CreatedAt, &createdUnit.UpdatedAt)
	if err != nil {
		r.logger.Error("Failed to create pipeline unit", zap.Error(err))
		return PipelineUnit{}, err
	}
	return createdUnit, nil
}

// GetPipelineUnit retrieves a pipeline unit by its ID, including its microservice dependencies.
func (r *PostgresRepository) GetPipelineUnit(ctx context.Context, id string) (PipelineUnit, error) {
	query := `SELECT id, macro_service_id, created_at, updated_at
		FROM pipeline_units WHERE id = $1`
	var unit PipelineUnit
	err := r.db.Pool.QueryRow(ctx, query, id).
		Scan(&unit.ID, &unit.MacroServiceID, &unit.CreatedAt, &unit.UpdatedAt)
	if err == pgx.ErrNoRows {
		return PipelineUnit{}, nil
	}
	if err != nil {
		r.logger.Error("Failed to get pipeline unit", zap.String("id", id), zap.Error(err))
		return PipelineUnit{}, err
	}

	// Fetch microservice dependencies
	unit.MicroServiceIDs, err = r.GetMicroServiceDependencies(ctx, id)
	if err != nil {
		return PipelineUnit{}, err
	}
	return unit, nil
}

func (r *PostgresRepository) GetPipelineUnitWithServices(ctx context.Context, id string) (PipelineUnit, Service, []Service, error) {
	query := `
		SELECT 
			pu.id, pu.macro_service_id, pu.created_at, pu.updated_at,
			ms.id, ms.gitlab_repo_id, ms.name, ms.url, ms.type, ms.created_at, ms.updated_at,
			mics.id, mics.gitlab_repo_id, mics.name, mics.url, mics.type, mics.created_at, mics.updated_at
		FROM pipeline_units pu
		LEFT JOIN services ms ON pu.macro_service_id = ms.id
		LEFT JOIN pipeline_dependencies pd ON pu.id = pd.pipeline_unit_id
		LEFT JOIN services mics ON pd.micro_service_id = mics.id
		WHERE pu.id = $1
		ORDER BY pd.order_index`
	rows, err := r.db.Pool.Query(ctx, query, id)
	if err != nil {
		r.logger.Error("Failed to get pipeline unit with services", zap.String("id", id), zap.Error(err))
		return PipelineUnit{}, Service{}, nil, err
	}
	defer rows.Close()

	var unit PipelineUnit
	var macroService Service
	var microServices []Service
	var firstRow = true

	for rows.Next() {
		var microService Service
		var msID, msGitLabRepoID, msName, msURL, msType *string
		var msCreatedAt, msUpdatedAt *time.Time
		var micsID, micsGitLabRepoID, micsName, micsURL, micsType *string
		var micsCreatedAt, micsUpdatedAt *time.Time

		err := rows.Scan(
			&unit.ID, &unit.MacroServiceID, &unit.CreatedAt, &unit.UpdatedAt,
			&msID, &msGitLabRepoID, &msName, &msURL, &msType, &msCreatedAt, &msUpdatedAt,
			&micsID, &micsGitLabRepoID, &micsName, &micsURL, &micsType, &micsCreatedAt, &micsUpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan pipeline unit with services", zap.String("id", id), zap.Error(err))
			return PipelineUnit{}, Service{}, nil, err
		}

		if firstRow {
			if msID != nil {
				macroService = Service{
					ID:           *msID,
					GitLabRepoID: *msGitLabRepoID,
					Name:         *msName,
					URL:          *msURL,
					Type:         ServiceType(*msType),
					CreatedAt:    *msCreatedAt,
					UpdatedAt:    *msUpdatedAt,
				}
			}
			firstRow = false
		}

		if micsID != nil {
			microService = Service{
				ID:           *micsID,
				GitLabRepoID: *micsGitLabRepoID,
				Name:         *micsName,
				URL:          *micsURL,
				Type:         ServiceType(*micsType),
				CreatedAt:    *micsCreatedAt,
				UpdatedAt:    *micsUpdatedAt,
			}
			microServices = append(microServices, microService)
		}
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating pipeline unit with services", zap.String("id", id), zap.Error(err))
		return PipelineUnit{}, Service{}, nil, err
	}

	if unit.ID == "" {
		return PipelineUnit{}, Service{}, nil, nil // No rows found
	}

	unit.MicroServiceIDs = make([]string, len(microServices))
	for i, svc := range microServices {
		unit.MicroServiceIDs[i] = svc.ID
	}

	return unit, macroService, microServices, nil
}

// AddMicroServiceDependency adds a microservice dependency to a pipeline unit with an order index.
func (r *PostgresRepository) AddMicroServiceDependency(ctx context.Context, pipelineUnitID, microServiceID string, orderIndex int) error {
	query := `INSERT INTO pipeline_dependencies (pipeline_unit_id, micro_service_id, order_index)
		VALUES ($1, $2, $3)`
	_, err := r.db.Pool.Exec(ctx, query, pipelineUnitID, microServiceID, orderIndex)
	if err != nil {
		r.logger.Error("Failed to add microservice dependency", zap.String("pipeline_unit_id", pipelineUnitID), zap.String("micro_service_id", microServiceID), zap.Error(err))
		return err
	}
	return nil
}

// ListAllAuthorizationRequests retrieves all authorization requests from the database.

// ListAllAuthorizationRequests retrieves all authorization requests with requester and approver names.
// ListAllAuthorizationRequests retrieves all authorization requests with requester, approver, and service names.
func (r *PostgresRepository) ListAllAuthorizationRequests(ctx context.Context) ([]AuthorizationRequest, error) {
	query := `
SELECT
    ar.id,
    ar.pipeline_run_id,
    ar.requester_id,
    u1.username AS requester_name,
    ar.approver_id AS approver_id,
    u2.username AS approver_name,
    ar.status,
    ar.created_at,
    ar.updated_at,
    ar.comment,
    s_macro.name AS macro_service_name,
    COALESCE(array_agg(DISTINCT s_micro.name), '{}') AS micro_service_names
FROM authorization_requests ar
JOIN users u1 
    ON ar.requester_id::uuid = u1.id
LEFT JOIN users u2 ON ar.approver_id = u2.id

JOIN pipeline_runs pr 
    ON ar.pipeline_run_id = pr.id
JOIN pipeline_units pu 
    ON pr.pipeline_unit_id = pu.id
LEFT JOIN services s_macro 
    ON pu.macro_service_id = s_macro.id
LEFT JOIN pipeline_dependencies pd 
    ON pu.id = pd.pipeline_unit_id
LEFT JOIN services s_micro 
    ON pd.micro_service_id = s_micro.id
GROUP BY ar.id, u1.username, u2.username, s_macro.name;
	`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		r.logger.Error("Failed to list all authorization requests", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var requests []AuthorizationRequest
	for rows.Next() {
		var req AuthorizationRequest
		var requesterName, approverName, macroServiceName sql.NullString
		var microServiceNames []string
		if err := rows.Scan(
			&req.ID,
			&req.PipelineRunID,
			&req.RequesterID,
			&requesterName,
			&req.ApproverID,
			&approverName,
			&req.Status,
			&req.CreatedAt,
			&req.UpdatedAt,
			&req.Comment,
			&macroServiceName,
			&microServiceNames,
		); err != nil {
			r.logger.Error("Failed to scan authorization request", zap.Error(err))
			return nil, err
		}
		req.RequesterName = requesterName.String
		req.ApproverName = approverName.String
		req.MacroServiceName = macroServiceName.String
		req.MicroServiceNames = microServiceNames
		requests = append(requests, req)
	}
	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating authorization requests", zap.Error(err))
		return nil, err
	}
	return requests, nil
}

func (r *PostgresRepository) GetAuthorizationRequestByID(ctx context.Context, id string) (*AuthorizationRequest, error) {
	query := `
SELECT
    ar.id,
    ar.pipeline_run_id,
    ar.requester_id,
    u1.username AS requester_name,
    ar.approver_id AS approver_id,
    u2.username AS approver_name,
    ar.status,
    ar.created_at,
    ar.updated_at,
    ar.comment,
    s_macro.name AS macro_service_name,
    COALESCE(array_agg(DISTINCT s_micro.name), '{}') AS micro_service_names
FROM authorization_requests ar
JOIN users u1 
    ON ar.requester_id::uuid = u1.id
LEFT JOIN users u2 
    ON ar.approver_id = u2.id
JOIN pipeline_runs pr 
    ON ar.pipeline_run_id = pr.id
JOIN pipeline_units pu 
    ON pr.pipeline_unit_id = pu.id
LEFT JOIN services s_macro 
    ON pu.macro_service_id = s_macro.id
LEFT JOIN pipeline_dependencies pd 
    ON pu.id = pd.pipeline_unit_id
LEFT JOIN services s_micro 
    ON pd.micro_service_id = s_micro.id
WHERE ar.id = $1
GROUP BY ar.id, u1.username, u2.username, s_macro.name;
`

	row := r.db.Pool.QueryRow(ctx, query, id)

	var req AuthorizationRequest
	var requesterName, approverName, macroServiceName sql.NullString
	var microServiceNames []string

	if err := row.Scan(
		&req.ID,
		&req.PipelineRunID,
		&req.RequesterID,
		&requesterName,
		&req.ApproverID,
		&approverName,
		&req.Status,
		&req.CreatedAt,
		&req.UpdatedAt,
		&req.Comment,
		&macroServiceName,
		&microServiceNames,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // not found
		}
		r.logger.Error("Failed to fetch authorization request by ID", zap.Error(err))
		return nil, err
	}

	req.RequesterName = requesterName.String
	req.ApproverName = approverName.String
	req.MacroServiceName = macroServiceName.String
	req.MicroServiceNames = microServiceNames

	return &req, nil
}


// ListAllExecutionHistories retrieves all execution history records with requester, approver, and service names.
func (r *PostgresRepository) ListAllExecutionHistories(ctx context.Context, id string) ([]ExecutionHistory, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	roles, err := r.rbac.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	isSuperAdmin := false
	for _, role := range roles {
		if role.Name == "super admin" {
			isSuperAdmin = true
			break
		}
	}

	// Base query
	query := `
		SELECT
			eh.id,
			eh.pipeline_run_id,
			eh.requester_id,
			u1.username AS requester_name,
			eh.approver_id,
			u2.username AS approver_name,
			eh.status,
			eh.started_at,
			eh.completed_at,
			eh.error_message,
			pu.id AS pipeline_unit_id,
			s_macro.name AS macro_service_name,
			COALESCE((
				SELECT array_agg(s_micro.name)
				FROM unnest(pr.selected_micro_service_ids) AS micro_id
				JOIN services s_micro ON s_micro.id::uuid = micro_id::uuid
			), '{}') AS micro_service_names
		FROM execution_history eh
		JOIN users u1 ON eh.requester_id::uuid = u1.id
		LEFT JOIN users u2 ON eh.approver_id::uuid = u2.id
		JOIN pipeline_runs pr ON eh.pipeline_run_id = pr.id
		JOIN pipeline_units pu ON pr.pipeline_unit_id = pu.id
		LEFT JOIN services s_macro ON pu.macro_service_id = s_macro.id
	`

	// Add condition if not super admin
	var rows pgx.Rows
	if isSuperAdmin {
		rows, err = r.db.Pool.Query(ctx, query)
	} else {
		query += ` WHERE eh.requester_id::uuid = $1`
		rows, err = r.db.Pool.Query(ctx, query, userID)
	}

	if err != nil {
		r.logger.Error("Failed to list execution histories", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var histories []ExecutionHistory
	for rows.Next() {
		var h ExecutionHistory
		var requesterName, approverName, macroServiceName, errorMessage sql.NullString
		var microServiceNames []string
		var completedAt sql.NullTime

		if err := rows.Scan(
			&h.ID,
			&h.PipelineRunID,
			&h.RequesterID,
			&requesterName,
			&h.ApproverID,
			&approverName,
			&h.Status,
			&h.StartedAt,
			&completedAt,
			&errorMessage,
			&h.PipelineUnitID,
			&macroServiceName,
			&microServiceNames,
		); err != nil {
			r.logger.Error("Failed to scan execution history", zap.Error(err))
			return nil, err
		}

		h.RequesterName = requesterName.String
		h.ApproverName = approverName.String
		h.CompletedAt = completedAt.Time
		h.ErrorMessage = errorMessage.String
		h.MacroServiceName = macroServiceName.String
		h.MicroServiceNames = microServiceNames

		histories = append(histories, h)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating execution histories", zap.Error(err))
		return nil, err
	}

	return histories, nil
}

func (r *PostgresRepository) GetExecutionHistoryByID(ctx context.Context, userIDStr, historyID string) (*ExecutionHistory, error) {
	// Parse user UUID
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	// Get user roles
	roles, err := r.rbac.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	isSuperAdmin := false
	for _, role := range roles {
		if role.Name == "super admin" {
			isSuperAdmin = true
			break
		}
	}

	// Base query (same as list but filtered by eh.id)
	query := `
		SELECT
			eh.id,
			eh.pipeline_run_id,
			eh.requester_id,
			u1.username AS requester_name,
			eh.approver_id,
			u2.username AS approver_name,
			eh.status,
			eh.started_at,
			eh.completed_at,
			eh.error_message,
			pu.id AS pipeline_unit_id,
			s_macro.name AS macro_service_name,
			COALESCE((
				SELECT array_agg(s_micro.name)
				FROM unnest(pr.selected_micro_service_ids) AS micro_id
				JOIN services s_micro ON s_micro.id::uuid = micro_id::uuid
			), '{}') AS micro_service_names
		FROM execution_history eh
		JOIN users u1 ON eh.requester_id::uuid = u1.id
		LEFT JOIN users u2 ON eh.approver_id::uuid = u2.id
		JOIN pipeline_runs pr ON eh.pipeline_run_id = pr.id
		JOIN pipeline_units pu ON pr.pipeline_unit_id = pu.id
		LEFT JOIN services s_macro ON pu.macro_service_id = s_macro.id
		WHERE eh.id = $1
	`

	// If not super admin, restrict to user’s own histories
	args := []interface{}{historyID}
	if !isSuperAdmin {
		query += ` AND eh.requester_id::uuid = $2`
		args = append(args, userID)
	}

	row := r.db.Pool.QueryRow(ctx, query, args...)

	var h ExecutionHistory
	var requesterName, approverName, macroServiceName, errorMessage sql.NullString
	var microServiceNames []string
	var completedAt sql.NullTime

	if err := row.Scan(
		&h.ID,
		&h.PipelineRunID,
		&h.RequesterID,
		&requesterName,
		&h.ApproverID,
		&approverName,
		&h.Status,
		&h.StartedAt,
		&completedAt,
		&errorMessage,
		&h.PipelineUnitID,
		&macroServiceName,
		&microServiceNames,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // not found
		}
		r.logger.Error("Failed to fetch execution history by ID", zap.Error(err))
		return nil, err
	}

	h.RequesterName = requesterName.String
	h.ApproverName = approverName.String
	h.CompletedAt = completedAt.Time
	h.ErrorMessage = errorMessage.String
	h.MacroServiceName = macroServiceName.String
	h.MicroServiceNames = microServiceNames

	return &h, nil
}


// GetMicroServiceDependencies retrieves the ordered list of microservice IDs for a pipeline unit.
func (r *PostgresRepository) GetMicroServiceDependencies(ctx context.Context, pipelineUnitID string) ([]string, error) {
	query := `SELECT micro_service_id FROM pipeline_dependencies WHERE pipeline_unit_id = $1 ORDER BY order_index`
	rows, err := r.db.Pool.Query(ctx, query, pipelineUnitID)
	if err != nil {
		r.logger.Error("Failed to get microservice dependencies", zap.String("pipeline_unit_id", pipelineUnitID), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var microServiceIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			r.logger.Error("Failed to scan microservice dependency", zap.Error(err))
			return nil, err
		}
		microServiceIDs = append(microServiceIDs, id)
	}
	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating microservice dependencies", zap.Error(err))
		return nil, err
	}
	return microServiceIDs, nil
}

// CreatePipelineRun creates a new pipeline run.
func (r *PostgresRepository) CreatePipelineRun(ctx context.Context, run PipelineRun) (PipelineRun, error) {
	if run.ID == "" {
		run.ID = uuid.New().String()
	}
	run.CreatedAt = time.Now()
	run.UpdatedAt = run.CreatedAt

	query := `INSERT INTO pipeline_runs (id, pipeline_unit_id, status, created_at, updated_at, gitlab_pipeline_id, approver_id, execution_time, selected_micro_service_ids)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id, pipeline_unit_id, status, created_at, updated_at, gitlab_pipeline_id, approver_id, execution_time, selected_micro_service_ids`
	var createdRun PipelineRun
	err := r.db.Pool.QueryRow(ctx, query, run.ID, run.PipelineUnitID, run.Status, run.CreatedAt, run.UpdatedAt, run.GitLabPipelineID, run.ApproverID, run.ExecutionTime, run.SelectedMicroServiceIDs).
		Scan(&createdRun.ID, &createdRun.PipelineUnitID, &createdRun.Status, &createdRun.CreatedAt, &createdRun.UpdatedAt, &createdRun.GitLabPipelineID, &createdRun.ApproverID, &createdRun.ExecutionTime, &createdRun.SelectedMicroServiceIDs)
	if err != nil {
		r.logger.Error("Failed to create pipeline run", zap.Error(err))
		return PipelineRun{}, err
	}
	return createdRun, nil
}

// GetPipelineRun retrieves a pipeline run by its ID.
func (r *PostgresRepository) GetPipelineRun(ctx context.Context, id string) (PipelineRun, error) {
	query := `SELECT id, pipeline_unit_id, status, created_at, updated_at, gitlab_pipeline_id, approver_id, execution_time, selected_micro_service_ids
        FROM pipeline_runs WHERE id = $1`
	var run PipelineRun
	err := r.db.Pool.QueryRow(ctx, query, id).
		Scan(&run.ID, &run.PipelineUnitID, &run.Status, &run.CreatedAt, &run.UpdatedAt, &run.GitLabPipelineID, &run.ApproverID, &run.ExecutionTime, &run.SelectedMicroServiceIDs)
	if err == pgx.ErrNoRows {
		return PipelineRun{}, nil
	}
	if err != nil {
		r.logger.Error("Failed to get pipeline run", zap.String("id", id), zap.Error(err))
		return PipelineRun{}, err
	}
	return run, nil
}

// UpdatePipelineRunStatus updates the status of a pipeline run.
func (r *PostgresRepository) UpdatePipelineRunStatus(ctx context.Context, id string, status PipelineStatus) error {
	query := `UPDATE pipeline_runs SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Pool.Exec(ctx, query, status, time.Now(), id)
	if err != nil {
		r.logger.Error("Failed to update pipeline run status", zap.String("id", id), zap.String("status", string(status)), zap.Error(err))
		return err
	}
	return nil
}

// UpdatePipelineRun updates a pipeline run's fields (e.g., GitLabPipelineID).
func (r *PostgresRepository) UpdatePipelineRun(ctx context.Context, run PipelineRun) error {
	query := `UPDATE pipeline_runs SET status = $1, updated_at = $2, gitlab_pipeline_id = $3, approver_id = $4, execution_time = $5
		WHERE id = $6`
	_, err := r.db.Pool.Exec(ctx, query, run.Status, time.Now(), run.GitLabPipelineID, run.ApproverID, run.ExecutionTime, run.ID)
	if err != nil {
		r.logger.Error("Failed to update pipeline run", zap.String("id", run.ID), zap.Error(err))
		return err
	}
	return nil
}

// CreateAuthorizationRequest creates a new authorization request.
func (r *PostgresRepository) CreateAuthorizationRequest(ctx context.Context, request AuthorizationRequest) (AuthorizationRequest, error) {
	if request.ID == "" {
		request.ID = uuid.New().String()
	}
	request.CreatedAt = time.Now()
	request.UpdatedAt = request.CreatedAt

	query := `INSERT INTO authorization_requests (id, pipeline_run_id, requester_id, approver_id, status, created_at, updated_at, comment)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, pipeline_run_id, requester_id, approver_id, status, created_at, updated_at, comment`
	var createdRequest AuthorizationRequest
	err := r.db.Pool.QueryRow(ctx, query, request.ID, request.PipelineRunID, request.RequesterID, request.ApproverID, request.Status, request.CreatedAt, request.UpdatedAt, request.Comment).
		Scan(&createdRequest.ID, &createdRequest.PipelineRunID, &createdRequest.RequesterID, &createdRequest.ApproverID, &createdRequest.Status, &createdRequest.CreatedAt, &createdRequest.UpdatedAt, &createdRequest.Comment)
	if err != nil {
		r.logger.Error("Failed to create authorization request", zap.Error(err))
		return AuthorizationRequest{}, err
	}
	return createdRequest, nil
}

// GetAuthorizationRequest retrieves an authorization request by its ID.
func (r *PostgresRepository) GetAuthorizationRequest(ctx context.Context, id string) (AuthorizationRequest, error) {
	query := `
	SELECT id,
       pipeline_run_id,
       requester_id,
       COALESCE(approver_id::text, '') AS approver_id,
       status,
       created_at,
       updated_at,
       COALESCE(comment, '') AS comment
FROM authorization_requests
WHERE id = $1;
`

	var request AuthorizationRequest
	err := r.db.Pool.QueryRow(ctx, query, id).
		Scan(&request.ID, &request.PipelineRunID, &request.RequesterID, &request.ApproverID, &request.Status, &request.CreatedAt, &request.UpdatedAt, &request.Comment)
	if err == pgx.ErrNoRows {
		return AuthorizationRequest{}, nil
	}
	if err != nil {
		r.logger.Error("Failed to get authorization request", zap.String("id", id), zap.Error(err))
		return AuthorizationRequest{}, err
	}

	return request, nil
}

// UpdateAuthorizationRequest updates the status and comment of an authorization request.
func (r *PostgresRepository) UpdateAuthorizationRequest(ctx context.Context, id string, status PipelineStatus, comment string, approverID ...string) error {
	var query string
	var args []interface{}

	if len(approverID) > 0 && approverID[0] != "" {
		// Update with approver_id
		query = `UPDATE authorization_requests SET status = $1, comment = $2, approver_id = $3, updated_at = $4 WHERE id = $5`
		args = []interface{}{status, comment, approverID[0], time.Now(), id}
	} else {
		// Update without approver_id (original behavior)
		query = `UPDATE authorization_requests SET status = $1, comment = $2, updated_at = $3 WHERE id = $4`
		args = []interface{}{status, comment, time.Now(), id}
	}

	_, err := r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to update authorization request", zap.String("id", id), zap.String("status", string(status)), zap.Error(err))
		return err
	}
	return nil
}

// ListAuthorizationRequestsByPipelineRun lists authorization requests for a pipeline run.
func (r *PostgresRepository) ListAuthorizationRequestsByPipelineRun(ctx context.Context, pipelineRunID string) ([]AuthorizationRequest, error) {
	query := `SELECT id, pipeline_run_id, requester_id, approver_id, status, created_at, updated_at, comment
		FROM authorization_requests WHERE pipeline_run_id = $1`
	rows, err := r.db.Pool.Query(ctx, query, pipelineRunID)
	if err != nil {
		r.logger.Error("Failed to list authorization requests", zap.String("pipeline_run_id", pipelineRunID), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var requests []AuthorizationRequest
	for rows.Next() {
		var a AuthorizationRequest
		if err := rows.Scan(&a.ID, &a.PipelineRunID, &a.RequesterID, &a.ApproverID, &a.Status, &a.CreatedAt, &a.UpdatedAt, &a.Comment); err != nil {
			r.logger.Error("Failed to scan authorization request", zap.Error(err))
			return nil, err
		}
		requests = append(requests, a)
	}
	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating authorization requests", zap.Error(err))
		return nil, err
	}
	return requests, nil
}

// CreateExecutionHistory creates a new execution history entry.
func (r *PostgresRepository) CreateExecutionHistory(ctx context.Context, history ExecutionHistory) (ExecutionHistory, error) {
	if history.ID == "" {
		history.ID = uuid.New().String()
	}
	history.StartedAt = time.Now()

	query := `INSERT INTO execution_history (id, pipeline_run_id, requester_id, approver_id, status, started_at, completed_at, execution_time, error_message)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, pipeline_run_id, requester_id, approver_id, status, started_at, completed_at, execution_time, error_message`
	var createdHistory ExecutionHistory
	err := r.db.Pool.QueryRow(ctx, query, history.ID, history.PipelineRunID, history.RequesterID, history.ApproverID, history.Status, history.StartedAt, history.CompletedAt, history.ExecutionTime, history.ErrorMessage).
		Scan(&createdHistory.ID, &createdHistory.PipelineRunID, &createdHistory.RequesterID, &createdHistory.ApproverID, &createdHistory.Status, &createdHistory.StartedAt, &createdHistory.CompletedAt, &createdHistory.ExecutionTime, &createdHistory.ErrorMessage)
	if err != nil {
		r.logger.Error("Failed to create execution history", zap.Error(err))
		return ExecutionHistory{}, err
	}
	return createdHistory, nil
}

// GetExecutionHistory retrieves an execution history entry by its ID.
func (r *PostgresRepository) GetExecutionHistory(ctx context.Context, id string) (ExecutionHistory, error) {
	query := `SELECT id, pipeline_run_id, requester_id, approver_id, status, started_at, completed_at, execution_time, error_message
		FROM execution_history WHERE id = $1`
	var history ExecutionHistory
	err := r.db.Pool.QueryRow(ctx, query, id).
		Scan(&history.ID, &history.PipelineRunID, &history.RequesterID, &history.ApproverID, &history.Status, &history.StartedAt, &history.CompletedAt, &history.ExecutionTime, &history.ErrorMessage)
	if err == pgx.ErrNoRows {
		return ExecutionHistory{}, nil
	}
	if err != nil {
		r.logger.Error("Failed to get execution history", zap.String("id", id), zap.Error(err))
		return ExecutionHistory{}, err
	}
	return history, nil
}

// ListExecutionHistoryByPipelineRun lists execution history for a pipeline run.
func (r *PostgresRepository) ListExecutionHistoryByPipelineRun(ctx context.Context, pipelineRunID string) ([]ExecutionHistory, error) {
	query := `SELECT id, pipeline_run_id, requester_id, approver_id, status, started_at, completed_at, execution_time, error_message
		FROM execution_history WHERE pipeline_run_id = $1`
	rows, err := r.db.Pool.Query(ctx, query, pipelineRunID)
	if err != nil {
		r.logger.Error("Failed to list execution history", zap.String("pipeline_run_id", pipelineRunID), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var histories []ExecutionHistory
	for rows.Next() {
		var h ExecutionHistory
		if err := rows.Scan(&h.ID, &h.PipelineRunID, &h.RequesterID, &h.ApproverID, &h.Status, &h.StartedAt, &h.CompletedAt, &h.ExecutionTime, &h.ErrorMessage); err != nil {
			r.logger.Error("Failed to scan execution history", zap.Error(err))
			return nil, err
		}
		histories = append(histories, h)
	}
	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating execution history", zap.Error(err))
		return nil, err
	}
	return histories, nil
}

// UpdateExecutionHistoryError updates the error message of an execution history entry.
func (r *PostgresRepository) UpdateExecutionHistoryError(ctx context.Context, id, errorMessage string) error {
	query := `UPDATE execution_history SET error_message = $1, completed_at = $2, updated_at = $3 WHERE id = $4`
	_, err := r.db.Pool.Exec(ctx, query, errorMessage, time.Now(), time.Now(), id)
	if err != nil {
		r.logger.Error("Failed to update execution history error", zap.String("id", id), zap.Error(err))
		return err
	}
	return nil
}
func (r *PostgresRepository) ListPipelineUnits(ctx context.Context) ([]PipelineUnit, error) {
	query := `SELECT id, macro_service_id, created_at, updated_at FROM pipeline_units`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query pipeline units: %w", err)
	}
	defer rows.Close()

	var units []PipelineUnit
	for rows.Next() {
		var unit PipelineUnit
		if err := rows.Scan(&unit.ID, &unit.MacroServiceID, &unit.CreatedAt, &unit.UpdatedAt); err != nil {
			return nil, err
		}
		units = append(units, unit)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return units, nil
}

// UpdateExecutionHistory updates an execution history entry with all relevant fields.
func (r *PostgresRepository) UpdateExecutionHistory(ctx context.Context, history ExecutionHistory) error {
	fmt.Println("\n\n\n\n This is the history given \n\n\n %v \n\n\n\n\n", history)

	query := `
		UPDATE execution_history 
		SET status = $1, completed_at = $2, execution_time = $3, error_message = $4
		WHERE id = $5
	`

	var completedAt interface{}
	if !history.CompletedAt.IsZero() {
		completedAt = history.CompletedAt
	} else {
		completedAt = nil
	}

	var executionTime interface{}
	if history.ExecutionTime != 0 {
		executionTime = history.ExecutionTime.Milliseconds() // ✅ correct type
	} else {
		executionTime = nil
	}

	var errorMessage interface{}
	if history.ErrorMessage != "" {
		errorMessage = history.ErrorMessage
	} else {
		errorMessage = nil
	}

	_, err := r.db.Pool.Exec(ctx, query,
		history.Status,
		completedAt,
		executionTime,
		errorMessage,
		history.ID,
	)
	if err != nil {
		r.logger.Error("Failed to update execution history",
			zap.String("id", history.ID),
			zap.String("status", string(history.Status)),
			zap.Error(err))
		return err
	}

	r.logger.Info("Updated execution history",
		zap.String("id", history.ID),
		zap.String("status", string(history.Status)),
		zap.Any("completed_at", completedAt),
		zap.Any("execution_time", executionTime),
		zap.Any("error_message", errorMessage))

	return nil
}

// ListPipelineRunsWithServices retrieves pipeline runs with service information, optionally filtered by run ID
func (r *PostgresRepository) ListPipelineRunsWithServices(ctx context.Context, runID string) ([]PipelineRunStatus, error) {
	var query string
	var args []interface{}

	if runID != "" {
		// Filter by specific run ID
		query = `
		SELECT
			pr.id,
			pr.pipeline_unit_id,
			pr.status,
			pr.created_at,
			pr.updated_at,
			pr.gitlab_pipeline_id,
			COALESCE(s_macro.name, '') AS macro_service_name,
			COALESCE(u1.username, '') AS requester_name,
			COALESCE(u2.username, '') AS approver_name,
			COALESCE((
				SELECT array_agg(s_micro.name ORDER BY pd.order_index)
				FROM unnest(pr.selected_micro_service_ids) AS micro_id
				JOIN services s_micro ON s_micro.id = micro_id
				JOIN pipeline_dependencies pd ON pd.micro_service_id = micro_id AND pd.pipeline_unit_id = pr.pipeline_unit_id
			), '{}') AS micro_service_names
		FROM pipeline_runs pr
		JOIN pipeline_units pu ON pr.pipeline_unit_id = pu.id
		LEFT JOIN services s_macro ON pu.macro_service_id = s_macro.id
		LEFT JOIN authorization_requests ar ON pr.id = ar.pipeline_run_id
		LEFT JOIN users u1 ON ar.requester_id::uuid = u1.id
		LEFT JOIN users u2 ON ar.approver_id::uuid = u2.id
		WHERE pr.id = $1
		ORDER BY pr.created_at DESC`
		args = []interface{}{runID}
	} else {
		// Get all pipeline runs
		query = `
		SELECT
			pr.id,
			pr.pipeline_unit_id,
			pr.status,
			pr.created_at,
			pr.updated_at,
			COALESCE(pr.gitlab_pipeline_id, 0) AS gitlab_pipeline_id,
			COALESCE(s_macro.name, '') AS macro_service_name,
			COALESCE(u1.username, '') AS requester_name,
			COALESCE(u2.username, '') AS approver_name,
			COALESCE((
				SELECT array_agg(s_micro.name ORDER BY pd.order_index)
				FROM unnest(pr.selected_micro_service_ids) AS micro_id
				JOIN services s_micro ON s_micro.id = micro_id
				JOIN pipeline_dependencies pd ON pd.micro_service_id = micro_id AND pd.pipeline_unit_id = pr.pipeline_unit_id
			), '{}') AS micro_service_names
		FROM pipeline_runs pr
		JOIN pipeline_units pu ON pr.pipeline_unit_id = pu.id
		LEFT JOIN services s_macro ON pu.macro_service_id = s_macro.id
		LEFT JOIN authorization_requests ar ON pr.id = ar.pipeline_run_id
		LEFT JOIN users u1 ON ar.requester_id::uuid = u1.id
		LEFT JOIN users u2 ON ar.approver_id::uuid = u2.id
		ORDER BY pr.created_at DESC`
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to list pipeline runs with services", zap.String("run_id", runID), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var pipelineStatuses []PipelineRunStatus
	for rows.Next() {
		var ps PipelineRunStatus
		var gitlabPipelineID sql.NullInt32
		var microServiceNames []string

		if err := rows.Scan(
			&ps.ID,
			&ps.PipelineUnitID,
			&ps.Status,
			&ps.CreatedAt,
			&ps.UpdatedAt,
			&gitlabPipelineID,
			&ps.MacroServiceName,
			&ps.RequesterName,
			&ps.ApproverName,
			&microServiceNames,
		); err != nil {
			r.logger.Error("Failed to scan pipeline run status", zap.Error(err))
			return nil, err
		}

		if gitlabPipelineID.Valid {
			ps.GitLabPipelineID = int(gitlabPipelineID.Int32)
		}
		ps.MicroServiceNames = microServiceNames

		pipelineStatuses = append(pipelineStatuses, ps)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating pipeline run statuses", zap.Error(err))
		return nil, err
	}

	r.logger.Info("Retrieved pipeline runs with services",
		zap.String("run_id", runID),
		zap.Int("count", len(pipelineStatuses)))

	return pipelineStatuses, nil
}
