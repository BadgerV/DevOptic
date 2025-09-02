package gitlab

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/badgerv/monitoring-api/internal/auth"
	"github.com/badgerv/monitoring-api/internal/emailservice"
	"github.com/badgerv/monitoring-api/internal/websocket"
	"github.com/google/uuid"
	"gitlab.com/gitlab-org/api/client-go"
	"go.uber.org/zap"
)

// PipelineService handles the business logic for GitLab pipeline orchestration.
type PipelineService struct {
	repo         Repository
	gitlabClient *gitlab.Client
	emailService *emailservice.EmailService
	logger       *zap.Logger
	wsHub        *websocket.Hub
	authRepo     auth.UserRepository
}

// NewPipelineService creates a new PipelineService instance.
func NewPipelineService(repo Repository, gitlabClient *gitlab.Client, emailService *emailservice.EmailService, logger *zap.Logger, wsHub *websocket.Hub, authRepo auth.UserRepository) *PipelineService {
	return &PipelineService{
		repo:         repo,
		gitlabClient: gitlabClient,
		emailService: emailService,
		logger:       logger,
		wsHub:        wsHub,
		authRepo:     authRepo,
	}
}

// CreateService registers a new macro or micro service.
func (s *PipelineService) CreateService(ctx context.Context, gitlabRepoID, name, url string, serviceType ServiceType) (Service, error) {
	service := Service{
		GitLabRepoID: gitlabRepoID,
		Name:         name,
		URL:          url,
		Type:         serviceType,
	}
	createdService, err := s.repo.CreateService(ctx, service)
	if err != nil {
		s.logger.Error("Failed to create service", zap.String("gitlab_repo_id", gitlabRepoID), zap.Error(err))
		return Service{}, err
	}
	return createdService, nil
}

// CreatePipelineUnit creates a pipeline unit definition with a macro service and microservice dependencies.
func (s *PipelineService) CreatePipelineUnit(ctx context.Context, macroServiceID string, microServiceIDs []string) (PipelineUnit, error) {
	// Add context timeout for database operations
	dbCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Validate macro service
	if macroServiceID != "" {
		macroService, err := s.repo.GetServiceByID(dbCtx, macroServiceID)
		if err != nil || macroService.Type != MacroService {
			s.logger.Error("Invalid macro service ID", zap.String("macro_service_id", macroServiceID), zap.Error(err))
			return PipelineUnit{}, fmt.Errorf("invalid macro service ID: %w", err)
		}
	}
	// Validate micro services
	for _, id := range microServiceIDs {
		service, err := s.repo.GetServiceByID(dbCtx, id)
		if err != nil || service.Type != MicroService {
			s.logger.Error("Invalid micro service ID", zap.String("micro_service_id", id), zap.Error(err))
			return PipelineUnit{}, fmt.Errorf("invalid micro service ID: %s", id)
		}
	}

	unit := PipelineUnit{
		MacroServiceID:  macroServiceID,
		MicroServiceIDs: microServiceIDs,
	}
	createdUnit, err := s.repo.CreatePipelineUnit(dbCtx, unit)
	if err != nil {
		s.logger.Error("Failed to create pipeline unit", zap.Error(err))
		return PipelineUnit{}, err
	}

	// Add microservice dependencies with order index
	for i, microID := range microServiceIDs {
		if err := s.repo.AddMicroServiceDependency(dbCtx, createdUnit.ID, microID, i); err != nil {
			s.logger.Error("Failed to add microservice dependency", zap.String("pipeline_unit_id", createdUnit.ID), zap.String("micro_service_id", microID), zap.Error(err))
			return PipelineUnit{}, err
		}
	}

	return createdUnit, nil
}

// ListAllExecutionHistories retrieves all execution history records for all pipeline runs.
func (s *PipelineService) ListAllExecutionHistories(ctx context.Context, id string) ([]ExecutionHistory, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	histories, err := s.repo.ListAllExecutionHistories(dbCtx, id)
	if err != nil {
		s.logger.Error("Failed to list all execution histories", zap.Error(err))
		return nil, err
	}
	return histories, nil
}

// ListAllAuthorizationRequests retrieves all authorization requests.
func (s *PipelineService) ListAllAuthorizationRequests(ctx context.Context) ([]AuthorizationRequest, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	requests, err := s.repo.ListAllAuthorizationRequests(dbCtx)
	if err != nil {
		s.logger.Error("Failed to list all authorization requests", zap.Error(err))
		return nil, err
	}
	return requests, nil
}

// ListPipelineUnitsWithServices retrieves all pipeline units with their macro and microservices.
func (s *PipelineService) ListPipelineUnitsWithServices(ctx context.Context) ([]PipelineUnit, []Service, [][]Service, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	units, err := s.repo.ListPipelineUnits(dbCtx)
	if err != nil {
		s.logger.Error("Failed to list pipeline units", zap.Error(err))
		return nil, nil, nil, err
	}

	var macroServices []Service
	var microServicesList [][]Service

	for _, unit := range units {
		_, macro, micros, err := s.repo.GetPipelineUnitWithServices(dbCtx, unit.ID)
		if err != nil {
			s.logger.Error("Failed to get services for pipeline unit", zap.String("pipeline_unit_id", unit.ID), zap.Error(err))
			return nil, nil, nil, err
		}

		unit.MicroServiceIDs, err = s.repo.GetMicroServiceDependencies(dbCtx, unit.ID)
		if err != nil {
			s.logger.Error("Failed to get microservice dependencies", zap.String("pipeline_unit_id", unit.ID), zap.Error(err))
			return nil, nil, nil, err
		}

		macroServices = append(macroServices, macro)
		microServicesList = append(microServicesList, micros)
	}

	return units, macroServices, microServicesList, nil
}

