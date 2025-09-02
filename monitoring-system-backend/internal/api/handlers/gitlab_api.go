package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/badgerv/monitoring-api/internal/auth"
	"github.com/badgerv/monitoring-api/internal/gitlab"
	"github.com/badgerv/monitoring-api/internal/websocket"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler provides HTTP endpoints for the GitLab PipelineService.
type Handler struct {
	service *gitlab.PipelineService
	logger  *zap.Logger
	wbHub   *websocket.Hub
	auth    *auth.Service
}

// NewGitlabHandler creates a new Handler instance.
func NewGitlabHandler(service *gitlab.PipelineService, logger *zap.Logger, wbHub *websocket.Hub, a *auth.Service) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
		wbHub:   wbHub,
		auth:    a,
	}
}

// CreateService handles the creation of a new macro or micro service.
func (h *Handler) CreateService(c *gin.Context) {
	var req struct {
		GitLabRepoID string             `json:"gitlab_repo_id" binding:"required"`
		Name         string             `json:"name" binding:"required"`
		URL          string             `json:"url" binding:"required"`
		Type         gitlab.ServiceType `json:"type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to decode request body", zap.Error(err))
		c.JSON(400, gin.H{"message": "Invalid request body"})
		return
	}

	service, err := h.service.CreateService(c.Request.Context(), req.GitLabRepoID, req.Name, req.URL, req.Type)
	if err != nil {
		h.logger.Error("Failed to create service", zap.Error(err))
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Success",
		"data":    service,
	})
}

// CreatePipelineUnit handles the creation of a new pipeline unit.
func (h *Handler) CreatePipelineUnit(c *gin.Context) {
	var req struct {
		MacroServiceID  string   `json:"macro_service_id"`
		MicroServiceIDs []string `json:"micro_service_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to decode request body", zap.Error(err))
		c.JSON(400, gin.H{"message": "Invalid request body"})
		return
	}

	unit, err := h.service.CreatePipelineUnit(c.Request.Context(), req.MacroServiceID, req.MicroServiceIDs)
	if err != nil {
		h.logger.Error("Failed to create pipeline unit", zap.Error(err))
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Success",
		"data":    unit,
	})
}

// TriggerPipelineUnit handles triggering a pipeline unit execution.
func (h *Handler) TriggerPipelineUnit(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		RequesterID             string   `json:"requester_id" binding:"required"`
		SelectedMicroServiceIDs []string `json:"selected_micro_service_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to decode request body", zap.Error(err))
		c.JSON(400, gin.H{"message": "Invalid request body"})
		return
	}

	run, err := h.service.TriggerPipelineUnit(c.Request.Context(), id, req.RequesterID, req.SelectedMicroServiceIDs)
	if err != nil {
		h.logger.Error("Failed to trigger pipeline unit", zap.String("pipeline_unit_id", id), zap.Error(err))
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Success",
		"data":    run,
	})
}

// ListAllAuthorizationRequests handles GET /api/gitlab/authorization-requests
func (h *Handler) ListAllAuthorizationRequests(c *gin.Context) {
	requests, err := h.service.ListAllAuthorizationRequests(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to list all authorization requests", zap.Error(err))
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"message": "Success",
		"data":    requests,
	})
}

// TriggerPipelineUnit handles triggering a pipeline unit execution.
func (h *Handler) GetPipelineServices(c *gin.Context) {
	id := c.Param("id")
	// var req struct {
	// 	PipelineUnitID string `json:"requester_id" binding:"required"`
	// }
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	h.logger.Error("Failed to decode request body", zap.Error(err))
	// 	c.JSON(400, gin.H{"message": "Invalid request body"})
	// 	return
	// }

	unit, macroservice, microservices, err := h.service.GetPipelineUnitWithServices(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get pipeline unit", zap.String("pipeline_unit_id", id), zap.Error(err))
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	type Data struct {
		Unit          any
		MacroService  any
		MicroServices any
	}

	data := Data{
		Unit:          unit,
		MacroService:  macroservice,
		MicroServices: microservices,
	}

	c.JSON(200, gin.H{
		"message": "Success",
		"data":    data,
	})
}

// ListAllExecutionHistories handles GET /api/gitlab/execution-histories
func (h *Handler) ListAllExecutionHistories(c *gin.Context) {
	// Extract and validate token
	authHeader := c.GetHeader("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" || token == authHeader {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
		return
	}

	// Validate token -> get user
	user, err := h.auth.ValidateToken(c.Request.Context(), token)
	if err != nil {
		h.logger.Error("Failed to validate token", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	if user == nil || user.User == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Fetch execution histories
	histories, err := h.service.ListAllExecutionHistories(c.Request.Context(), user.User.ID.String())
	if err != nil {
		h.logger.Error("Failed to list all execution histories", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch execution histories"})
		return
	}

	// Success
	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    histories,
	})
}


// ApprovePipelineRun handles approving a pipeline run.
func (h *Handler) ApprovePipelineRun(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		ApproverID string `json:"approver_id" binding:"required"`
		Comment    string `json:"comment" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to decode request body", zap.Error(err))
		c.JSON(400, gin.H{"message": "Invalid request body"})
		return
	}

	if err := h.service.ApprovePipelineRun(c.Request.Context(), id, req.ApproverID, req.Comment); err != nil {
		h.logger.Error("Failed to approve pipeline run", zap.String("auth_request_id", id), zap.Error(err))
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Success"})
}

