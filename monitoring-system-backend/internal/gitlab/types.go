package gitlab

import (
	"time"

	"github.com/google/uuid"
)

// PipelineStatus represents the status of a pipeline unit or run.
type PipelineStatus string

const (
	StatusPending   PipelineStatus = "pending"
	StatusAccepted  PipelineStatus = "accepted"
	StatusRejected  PipelineStatus = "rejected"
	StatusRunning   PipelineStatus = "running"
	StatusCompleted PipelineStatus = "completed"
)

// ServiceType distinguishes between macro and micro services.
type ServiceType string

const (
	MacroService ServiceType = "macro"
	MicroService ServiceType = "micro"
)

// Service represents a GitLab repository registered as a macro or micro service.
type Service struct {
	ID           string      `json:"id"`
	GitLabRepoID string      `json:"gitlab_repo_id"`
	Name         string      `json:"name"`
	URL          string      `json:"url"`
	Type         ServiceType `json:"type"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

// PipelineUnit represents a pipeline definition with one macro service and multiple micro service dependencies.
type PipelineUnit struct {
	ID              string    `json:"id"`
	MacroServiceID  string    `json:"macro_service_id"`
	MicroServiceIDs []string  `json:"micro_service_ids"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// PipelineRun represents a single execution attempt of a pipeline unit.
type PipelineRun struct {
	ID                      string         `json:"id"`
	PipelineUnitID          string         `json:"pipeline_unit_id"`
	Status                  PipelineStatus `json:"status"`
	CreatedAt               time.Time      `json:"created_at"`
	UpdatedAt               time.Time      `json:"updated_at"`
	GitLabPipelineID        int            `json:"gitlab_pipeline_id"`
	ApproverID              uuid.UUID      `json:"approver_id"`
	ExecutionTime           time.Duration  `json:"execution_time"`
	SelectedMicroServiceIDs []string       `json:"selected_micro_service_ids"` // New field for selected microservices
}

// AuthorizationRequest represents a request for pipeline run approval.
type AuthorizationRequest struct {
	ID                string         `json:"id"`
	PipelineRunID     string         `json:"pipeline_run_id"`
	RequesterID       string         `json:"requester_id"`
	RequesterName     string         `json:"requester_name"`
	ApproverID        *uuid.UUID         `json:"approver_id"`
	ApproverName      string         `json:"approver_name"`
	Status            PipelineStatus `json:"status"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	Comment           string         `json:"comment"`
	MacroServiceName  string         `json:"macro_service_name"`
	MicroServiceNames []string       `json:"micro_service_names"`
}

// ExecutionHistory captures the execution details of a pipeline run.
type ExecutionHistory struct {
	ID                string         `json:"id"`
	PipelineRunID     string         `json:"pipeline_run_id"`
	RequesterID       string         `json:"requester_id"`
	RequesterName     string         `json:"requester_name"`
	ApproverID        uuid.UUID         `json:"approver_id"`
	ApproverName      string         `json:"approver_name"`
	Status            PipelineStatus `json:"status"`
	StartedAt         time.Time      `json:"started_at"`
	CompletedAt       time.Time      `json:"completed_at"`
	ErrorMessage      string         `json:"error_message"`
	MacroServiceName  string         `json:"macro_service_name"`
	MicroServiceNames []string       `json:"micro_service_names"`
	ExecutionTime     time.Duration  `json:"execution_time"`
	PipelineUnitID    uuid.UUID      `json:"pipeline_unit_id"`
}

// WebSocketMessage defines the structure for real-time pipeline updates.
type WebSocketMessage struct {
	Type          string         `json:"type"` // e.g., "status_update", "approval_update"
	PipelineRunID string         `json:"pipeline_run_id"`
	Status        PipelineStatus `json:"status"`
	Timestamp     time.Time      `json:"timestamp"`
	Message       string         `json:"message"` // Additional details
}

// PipelineStatusResponse represents the structured response for pipeline statuses
type PipelineStatusResponse struct {
	Running   []PipelineRunStatus `json:"running"`
	Pending   []PipelineRunStatus `json:"pending"`
	Completed []PipelineRunStatus `json:"completed"`
	Failed    []PipelineRunStatus `json:"failed"`
	Total     int                 `json:"total"`
}

// PipelineRunStatus represents essential pipeline run information for status display
type PipelineRunStatus struct {
	ID                string    `json:"id"`
	PipelineUnitID    string    `json:"pipeline_unit_id"`
	Status            string    `json:"status"`
	MacroServiceName  string    `json:"macro_service_name,omitempty"`
	MicroServiceNames []string  `json:"micro_service_names"`
	RequesterName     string    `json:"requester_name,omitempty"`
	ApproverName      string    `json:"approver_name,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	GitLabPipelineID  int       `json:"gitlab_pipeline_id,omitempty"`
}