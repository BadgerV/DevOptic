package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/badgerv/monitoring-api/internal/monitor"
)

type API struct {
	Monitor   *monitor.Service
	scheduler *monitor.Scheduler
	mu        sync.Mutex // Protect scheduler access
}

type Endpoint struct {
	ID                  int      `json:"id,omitempty"`
	ServiceName         string   `json:"service_name"`
	URL                 string   `json:"url"`
	ServerName          string   `json:"server_name"`
	APIMethod           string   `json:"api_method"`
	ExpectedCode        int      `json:"expected_status_code"`
	GitlabURL           *string  `json:"gitlab_url,omitempty"`
	DockerContainerName *string  `json:"docker_container_name,omitempty"`
	KubernetesPodName   *string  `json:"kubernetes_pod_name,omitempty"`
	Tags                []string `json:"tags,omitempty"`
	Description         *string  `json:"description,omitempty"`
	LastChangedBy       *string  `json:"last_changed_by,omitempty"`
}

func NewMonitorHandle(m *monitor.Service) *API {
	return &API{
		Monitor: m,
	}
}

func (a *API) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (a *API) StartEndPointChecks(c *gin.Context) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Check if scheduler is already running
	if a.scheduler != nil {
		log.Println("Endpoints check already running")
		c.JSON(http.StatusConflict, gin.H{"message": "Endpoints check already running"})
		return
	}

	ctx := context.Background()

	scheduler, err := a.Monitor.StartScheduler(ctx)
	if err != nil {
		log.Println("Endpoints check failed to start:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Endpoints check failed"})
		return
	}

	a.scheduler = scheduler
	log.Println("Endpoints check started successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Endpoints check started successfully"})
}

func (a *API) StopEndPointChecks(c *gin.Context) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.scheduler == nil {
		log.Println("No endpoints check running")
		c.JSON(http.StatusConflict, gin.H{"message": "No endpoints check running"})
		return
	}

	a.scheduler.Stop()
	a.scheduler = nil

	log.Println("Endpoints check stopped successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Endpoints stopped successfully"})
}

// GetSchedulerStatus returns whether the scheduler is running
func (a *API) GetSchedulerStatus(c *gin.Context) {
	a.mu.Lock()
	defer a.mu.Unlock()

	isRunning := a.scheduler != nil
	c.JSON(http.StatusOK, gin.H{
		"scheduler_running": isRunning,
	})
}

func (a *API) GetEndpointDetailByID(c *gin.Context) {
	idParam := c.Param("id") // from URL, e.g. /endpoints/:id
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}

	endpointDetail, err := a.Monitor.GetEndpointByID(c.Request.Context(), id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    endpointDetail,
	})
}

func (a *API) GetAllEndpointEssentials(c *gin.Context) {
	endpointChecks, err := a.Monitor.GetAllEndpointEssentials(c.Request.Context())

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    endpointChecks,
	})
}

func (a *API) GetAggregateStats(c *gin.Context) {
	endpointChecks, err := a.Monitor.GetAggregateStats(c.Request.Context())

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    endpointChecks,
	})
}

func (a *API) CreateEndpoint(c *gin.Context) {
	var ep Endpoint
	if err := c.ShouldBindJSON(&ep); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}
	monitorEp := &monitor.Endpoint{
		ID:                  ep.ID,
		ServiceName:         ep.ServiceName,
		ServerName:          ep.ServerName,
		URL:                 ep.URL,
		APIMethod:           ep.APIMethod,
		ExpectedCode:        ep.ExpectedCode,
		GitlabURL:           ep.GitlabURL,
		DockerContainerName: ep.DockerContainerName,
		KubernetesPodName:   ep.KubernetesPodName,
		Tags:                ep.Tags,
		Description:         ep.Description,
		LastChangedBy:       ep.LastChangedBy,
	}

	createdEp, err := a.Monitor.CreateEndpoint(c.Request.Context(), monitorEp)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    createdEp,
	})
}

// CheckEndpointHandler handles requests to check an endpoint's status
func (a *API) CheckEndpointHandler(c *gin.Context) {
    // Parse endpoint ID from URL
    idStr := c.Param("id")
    endpointID, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Invalid endpoint ID",
        })
        return
    }

    // Get endpoint details
    endpoint, err := a.Monitor.GetEndpointByID(context.Background(), endpointID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "message": "Endpoint not found",
            "error":   err.Error(),
        })
        return
    }

    // Call CheckEndpointStatus and get latency
    latency, err := a.Monitor.CheckEndpointStatus(
        context.Background(),
        endpoint.URL,
        endpoint.APIMethod,
        endpoint.ExpectedCode,
    )

    // Build response
    status := "unreachable"
    if err == nil {
        status = "up"
    }

    c.JSON(http.StatusOK, gin.H{
        "endpoint_id": endpointID,
        "service_name": endpoint.ServiceName,
        "server_name": endpoint.ServerName,
        "status": status,
        "latency_ms": latency.Milliseconds(), // return latency in milliseconds
        "error": func() string {
            if err != nil {
                return err.Error()
            }
            return ""
        }(),
    })
}
