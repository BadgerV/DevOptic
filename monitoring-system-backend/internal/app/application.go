package app

import (
	"github.com/badgerv/monitoring-api/internal/storage"
	"github.com/gin-gonic/gin"
)

type Application struct {
	DB     *storage.DB
	Router *gin.Engine
}