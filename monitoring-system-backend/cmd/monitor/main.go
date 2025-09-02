package main

import (
	"context"
	"log"

	"github.com/badgerv/monitoring-api/internal/app"
	"github.com/badgerv/monitoring-api/internal/websocket"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction() // or zap.NewDevelopment() for dev mode
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // flushes buffer, if any

	wbHub := websocket.NewHub(logger)

	// run websocket hub in background goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go wbHub.Run(ctx) // Start WebSocket hub in a goroutine

	app := app.BootStrap(wbHub)
	defer app.DB.Close()

	//create router
	if err := app.Router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