// GetPipelineUnitWithServices retrieves a pipeline unit with its macro and microservices.
func (s *PipelineService) GetPipelineUnitWithServices(ctx context.Context, id string) (PipelineUnit, Service, []Service, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return s.repo.GetPipelineUnitWithServices(dbCtx, id)
}

// TriggerPipelineUnit triggers an execution of a pipeline unit, initiating the approval process.
func (s *PipelineService) TriggerPipelineUnit(ctx context.Context, pipelineUnitID, requesterID string, selectedMicroServiceIDs []string) (PipelineRun, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	unit, err := s.repo.GetPipelineUnit(dbCtx, pipelineUnitID)
	if err != nil {
		s.logger.Error("Failed to get pipeline unit", zap.String("pipeline_unit_id", pipelineUnitID), zap.Error(err))
		return PipelineRun{}, err
	}

	// Validate selected microservices are a subset of the unit's microservices
	validMicroIDs := make(map[string]bool)
	for _, id := range unit.MicroServiceIDs {
		validMicroIDs[id] = true
	}
	for _, id := range selectedMicroServiceIDs {
		if !validMicroIDs[id] {
			s.logger.Error("Invalid selected microservice ID", zap.String("micro_service_id", id), zap.String("pipeline_unit_id", pipelineUnitID))
			return PipelineRun{}, fmt.Errorf("invalid selected microservice ID: %s", id)
		}
	}

	run := PipelineRun{
		PipelineUnitID:          pipelineUnitID,
		Status:                  StatusPending,
		SelectedMicroServiceIDs: selectedMicroServiceIDs,
	}

	createdRun, err := s.repo.CreatePipelineRun(dbCtx, run)
	if err != nil {
		s.logger.Error("Failed to create pipeline run", zap.String("pipeline_unit_id", pipelineUnitID), zap.Error(err))
		return PipelineRun{}, err
	}

	// Create authorization request for the run
	authRequest := AuthorizationRequest{
		PipelineRunID: createdRun.ID,
		RequesterID:   requesterID,
		ApproverID:    nil,
		Status:        StatusPending,
	}

	req, err := s.repo.CreateAuthorizationRequest(dbCtx, authRequest)
	if err != nil {
		s.logger.Error("Failed to create authorization request",
			zap.String("pipeline_run_id", createdRun.ID),
			zap.Error(err),
		)
		return PipelineRun{}, err
	}

	fullAuthRequest, err := s.repo.GetAuthorizationRequestByID(ctx, req.ID)

	go func(f *AuthorizationRequest) {
		ctx := context.Background()
		userID, _ := uuid.Parse(authRequest.RequesterID)
		userDeliveryEmail, _ := s.authRepo.GetDeliveryEmail(ctx, userID)

		fmt.Println("\n\n\n - delivery email >>>>>>>>>>>>>>", userDeliveryEmail, "delivery \n\n\n\n")

		htmlDoc, _ := s.RenderAuthorizationRequestToHTML(fullAuthRequest)
		if err := s.emailService.SendHTML(
			"Pipeline Triggered", htmlDoc, []string{userDeliveryEmail},
		); err != nil {
			s.logger.Error("Failed to send email notification",
				zap.String("pipeline_run_id", userDeliveryEmail),
				zap.Error(err))
		}
	}(fullAuthRequest)

	// Broadcast WebSocket update
	if err != nil {
		s.logger.Error("Failed to marshal WebSocket message", zap.String("pipeline_run_id", createdRun.ID), zap.Error(err))
	} else {
		s.broadcastPipelineStatusChange(ctx, createdRun.ID, StatusPending, "Pipeline execution triggered, awaiting approval")
	}

	return createdRun, nil
}