// RejectPipelineRun handles rejecting a pipeline run.
func (h *Handler) RejectPipelineRun(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Comment string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to decode request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	authHeader := c.GetHeader("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
		c.Abort()
		return
	}
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing Authorization header"})
		return
	}

	user, err := h.auth.ValidateToken(c.Request.Context(), token)
	if err != nil {
		h.logger.Error("Failed to validate token", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		return
	}

	// Convert UUID to string properly
	userID := user.User.ID.String()

	fmt.Println(userID)

	if err := h.service.RejectPipelineRun(c.Request.Context(), id, userID, req.Comment); err != nil {
		h.logger.Error("Failed to reject pipeline run", zap.String("pipeline_run_id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

// GetPipelineRunStatus handles retrieving the status of a pipeline run.
func (h *Handler) GetPipelineRunStatus(c *gin.Context) {
	id := c.Param("id")
	status, err := h.service.GetPipelineRunStatus(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get pipeline run status", zap.String("pipeline_run_id", id), zap.Error(err))
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Success",
		"data": gin.H{
			"status": status,
		},
	})
}

// ListExecutionHistory handles retrieving the execution history for a pipeline run.
func (h *Handler) ListExecutionHistory(c *gin.Context) {
	id := c.Param("id")
	history, err := h.service.ListExecutionHistory(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to list execution history", zap.String("pipeline_run_id", id), zap.Error(err))
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Success",
		"data":    history,
	})
}

func (h *Handler) ListPipelineUnits(c *gin.Context) {
	units, macroServices, microServicesList, err := h.service.ListPipelineUnitsWithServices(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to list pipeline units", zap.Error(err))
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	response := make([]gin.H, len(units))
	for i, unit := range units {
		response[i] = gin.H{
			"pipeline_unit":  unit,
			"macro_service":  macroServices[i],
			"micro_services": microServicesList[i],
		}
	}
	c.JSON(200, gin.H{
		"message": "Success",
		"data":    response,
	})
}

func (h *Handler) GetPipelineStatusByRunID(c *gin.Context) {
	id := c.Param("id")

	pipelines, err := h.service.GetPipelineStatusByRunID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to list pipeline", zap.Error(err))
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Success",
		"data":    pipelines,
	})
}

func (h *Handler) GetPipelineUnit(c *gin.Context) {
	id := c.Param("id")
	unit, macroService, microServices, err := h.service.GetPipelineUnitWithServices(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get pipeline unit", zap.String("id", id), zap.Error(err))
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	if unit.ID == "" {
		c.JSON(404, gin.H{"message": "Pipeline unit not found"})
		return
	}
	c.JSON(200, gin.H{
		"message": "Success",
		"data": gin.H{
			"pipeline_unit":  unit,
			"macro_service":  macroService,
			"micro_services": microServices,
		},
	})
}

// ListAllServices handles GET /api/gitlab/services
func (h *Handler) ListAllServices(c *gin.Context) {
	microServices, err := h.service.ListAllServices(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to list all services", zap.Error(err))
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"message": "Success",
		"data": gin.H{
			"services": microServices,
		},
	})
}

// ListAllServices handles GET /api/gitlab/services
func (h *Handler) GetServiceByID(c *gin.Context) {
	microServices, err := h.service.GetServiceByID(c.Request.Context(), "458730af-6988-461d-8052-630e4c2e98ba")
	if err != nil {
		h.logger.Error("Failed to list service", zap.Error(err))
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"message": "Success",
		"data": gin.H{
			"services": microServices,
		},
	})
}

// HandleWebSocket handles WebSocket connections for pipeline run updates.
func (h *Handler) HandleWebSocket(c *gin.Context) {

	id := c.Param("id")

	fmt.Printf("WebSocket route hit - entity_id: %s\n", id) // Add this line
	h.wbHub.HandleWebSocket(c.Writer, c.Request, id)
}
