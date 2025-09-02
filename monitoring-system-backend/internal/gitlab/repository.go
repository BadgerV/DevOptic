package gitlab

import (
	"context"
)

// Repository defines the interface for interacting with GitLab pipeline data.
type Repository interface {
	// Service management
	CreateService(ctx context.Context, service Service) (Service, error)
	GetServiceByID(ctx context.Context, id string) (Service, error)
	GetServiceByGitLabRepoID(ctx context.Context, gitlabRepoID string) (Service, error)
	ListServicesByType(ctx context.Context, serviceType ServiceType) ([]Service, error)

	// PipelineUnit management
	CreatePipelineUnit(ctx context.Context, unit PipelineUnit) (PipelineUnit, error)
	GetPipelineUnit(ctx context.Context, id string) (PipelineUnit, error)
	AddMicroServiceDependency(ctx context.Context, pipelineUnitID, microServiceID string, orderIndex int) error
	GetMicroServiceDependencies(ctx context.Context, pipelineUnitID string) ([]string, error)

	// PipelineRun management
	CreatePipelineRun(ctx context.Context, run PipelineRun) (PipelineRun, error)
	GetPipelineRun(ctx context.Context, id string) (PipelineRun, error)
	UpdatePipelineRunStatus(ctx context.Context, id string, status PipelineStatus) error
	UpdatePipelineRun(ctx context.Context, run PipelineRun) error

	// AuthorizationRequest management
	CreateAuthorizationRequest(ctx context.Context, request AuthorizationRequest) (AuthorizationRequest, error)
	GetAuthorizationRequest(ctx context.Context, id string) (AuthorizationRequest, error)
	UpdateAuthorizationRequest(ctx context.Context, id string, status PipelineStatus, comment string, approverID ...string) error
	ListAuthorizationRequestsByPipelineRun(ctx context.Context, pipelineRunID string) ([]AuthorizationRequest, error)

	// ExecutionHistory management
	CreateExecutionHistory(ctx context.Context, history ExecutionHistory) (ExecutionHistory, error)
	GetExecutionHistory(ctx context.Context, id string) (ExecutionHistory, error)
	ListExecutionHistoryByPipelineRun(ctx context.Context, pipelineRunID string) ([]ExecutionHistory, error)
	UpdateExecutionHistoryError(ctx context.Context, id, errorMessage string) error
	UpdateExecutionHistory(ctx context.Context, history ExecutionHistory) error

	// ListAllAuthorizationRequests retrieves all authorization requests.
	ListAllAuthorizationRequests(ctx context.Context) ([]AuthorizationRequest, error)

	ListAllExecutionHistories(ctx context.Context, id string) ([]ExecutionHistory, error)
	ListPipelineUnits(ctx context.Context) ([]PipelineUnit, error)
	GetPipelineUnitWithServices(ctx context.Context, unitID string) (PipelineUnit, Service, []Service, error)

	ListPipelineRunsWithServices(ctx context.Context, runID string) ([]PipelineRunStatus, error)

	GetExecutionHistoryByID(ctx context.Context, userIDStr, historyID string) (*ExecutionHistory, error)
	GetAuthorizationRequestByID(ctx context.Context, id string) (*AuthorizationRequest, error)
}