// ApprovePipelineRun approves a pipeline run and triggers execution in the background.
func (s *PipelineService) ApprovePipelineRun(ctx context.Context, authRequestID, approverID string, comment string) error {
	// Use a separate context with timeout for database operations
	dbCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	authRequest, err := s.repo.GetAuthorizationRequest(dbCtx, authRequestID)
	if err != nil {
		s.logger.Error("Failed to get authorization request", zap.String("auth_request_id", authRequestID), zap.Error(err))
		return err
	}
	if authRequest.Status != StatusPending {
		s.logger.Info("Fetched authorization request", zap.Any("auth_request", authRequest))
		s.logger.Error("Authorization request not pending", zap.String("auth_request_id", authRequestID), zap.String("status", string(authRequest.Status)))
		return fmt.Errorf("authorization request is not pending")
	}

	// Update authorization request
	if err := s.repo.UpdateAuthorizationRequest(dbCtx, authRequestID, StatusAccepted, comment, approverID); err != nil {
		s.logger.Error("Failed to update authorization request", zap.String("auth_request_id", authRequestID), zap.Error(err))
		return err
	}

	// Update pipeline run
	run, err := s.repo.GetPipelineRun(dbCtx, authRequest.PipelineRunID)
	if err != nil {
		s.logger.Error("Failed to get pipeline run", zap.String("pipeline_run_id", authRequest.PipelineRunID), zap.Error(err))
		return err
	}
	if err := s.repo.UpdatePipelineRunStatus(dbCtx, run.ID, StatusAccepted); err != nil {
		s.logger.Error("Failed to update pipeline run status", zap.String("pipeline_run_id", run.ID), zap.Error(err))
		return err
	}

	approverUUID, err := uuid.Parse(approverID)
	if err != nil {
		return err // or handle the error properly
	}

	// Create execution history for the run
	history := ExecutionHistory{
		PipelineRunID: run.ID,
		RequesterID:   authRequest.RequesterID,
		ApproverID:    approverUUID,
		Status:        StatusRunning,
		StartedAt:     time.Now(),
	}
	history, err = s.repo.CreateExecutionHistory(dbCtx, history)
	if err != nil {
		s.logger.Error("Failed to create execution history", zap.String("pipeline_run_id", run.ID), zap.Error(err))
		return err
	}

	// Broadcast WebSocket update
	message := struct {
		PipelineRunID string         `json:"pipeline_run_id"`
		Status        PipelineStatus `json:"status"`
		Message       string         `json:"message"`
		Timestamp     time.Time      `json:"timestamp"`
		ApproverID    string         `json:"approver_id"`
	}{
		PipelineRunID: run.ID,
		Status:        StatusAccepted,
		Message:       "Pipeline run approved",
		Timestamp:     time.Now(),
		ApproverID:    approverID,
	}
	payload, err := json.Marshal(message)
	if err != nil {
		s.logger.Error("Failed to marshal WebSocket message", zap.String("pipeline_run_id", run.ID), zap.Error(err))
	} else {
		s.logger.Info("Broadcasting pipeline approval", zap.String("pipeline_run_id", run.ID), zap.String("payload", string(payload)))
		s.broadcastPipelineStatusChange(ctx, run.ID, StatusAccepted, "Pipeline run approved")
	}

	// Fetch the associated pipeline unit for execution
	unit, err := s.repo.GetPipelineUnit(dbCtx, run.PipelineUnitID)
	if err != nil {
		s.logger.Error("Failed to get pipeline unit", zap.String("pipeline_unit_id", run.PipelineUnitID), zap.Error(err))
		s.repo.UpdateExecutionHistoryError(dbCtx, history.ID, err.Error())
		s.repo.UpdatePipelineRunStatus(dbCtx, run.ID, StatusRejected)
		return err
	}

	// CRITICAL FIX: Create a new background context that won't be cancelled when HTTP request ends
	// This prevents "context canceled" errors in the background goroutine
	backgroundCtx := context.Background()

	userID, err := uuid.Parse(authRequest.RequesterID)

	if err != nil {
		s.logger.Error("\n\nFailed to create pipeline run invalid userID - ", zap.String("auth_request_id", authRequestID), zap.Error(err))
	}
	userDeliveryEmail, err := s.authRepo.GetDeliveryEmail(ctx, userID)

	if err != nil {
		s.logger.Error("\n\nFailed to create pipeline run - Unable to get user delivery email address - ", zap.String("auth_request_id", authRequestID), zap.Error(err))
	}

	fullAuthRequest, _ := s.repo.GetAuthorizationRequestByID(ctx, authRequest.ID)

	htlmDoc, _ := s.RenderAuthorizationRequestToHTML(fullAuthRequest)

	if err != nil {
		s.logger.Error("\n\nFailed to create pipeline run - Unable to parse HTML document - ", zap.String("auth_request_id", authRequestID), zap.Error(err))
	}

	if err := s.emailService.SendHTML(
		"Pipeline Has Been Approved", htlmDoc, []string{userDeliveryEmail},
	); err != nil {
		s.logger.Error("Failed to send email notification",
			zap.String("pipeline_run_id", userDeliveryEmail),
			zap.Error(err))
	}

	// Start pipeline execution in the background with a fresh context
	go func() {
		// Create a context with timeout for the entire pipeline execution
		execCtx, execCancel := context.WithTimeout(backgroundCtx, 1*time.Hour)
		defer execCancel()

		s.broadcastPipelineStatusChange(execCtx, run.ID, StatusRunning, "Pipeline run approved and execution started")

		if err := s.executePipelineChain(execCtx, &run, &unit, &history); err != nil {
			s.logger.Error("Pipeline execution failed", zap.String("pipeline_run_id", run.ID), zap.Error(err))
		}
	}()
	return nil
}

