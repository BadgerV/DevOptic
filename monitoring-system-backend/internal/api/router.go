package api

import (
	"github.com/badgerv/monitoring-api/internal/api/handlers"
	"github.com/badgerv/monitoring-api/internal/auth"
	"github.com/badgerv/monitoring-api/internal/gitlab"
	"github.com/badgerv/monitoring-api/internal/rbac"
	"github.com/badgerv/monitoring-api/internal/websocket"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ApiRouter(
	mh *handlers.API,
	ah *handlers.AuthAPI,
	authMiddleware gin.HandlerFunc,
	rbacService *rbac.Service,
	gitlabService *gitlab.PipelineService,
	wbHub *websocket.Hub,
	au *auth.Service,

) *gin.Engine {
	// Initialize Gin router
	r := gin.Default()

	// Add CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://192.168.20.20:3000", "http://192.9.201.92:3000"}, // React frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}))

	// ================== Monitoring Endpoints ==================
	monitor := r.Group("/api/v1/monitor")
	{

		// Admin only (multiple roles allowed)
		monitor.Use(authMiddleware, rbacService.RequireRole("admin", "super admin", "devops", "senior-developer", "developer", "qa-engineer"))
		{
			// Public (no auth needed)
			monitor.GET("/get-overall-stats", mh.GetAggregateStats)
			monitor.GET("/check-scheduler-status", mh.GetSchedulerStatus)
			monitor.GET("/get-endpoint-by-id/:id", mh.GetEndpointDetailByID)
			monitor.GET("/get-endpoint-essentials", mh.GetAllEndpointEssentials)
			monitor.GET("/:id/check", mh.CheckEndpointHandler)
		}

		monitor.Use(authMiddleware, rbacService.RequireRole("admin", "super admin", "devops"))
		{
			monitor.POST("/start-checks", mh.StartEndPointChecks)
			monitor.POST("/stop-checks", mh.StopEndPointChecks)
			monitor.POST("/create-endpoint", mh.CreateEndpoint)
		}

	}

	// ================== Auth Endpoints ==================
	auth := r.Group("/api/v1/auth")
	{
		// Public
		auth.GET("/health", ah.HealthCheck)
		auth.POST("/register", ah.Register)
		auth.POST("/login", ah.Login)

		// Protected
		auth.Use(authMiddleware)
		{
			auth.POST("/set-delivery-email", ah.SetDeliveryEmail)
			auth.GET("/get-delivery-email", ah.GetDeliveryEmail)
			auth.POST("/change-password", ah.ChangePassword)
			auth.POST("/logout", ah.Logout)
			auth.GET("/validate", ah.Validate)
			auth.GET("/user/:id", ah.GetUserByID)
		}
	}

	// ================== RBAC Endpoints ==================
	rbacHandler := handlers.NewRBACHandler(rbacService)
	rbacRoutes := r.Group("/api/v1/rbac")
	{
		// Health check
		rbacRoutes.GET("/health", rbacHandler.HealthCheck)

		rbacRoutes.GET("/:user_id/is-super-admin", rbacHandler.CheckIfSuperAdmin)

		// Role management (protected, super admin only)
		rbacRoutes.Use(authMiddleware, rbacService.RequireRole("super admin"))
		{
			rbacRoutes.POST("/roles", rbacHandler.CreateRole)
			rbacRoutes.GET("/roles", rbacHandler.GetAllRoles)

			// Permissions
			rbacRoutes.POST("/permissions", rbacHandler.CreatePermission)

			// Role ↔ User
			rbacRoutes.POST("/assign-role", rbacHandler.AssignRoleToUser)
			rbacRoutes.POST("/remove-role", rbacHandler.RemoveRoleFromUser)

			// Role ↔ Permission
			rbacRoutes.POST("/assign-permission", rbacHandler.AssignPermissionToRole)

			// User Permission Queries
			rbacRoutes.GET("/user/:id/permissions", rbacHandler.GetUserPermissions)
			rbacRoutes.POST("/check-permission", rbacHandler.CheckPermission)
			rbacRoutes.GET("/check-permission/:user_id", rbacHandler.CheckUserPermission)
			rbacRoutes.GET("/get-all-users-usernames-id", rbacHandler.GetAllUsersUsernameAndID)
		}
	}

	logger, err := zap.NewProduction() // or zap.NewDevelopment() for dev mode
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // flushes buffer, if any

	gitlabHandler := handlers.NewGitlabHandler(gitlabService, logger, wbHub, au)

	gitlabRoutes := r.Group("/api/v1/gitlab", authMiddleware, rbacService.RequireRole("super admin"))
	{
		gitlabRoutes.POST("/services", gitlabHandler.CreateService)
		gitlabRoutes.POST("/pipeline-units", gitlabHandler.CreatePipelineUnit)
		gitlabRoutes.POST("/authorization-requests/:id/approve", gitlabHandler.ApprovePipelineRun)
		gitlabRoutes.POST("/authorization-requests/:id/reject", gitlabHandler.RejectPipelineRun)
	}

	// separate group for senior-developer role
	// seniorDevRoutes := r.Group("/api/v1/gitlab")
	seniorDevRoutes := r.Group("/api/v1/gitlab", authMiddleware, rbacService.RequireRole("senior-developer", "super admin", "developer", "qa-engineer"))
	{
		seniorDevRoutes.GET("/authorization-requests", gitlabHandler.ListAllAuthorizationRequests)
		seniorDevRoutes.GET("/services", gitlabHandler.ListAllServices)
		seniorDevRoutes.POST("/pipeline-units/:id/trigger", gitlabHandler.TriggerPipelineUnit)
		seniorDevRoutes.GET("/pipeline-runs/:id/status", gitlabHandler.GetPipelineRunStatus)
		seniorDevRoutes.GET("/pipeline-runs/:id/history", gitlabHandler.ListExecutionHistory)
		seniorDevRoutes.GET("/pipeline-runs/history", gitlabHandler.ListAllExecutionHistories)
		seniorDevRoutes.GET("/pipeline-unit/get/:id", gitlabHandler.GetPipelineServices)
		seniorDevRoutes.GET("/pipeline-units", gitlabHandler.ListPipelineUnits)
		seniorDevRoutes.GET("/pipeline-units/:id", gitlabHandler.GetPipelineUnit)
		seniorDevRoutes.GET("/services-one", gitlabHandler.GetServiceByID)
		seniorDevRoutes.GET("/pipeline-status/:id", gitlabHandler.GetPipelineStatusByRunID)
	}

	// separate group for senior-developer role
	// webSockerRoutes := r.Group("/api/v1/gitlab")
	webSockerRoutes := r.Group("/api/v1/gitlab", rbacService.RequireRoleForWebsocket("senior-developer", "super admin"))
	{
		webSockerRoutes.GET("/ws/pipeline-runs/:id", gitlabHandler.HandleWebSocket)
	}

	return r
}
