package main

import (
	"example.com/database"
	"example.com/log"
	"example.com/server"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const shutDownTimeOut = 10 * time.Second

func main() {
	log.Init()
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// create server instance
	srv := server.SetupRoutes()
	if err := database.ConnectAndMigrate(
		os.Getenv("host"),
		os.Getenv("port"),
		os.Getenv("databaseName"),
		os.Getenv("user"),
		os.Getenv("password"),
		database.SSLModeDisable); err != nil {
		logrus.Fatalf("Failed to initialize and migrate database with error: %+v", err)
	}
	log.Info("migration successful!!")

	go func() {
		if err := srv.Run(":8084"); err != nil && err != http.ErrServerClosed {
			logrus.Panicf("Failed to run server with error: %+v", err)
		}
	}()
	log.Info("Server started at :8084")

	<-done

	log.Info("shutting down server")
	if err := database.ShutdownDatabase(); err != nil {
		logrus.WithError(err).Error("failed to close database connection")
	}

	if err := srv.Shutdown(shutDownTimeOut); err != nil {
		logrus.WithError(err).Panic("failed to gracefully shutdown server")
	}
}