// RejectPipelineRun rejects a pipeline run.
func (s *PipelineService) RejectPipelineRun(ctx context.Context, authRequestID, approverID, comment string) error {
	dbCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	fmt.Println(approverID)

	authRequest, err := s.repo.GetAuthorizationRequest(dbCtx, authRequestID)
	if err != nil {
		s.logger.Error("Failed to get authorization request", zap.String("auth_request_id", authRequestID), zap.Error(err))
		return err
	}
	if authRequest.Status != StatusPending {
		s.logger.Error("Authorization request not pending", zap.String("auth_request_id", authRequestID), zap.String("status", string(authRequest.Status)))
		return fmt.Errorf("authorization request is not pending")
	}

	// Update authorization request
	if err := s.repo.UpdateAuthorizationRequest(dbCtx, authRequestID, StatusRejected, comment); err != nil {
		s.logger.Error("Failed to update authorization request", zap.String("auth_request_id", authRequestID), zap.Error(err))
		return err
	}

	// Update pipeline run status
	if err := s.repo.UpdatePipelineRunStatus(dbCtx, authRequest.PipelineRunID, StatusRejected); err != nil {
		s.logger.Error("Failed to update pipeline run status", zap.String("pipeline_run_id", authRequest.PipelineRunID), zap.Error(err))
		return err
	}

	approverUUID, err := uuid.Parse(approverID)
	if err != nil {
		return err // or handle the error properly
	}

	// Create execution history
	history := ExecutionHistory{
		PipelineRunID: authRequest.PipelineRunID,
		RequesterID:   authRequest.RequesterID,
		ApproverID:    approverUUID,
		Status:        StatusRejected,
		StartedAt:     time.Now(),
		CompletedAt:   time.Now(),
		ErrorMessage:  comment,
	}
	if _, err := s.repo.CreateExecutionHistory(dbCtx, history); err != nil {
		s.logger.Error("Failed to create execution history", zap.String("pipeline_run_id", authRequest.PipelineRunID), zap.Error(err))
		return err
	}

	// Broadcast WebSocket update
	message := struct {
		PipelineRunID string         `json:"pipeline_run_id"`
		Status        PipelineStatus `json:"status"`
		Message       string         `json:"message"`
		Timestamp     time.Time      `json:"timestamp"`
		ApproverID    string         `json:"approver_id"`
		Comment       string         `json:"comment,omitempty"`
	}{
		PipelineRunID: authRequest.PipelineRunID,
		Status:        StatusRejected,
		Message:       fmt.Sprintf("Pipeline run rejected: %s", comment),
		Timestamp:     time.Now(),
		ApproverID:    approverID,
		Comment:       comment,
	}
	payload, err := json.Marshal(message)
	if err != nil {
		s.logger.Error("Failed to marshal WebSocket message", zap.String("pipeline_run_id", authRequest.PipelineRunID), zap.Error(err))
	} else {
		s.logger.Info("Broadcasting pipeline rejection", zap.String("pipeline_run_id", authRequest.PipelineRunID), zap.String("payload", string(payload)))
		s.broadcastPipelineStatusChange(ctx, authRequest.PipelineRunID, StatusRejected, fmt.Sprintf("Pipeline run rejected: %s", comment))
	}

	userID, err := uuid.Parse(authRequest.RequesterID)

	if err != nil {
		s.logger.Error("\n\nFailed to create pipeline run invalid userID - ", zap.String("auth_request_id", authRequestID), zap.Error(err))
	}
	userDeliveryEmail, err := s.authRepo.GetDeliveryEmail(ctx, userID)

	if err != nil {
		s.logger.Error("\n\nFailed to create pipeline run - Unable to get user delivery email address - ", zap.String("auth_request_id", authRequestID), zap.Error(err))
	}

	fullAuthRequest, _ := s.repo.GetAuthorizationRequestByID(ctx, authRequest.ID)

	htlmDoc, _ := s.RenderAuthorizationRequestToHTML(fullAuthRequest)

	if err != nil {
		s.logger.Error("\n\nFailed to create pipeline run - Unable to parse HTML document - ", zap.String("auth_request_id", authRequestID), zap.Error(err))
	}

	go func(userDeliveryEmail, htlmDoc string) {
		if err := s.emailService.SendHTML(
			"Pipeline Has Been Rejected", htlmDoc, []string{userDeliveryEmail},
		); err != nil {
			s.logger.Error("Failed to send email notification",
				zap.String("pipeline_run_id", userDeliveryEmail),
				zap.Error(err))
		}
	}(userDeliveryEmail, htlmDoc)

	return nil
}

// GetPipelineRunStatus retrieves the current status of a pipeline run.
func (s *PipelineService) GetPipelineRunStatus(ctx context.Context, pipelineRunID string) (PipelineStatus, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	run, err := s.repo.GetPipelineRun(dbCtx, pipelineRunID)
	if err != nil {
		s.logger.Error("Failed to get pipeline run", zap.String("pipeline_run_id", pipelineRunID), zap.Error(err))
		return "", err
	}
	return run.Status, nil
}

// ListExecutionHistory retrieves the execution history for a pipeline run.
func (s *PipelineService) ListExecutionHistory(ctx context.Context, pipelineRunID string) ([]ExecutionHistory, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	history, err := s.repo.ListExecutionHistoryByPipelineRun(dbCtx, pipelineRunID)
	if err != nil {
		s.logger.Error("Failed to list execution history", zap.String("pipeline_run_id", pipelineRunID), zap.Error(err))
		return nil, err
	}
	return history, nil
}

// ListAllServices retrieves all services, categorized into microservices and one macroservice.
func (s *PipelineService) ListAllServices(ctx context.Context) ([]Service, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Fetch all services
	services, err := s.repo.ListServicesByType(dbCtx, "") // Empty type to get all services
	if err != nil {
		s.logger.Error("Failed to list all services", zap.Error(err))
		return nil, err
	}

	var microServices []Service
	var macroService *Service
	for _, svc := range services {
		if svc.Type == MicroService {
			microServices = append(microServices, svc)
		} else if svc.Type == MacroService && macroService == nil {
			microServices = append(microServices, svc)
		}
	}

	return microServices, nil
}

// GetServiceByID retrieves a service by ID - FIXED: removed hardcoded ID
func (s *PipelineService) GetServiceByID(ctx context.Context, serviceID string) (Service, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	return s.repo.GetServiceByID(dbCtx, serviceID)
}

