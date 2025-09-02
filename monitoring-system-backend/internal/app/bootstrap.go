package app

import (
	"log"

	"github.com/badgerv/monitoring-api/internal/api"
	"github.com/badgerv/monitoring-api/internal/api/handlers"
	"github.com/badgerv/monitoring-api/internal/auth"
	"github.com/badgerv/monitoring-api/internal/config"
	"github.com/badgerv/monitoring-api/internal/emailservice"
	"github.com/badgerv/monitoring-api/internal/gitlab"
	"github.com/badgerv/monitoring-api/internal/monitor"
	"github.com/badgerv/monitoring-api/internal/rbac"
	"github.com/badgerv/monitoring-api/internal/websocket"

	// "github.com/badgerv/monitoring-api/internal/rbac"
	"os"

	"github.com/badgerv/monitoring-api/internal/storage"
	"go.uber.org/zap"

	gitlabapi "gitlab.com/gitlab-org/api/client-go"
)

func BootStrap(wbHub *websocket.Hub) *Application {
	// Load Environment Variables
	config.InitEnvVariables()

	//load gitlab env vairbale
	gitlab_token := os.Getenv("GITLAB_TOKEN")

	// Initialize DB
	db := storage.NewDB()

	// --- Monitor setup ---
	monitorRepo := monitor.NewPostgresRepository(db)
	monitorService := monitor.NewService(db, monitorRepo)
	monitorApiHandler := handlers.NewMonitorHandle(monitorService)

	//Rbac setup
	rbacRepo := rbac.NewPostgresRepository(db)
	rbacService := rbac.NewService(rbacRepo)

	// --- Auth setup ---
	userRepo := auth.NewPostgresUserRepository(db)
	sessionRepo := auth.NewPostgresSessionRepository(db)

	authService := auth.NewService(userRepo, sessionRepo)
	authApiHandler := handlers.NewAuthHandler(authService)

	// --Gitlab Service --

	logger, err := zap.NewProduction() // or zap.NewDevelopment() for dev mode
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // flushes buffer, if any

	git, err := gitlabapi.NewClient(gitlab_token, gitlabapi.WithBaseURL("https://gitlab.remita.net/api/v4"))
	if err != nil {
		log.Fatalf("Failed to create GitLab client: %v", err)
	}

	gitlabRepo, _ := gitlab.NewPostgresRepository(db, logger, rbacService)
	gitlabService := gitlab.NewPipelineService(gitlabRepo, git, emailservice.NewEmailService(), logger, wbHub, userRepo)

	// --- Router ---
	router := api.ApiRouter(monitorApiHandler, authApiHandler, authService.AuthMiddleware(), rbacService, gitlabService, wbHub, authService)

	return &Application{DB: db, Router: router}
}