// executePipelineChain triggers and monitors pipelines for microservices and the macroservice in sequence.
func (s *PipelineService) executePipelineChain(ctx context.Context, run *PipelineRun, unit *PipelineUnit, history *ExecutionHistory) error {
	// Fetch microservices with retry
	ids := append(run.SelectedMicroServiceIDs, unit.MacroServiceID)
	microServices, err := s.fetchMicroServicesWithRetry(ctx, ids, 3)

	if err != nil {
		s.logger.Error("Failed to fetch microservices after retries", zap.String("pipeline_run_id", run.ID), zap.Error(err))
		return s.handlePipelineError(ctx, run, history, "", fmt.Sprintf("Failed to fetch microservices: %v", err))
	}

	// Process each microservice
	for i, microService := range microServices {
		// Broadcast pipeline start
		s.broadcastPipelineStatusChange(ctx, run.ID, StatusRunning, fmt.Sprintf("Starting pipeline for microservice %s", microService.Name))

		// Trigger microservice pipeline
		ref := "development"
		refPtr := &ref // Convert string to *string

		// Add pipeline variables only if this is the last microservice
		var variables []*gitlab.PipelineVariableOptions
		if i == len(microServices)-1 {
			deployEnvKey := "DEPLOY_ENV" // Create string variable
			deployEnvValue := "QA"       // Create string variable
			variables = []*gitlab.PipelineVariableOptions{
				{
					Key:   &deployEnvKey,   // Use *string
					Value: &deployEnvValue, // Use *string
				},
			}
		}
		variablesPtr := &variables // Convert slice to *slice

		pipeline, _, err := s.gitlabClient.Pipelines.CreatePipeline(
			microService.GitLabRepoID,
			&gitlab.CreatePipelineOptions{
				Ref:       refPtr,       // Use *string
				Variables: variablesPtr, // Use *[]*gitlab.PipelineVariableOptions
			},
		)
		if err != nil {
			s.logger.Error("Failed to trigger microservice pipeline",
				zap.String("pipeline_run_id", run.ID),
				zap.String("micro_service_id", microService.ID),
				zap.Error(err),
			)
			return s.handlePipelineError(
				ctx, run, history, microService.ID,
				fmt.Sprintf("Microservice %s pipeline failed: %v", microService.Name, err),
			)
		}

		// Update pipeline run with GitLab pipeline ID
		updateRunCtx, updateRunCancel := context.WithTimeout(ctx, 30*time.Second)
		defer updateRunCancel()
		run.GitLabPipelineID = pipeline.ID
		if err := s.repo.UpdatePipelineRun(updateRunCtx, *run); err != nil {
			s.logger.Error("Failed to update pipeline run with GitLab ID",
				zap.String("pipeline_run_id", run.ID),
				zap.Error(err),
			)
			return s.handlePipelineError(
				ctx, run, history, microService.ID,
				fmt.Sprintf("Failed to update pipeline run: %v", err),
			)
		}

		// Poll for microservice pipeline status
		if err := s.pollPipelineStatus(ctx, microService.GitLabRepoID, pipeline.ID); err != nil {
			s.logger.Error("Microservice pipeline failed or timed out",
				zap.String("pipeline_run_id", run.ID),
				zap.Int("gitlab_pipeline_id", pipeline.ID),
				zap.Error(err),
			)
			return s.handlePipelineError(
				ctx, run, history, microService.ID,
				fmt.Sprintf("Microservice %s pipeline failed: %v", microService.Name, err),
			)
		}
	}

	// Update to completed
	finalUpdateCtx, finalUpdateCancel := context.WithTimeout(ctx, 30*time.Second)
	defer finalUpdateCancel()

	history.Status = StatusCompleted
	history.CompletedAt = time.Now()
	// history.ExecutionTime = history.CompletedAt.Sub(history.StartedAt)
	ms := time.Since(history.StartedAt).Milliseconds() // int64
	executionTime := time.Duration(ms) * time.Millisecond

	history.ExecutionTime = executionTime
	if err := s.repo.UpdateExecutionHistory(finalUpdateCtx, *history); err != nil {
		s.logger.Error("Failed to update execution history",
			zap.String("history_id", history.ID),
			zap.Error(err),
		)
		// Log but don't fail the pipeline since the execution completed
	}

	if err := s.repo.UpdatePipelineRunStatus(finalUpdateCtx, run.ID, StatusCompleted); err != nil {
		s.logger.Error("Failed to update pipeline run to completed",
			zap.String("pipeline_run_id", run.ID),
			zap.Error(err),
		)
		s.broadcastPipelineStatusChange(ctx, run.ID, StatusRejected, fmt.Sprintf("Failed to update pipeline run status: %v", err))
		return err
	}

	go func(history *ExecutionHistory) {
		ctx := context.Background()
		historyFromID, err := s.repo.GetExecutionHistoryByID(ctx, history.RequesterID, history.ID)
		if err != nil {
			s.logger.Error("Failed to send email",
				zap.String("pipeline_run_id", run.ID),
				zap.Error(err),
			)
		}

		htmlDoc, _ := s.RenderExecutionHistoryToHTML(historyFromID)

		userID, _ := uuid.Parse(historyFromID.RequesterID)

		userDeliveryEmail, _ := s.authRepo.GetDeliveryEmail(ctx, userID)

		if err := s.emailService.SendHTML(
			"Pipeline Has Run And Completed Successful", htmlDoc, []string{userDeliveryEmail},
		); err != nil {
			s.logger.Error("Failed to send email notification",
				zap.String("pipeline_run_id", userDeliveryEmail),
				zap.Error(err))
		}
	}(history)

	// Broadcast success
	s.broadcastPipelineStatusChange(ctx, run.ID, StatusCompleted, "Pipeline run completed successfully")

	return nil
}

// handlePipelineError centralizes error handling for pipeline failures
func (s *PipelineService) handlePipelineError(ctx context.Context, run *PipelineRun, history *ExecutionHistory, microServiceID, errorMessage string) error {
	updateCtx, updateCancel := context.WithTimeout(ctx, 30*time.Second)
	defer updateCancel()

	history.Status = StatusRejected
	history.CompletedAt = time.Now()
	history.ExecutionTime = history.CompletedAt.Sub(history.StartedAt)
	history.ErrorMessage = errorMessage
	if err := s.repo.UpdateExecutionHistory(updateCtx, *history); err != nil {
		s.logger.Error("Failed to update execution history",
			zap.String("history_id", history.ID),
			zap.Error(err),
		)
	}

	if err := s.repo.UpdatePipelineRunStatus(updateCtx, run.ID, StatusRejected); err != nil {
		s.logger.Error("Failed to update pipeline run status",
			zap.String("pipeline_run_id", run.ID),
			zap.Error(err),
		)
	}

	s.broadcastPipelineUpdate(ctx, run.ID, microServiceID, StatusRejected, errorMessage)
	return fmt.Errorf(errorMessage)
}

// fetchMicroServicesWithRetry fetches microservices with retry logic to handle temporary failures
func (s *PipelineService) fetchMicroServicesWithRetry(ctx context.Context, ids []string, maxRetries int) ([]Service, error) {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: wait 1s, 2s, 4s between retries
			waitTime := time.Duration(1<<(attempt-1)) * time.Second
			s.logger.Info("Retrying fetchMicroServices",
				zap.Int("attempt", attempt),
				zap.Duration("wait_time", waitTime),
				zap.Int("service_count", len(ids)))

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(waitTime):
				// Continue with retry
			}
		}

		services, err := s.fetchMicroServices(ctx, ids)
		if err == nil {
			return services, nil
		}

		lastErr = err
		s.logger.Warn("fetchMicroServices attempt failed",
			zap.Int("attempt", attempt),
			zap.Error(err))
	}

	return nil, fmt.Errorf("failed to fetch microservices after %d attempts: %w", maxRetries+1, lastErr)
}

// fetchMicroServices fetches full Service details for microservice IDs with proper timeout handling.
func (s *PipelineService) fetchMicroServices(ctx context.Context, ids []string) ([]Service, error) {
	var services []Service

	for i, id := range ids {
		// Create a timeout context for each individual service fetch
		serviceCtx, cancel := context.WithTimeout(ctx, 10*time.Second)

		svc, err := s.repo.GetServiceByID(serviceCtx, id)
		cancel() // Always cancel to free resources

		if err != nil {
			s.logger.Error("Failed to fetch microservice",
				zap.String("service_id", id),
				zap.Int("service_index", i),
				zap.Error(err))
			return nil, fmt.Errorf("failed to fetch microservice %s (index %d): %w", id, i, err)
		}

		if svc.ID == "" {
			s.logger.Error("Microservice not found", zap.String("service_id", id))
			return nil, fmt.Errorf("microservice ID %s not found", id)
		}

		services = append(services, svc)
		s.logger.Debug("Successfully fetched microservice",
			zap.String("service_id", id),
			zap.String("service_name", svc.Name))
	}

	s.logger.Info("Successfully fetched all microservices",
		zap.Int("count", len(services)),
		zap.Strings("service_ids", ids))

	return services, nil
}

// pollPipelineStatus polls GitLab for pipeline completion with better context handling.
func (s *PipelineService) pollPipelineStatus(ctx context.Context, projectID string, pipelineID int) error {
	const pollInterval = 10 * time.Second
	const maxDuration = 30 * time.Minute
	start := time.Now()

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Pipeline polling cancelled by context",
				zap.String("project_id", projectID),
				zap.Int("pipeline_id", pipelineID))
			return ctx.Err()
		case <-ticker.C:
			if time.Since(start) > maxDuration {
				s.logger.Error("Pipeline polling timed out",
					zap.String("project_id", projectID),
					zap.Int("pipeline_id", pipelineID),
					zap.Duration("duration", time.Since(start)))
				return fmt.Errorf("pipeline %d timed out after %v", pipelineID, maxDuration)
			}

			pipeline, _, err := s.gitlabClient.Pipelines.GetPipeline(projectID, pipelineID)
			if err != nil {
				s.logger.Error("Failed to get pipeline status",
					zap.String("project_id", projectID),
					zap.Int("pipeline_id", pipelineID),
					zap.Error(err))
				return fmt.Errorf("failed to get pipeline %d status: %w", pipelineID, err)
			}

			s.logger.Debug("Pipeline status check",
				zap.String("project_id", projectID),
				zap.Int("pipeline_id", pipelineID),
				zap.String("status", pipeline.Status))

			switch pipeline.Status {
			case "success":
				s.logger.Info("Pipeline completed successfully",
					zap.String("project_id", projectID),
					zap.Int("pipeline_id", pipelineID))
				return nil
			case "failed", "canceled":
				s.logger.Error("Pipeline failed",
					zap.String("project_id", projectID),
					zap.Int("pipeline_id", pipelineID),
					zap.String("status", pipeline.Status))
				return fmt.Errorf("pipeline %d failed with status %s", pipelineID, pipeline.Status)
			}
			// Continue polling for running, pending, etc.
		}
	}
}

// broadcastPipelineUpdate sends a WebSocket update with detailed pipeline status.
func (s *PipelineService) broadcastPipelineUpdate(ctx context.Context, runID, serviceID string, status PipelineStatus, message string) {
	msg := struct {
		PipelineRunID    string         `json:"pipeline_run_id"`
		CurrentServiceID string         `json:"current_service_id,omitempty"`
		Status           PipelineStatus `json:"status"`
		Message          string         `json:"message"`
		Timestamp        time.Time      `json:"timestamp"`
	}{
		PipelineRunID:    runID,
		CurrentServiceID: serviceID,
		Status:           status,
		Message:          message,
		Timestamp:        time.Now(),
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		s.logger.Error("Failed to marshal WebSocket message", zap.String("pipeline_run_id", runID), zap.String("service_id", serviceID), zap.Error(err))
		return
	}
	s.logger.Info("Broadcasting pipeline update", zap.String("pipeline_run_id", runID), zap.String("service_id", serviceID), zap.String("payload", string(payload)))
	s.broadcastPipelineStatusChange(ctx, runID, status, message)

}

func (s *PipelineService) GetAllPipelineStatuses(ctx context.Context) (PipelineStatusResponse, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	pipelineStatuses, err := s.repo.ListPipelineRunsWithServices(dbCtx, "")
	if err != nil {
		s.logger.Error("Failed to get all pipeline statuses", zap.Error(err))
		return PipelineStatusResponse{}, err
	}

	response := PipelineStatusResponse{
		Running:   []PipelineRunStatus{},
		Pending:   []PipelineRunStatus{},
		Completed: []PipelineRunStatus{},
		Failed:    []PipelineRunStatus{},
		Total:     len(pipelineStatuses),
	}

	// Group pipelines by status
	for _, pipeline := range pipelineStatuses {
		switch pipeline.Status {
		case "running":
			response.Running = append(response.Running, pipeline)
		case "pending":
			response.Pending = append(response.Pending, pipeline)
		case "completed", "accepted":
			response.Completed = append(response.Completed, pipeline)
		case "failed", "rejected":
			response.Failed = append(response.Failed, pipeline)
		default:
			response.Pending = append(response.Pending, pipeline)
		}
	}

	s.logger.Info("Retrieved pipeline statuses",
		zap.Int("total", response.Total),
		zap.Int("running", len(response.Running)),
		zap.Int("pending", len(response.Pending)),
		zap.Int("completed", len(response.Completed)),
		zap.Int("failed", len(response.Failed)))

	return response, nil
}

// GetPipelineStatusByRunID returns the status of a specific pipeline run
func (s *PipelineService) GetPipelineStatusByRunID(ctx context.Context, runID string) (PipelineRunStatus, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	pipelineStatuses, err := s.repo.ListPipelineRunsWithServices(dbCtx, runID)
	if err != nil {
		s.logger.Error("Failed to get pipeline status by run ID", zap.String("run_id", runID), zap.Error(err))
		return PipelineRunStatus{}, err
	}

	if len(pipelineStatuses) == 0 {
		return PipelineRunStatus{}, fmt.Errorf("pipeline run not found: %s", runID)
	}

	s.logger.Info("Retrieved pipeline status by run ID", zap.String("run_id", runID), zap.String("status", pipelineStatuses[0].Status))
	return pipelineStatuses[0], nil
}

func (s *PipelineService) broadcastPipelineStatusChange(ctx context.Context, runID string, status PipelineStatus, message string) {
	// Get detailed pipeline information for the broadcast
	pipelineInfo, err := s.GetPipelineStatusByRunID(ctx, runID)
	if err != nil {
		s.logger.Error("Failed to get pipeline info for broadcast", zap.String("run_id", runID), zap.Error(err))
		// Fallback to basic message
		s.broadcastBasicUpdate(runID, status, message)
		return
	}

	// Enhanced message with full pipeline context
	enhancedMessage := struct {
		Type             string            `json:"type"`
		PipelineRunID    string            `json:"pipeline_run_id"`
		PipelineUnitID   string            `json:"pipeline_unit_id"`
		Status           PipelineStatus    `json:"status"`
		Message          string            `json:"message"`
		MacroServiceName string            `json:"macro_service_name,omitempty"`
		MicroServices    []string          `json:"micro_service_names"`
		RequesterName    string            `json:"requester_name,omitempty"`
		ApproverName     string            `json:"approver_name,omitempty"`
		Timestamp        time.Time         `json:"timestamp"`
		PipelineInfo     PipelineRunStatus `json:"pipeline_info"`
	}{
		Type:             "pipeline_status_change",
		PipelineRunID:    runID,
		PipelineUnitID:   pipelineInfo.PipelineUnitID,
		Status:           status,
		Message:          message,
		MacroServiceName: pipelineInfo.MacroServiceName,
		MicroServices:    pipelineInfo.MicroServiceNames,
		RequesterName:    pipelineInfo.RequesterName,
		ApproverName:     pipelineInfo.ApproverName,
		Timestamp:        time.Now(),
		PipelineInfo:     pipelineInfo,
	}

	payload, err := json.Marshal(enhancedMessage)
	if err != nil {
		s.logger.Error("Failed to marshal enhanced WebSocket message", zap.String("pipeline_run_id", runID), zap.Error(err))
		s.broadcastBasicUpdate(runID, status, message)
		return
	}

	s.logger.Info("Broadcasting enhanced pipeline status change",
		zap.String("pipeline_run_id", runID),
		zap.String("status", string(status)),
		zap.String("message", message))

	s.wsHub.Broadcast(websocket.Message{
		Type:      "pipeline_status_change",
		ID:        runID,
		Payload:   string(payload),
		Timestamp: time.Now(),
	})
}

// broadcastBasicUpdate is a fallback for when detailed pipeline info cannot be retrieved
func (s *PipelineService) broadcastBasicUpdate(runID string, status PipelineStatus, message string) {
	basicMessage := struct {
		Type          string         `json:"type"`
		PipelineRunID string         `json:"pipeline_run_id"`
		Status        PipelineStatus `json:"status"`
		Message       string         `json:"message"`
		Timestamp     time.Time      `json:"timestamp"`
	}{
		Type:          "pipeline_status_change",
		PipelineRunID: runID,
		Status:        status,
		Message:       message,
		Timestamp:     time.Now(),
	}

	payload, err := json.Marshal(basicMessage)
	if err != nil {
		s.logger.Error("Failed to marshal basic WebSocket message", zap.String("pipeline_run_id", runID), zap.Error(err))
		return
	}

	s.wsHub.Broadcast(websocket.Message{
		Type:      "pipeline_status_change",
		ID:        runID,
		Payload:   string(payload),
		Timestamp: time.Now(),
	})
}

// parseProjectID converts a macro service ID to a GitLab project ID (assumes string ID for simplicity).
// func (s *PipelineService) parseProjectID(macroServiceID string) (string, error) {
// if macroServiceID == "" {
// 	s.logger.Error("Macro service ID is empty")
// 	return "", fmt.Errorf("macro service ID is empty")
// }
// return macroServiceID, nil.CompletedAt.Sub(history.StartedAt)
// 	history.ErrorMessage = err.Error()
// 	if updateErr := s.repo.UpdateExecutionHistory(updateCtx, *history); updateErr != nil {
// 		s.logger.Error("Failed to update execution history", zap.String("history_id", history.ID), zap.Error(updateErr))
// 	}
// 	s.repo.UpdatePipelineRunStatus(updateCtx, run.ID, StatusRejected)
// 	s.broadcastPipelineUpdate(ctx, run.ID, "", StatusRejected, fmt.Sprintf("Failed to fetch microservices: %v", err))

// 	return err
// }

// // Trigger and monitor microservices sequentially
// for _, micro := range microServices {
// 	s.broadcastPipelineUpdate(ctx, run.ID, micro.ID, StatusRunning, fmt.Sprintf("Starting pipeline for microservice %s", micro.Name))
// 	ref := "development"
// 	pipeline, _, err := s.gitlabClient.Pipelines.CreatePipeline(
// 		micro.GitLabRepoID,
// 		&gitlab.CreatePipelineOptions{
// 			Ref: &ref,
// 		},
// 	)
// 	if err != nil {
// 		s.logger.Error("Failed to trigger microservice pipeline",
// 			zap.String("pipeline_run_id", run.ID),
// 			zap.String("micro_service_id", micro.ID),
// 			zap.Error(err),
// 		)

// 		// Update with proper context
// 		updateCtx, updateCancel := context.WithTimeout(context.Background(), 30*time.Second)
// 		defer updateCancel()

// 		history.Status = StatusRejected
// 		history.CompletedAt = time.Now()
// 		history.ExecutionTime = history.CompletedAt.Sub(history.StartedAt)
// 		history.ErrorMessage = fmt.Sprintf("Microservice %s failed: %v", micro.ID, err)
// 		if updateErr := s.repo.UpdateExecutionHistory(updateCtx, *history); updateErr != nil {
// 			s.logger.Error("Failed to update execution history", zap.String("history_id", history.ID), zap.Error(updateErr))
// 		}
// 		s.repo.UpdatePipelineRunStatus(updateCtx, run.ID, StatusRejected)
// 		s.broadcastPipelineUpdate(ctx, run.ID, micro.ID, StatusRejected, fmt.Sprintf("Microservice %s pipeline failed: %v", micro.Name, err))
// 		return err
// 	}

// 	// Poll for pipeline status
// 	if err := s.pollPipelineStatus(ctx, micro.GitLabRepoID, pipeline.ID); err != nil {
// 		s.logger.Error("Microservice pipeline failed or timed out",
// 			zap.String("pipeline_run_id", run.ID),
// 			zap.String("micro_service_id", micro.ID),
// 			zap.Int("gitlab_pipeline_id", pipeline.ID),
// 			zap.Error(err),
// 		)

// 		updateCtx, updateCancel := context.WithTimeout(context.Background(), 30*time.Second)
// 		defer updateCancel()

// 		history.Status = StatusRejected
// 		history.CompletedAt = time.Now()
// 		history.ExecutionTime = history.CompletedAt.Sub(history.StartedAt)
// 		history.ErrorMessage = fmt.Sprintf("Microservice %s pipeline failed: %v", micro.ID, err)
// 		if updateErr := s.repo.UpdateExecutionHistory(updateCtx, *history); updateErr != nil {
// 			s.logger.Error("Failed to update execution history", zap.String("history_id", history.ID), zap.Error(updateErr))
// 		}
// 		s.repo.UpdatePipelineRunStatus(updateCtx, run.ID, StatusRejected)
// 		s.broadcastPipelineUpdate(ctx, run.ID, micro.ID, StatusRejected, fmt.Sprintf("Microservice %s pipeline failed: %v", micro.Name, err))
// 		return err
// 	}
// 	s.broadcastPipelineUpdate(ctx, run.ID, micro.ID, StatusCompleted, fmt.Sprintf("Microservice %s pipeline completed successfully", micro.Name))
// }

// // Check for macroservice
// if unit.MacroServiceID == "" {
// 	s.logger.Info("No macroservice for pipeline unit", zap.String("pipeline_unit_id", unit.ID))

// 	updateCtx, updateCancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer updateCancel()

// 	history.Status = StatusCompleted
// 	history.CompletedAt = time.Now()
// 	history.ExecutionTime = history.CompletedAt.Sub(history.StartedAt)
// 	if err := s.repo.UpdateExecutionHistory(updateCtx, *history); err != nil {
// 		s.logger.Error("Failed to update execution history", zap.String("history_id", history.ID), zap.Error(err))
// 	}
// 	s.repo.UpdatePipelineRunStatus(updateCtx, run.ID, StatusCompleted)
// 	s.broadcastPipelineUpdate(ctx, run.ID, "", StatusCompleted, "Pipeline run completed successfully (no macroservice)")
// 	return nil
// }

// // Fetch macroservice details with timeout context
// macroCtx, macroCancel := context.WithTimeout(ctx, 30*time.Second)
// defer macroCancel()

// macro, err := s.repo.GetServiceByID(macroCtx, unit.MacroServiceID)
// if err != nil {
// 	s.logger.Error("Failed to fetch macroservice", zap.String("macro_service_id", unit.MacroServiceID), zap.Error(err))

// 	updateCtx, updateCancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer updateCancel()

// 	history.Status = StatusRejected
// 	history.CompletedAt = time.Now()
// 	history.ExecutionTime = history
