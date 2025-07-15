package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/yoanesber/go-consumer-api-with-jwt/config/database"
	"github.com/yoanesber/go-consumer-api-with-jwt/pkg/diagnostics"
	"github.com/yoanesber/go-consumer-api-with-jwt/pkg/logger"
	validation "github.com/yoanesber/go-consumer-api-with-jwt/pkg/util/validation-util"
	"github.com/yoanesber/go-consumer-api-with-jwt/routes"
)

var (
	validatorInitialized bool
	dbInitialized        bool
)

func init() {
	logger.Init()
}

func main() {
	// Create base context with cancel for graceful shutdown
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Get environment variables
	env := os.Getenv("ENV")
	port := os.Getenv("PORT")
	isSSL := os.Getenv("IS_SSL")
	apiVersion := os.Getenv("API_VERSION")
	sslKeys := os.Getenv("SSL_KEYS")
	sslCert := os.Getenv("SSL_CERT")

	if env == "" || port == "" || isSSL == "" || apiVersion == "" || sslKeys == "" || sslCert == "" {
		logger.Panic("One or more required environment variables are not set", log.Fields{
			"ENV":         env,
			"PORT":        port,
			"IS_SSL":      isSSL,
			"API_VERSION": apiVersion,
			"SSL_KEYS":    sslKeys,
			"SSL_CERT":    sslCert,
		})
		return
	}

	// Set Gin mode
	gin.SetMode(gin.DebugMode)
	if env == "PRODUCTION" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Setup router
	r := routes.SetupRouter()
	r.SetTrustedProxies(nil) // Set trusted proxies to nil to avoid issues with forwarded headers

	// Log memory stats before initialization
	diagnostics.LogMemoryStats("Before initialization")

	// Init all dependencies
	initializeDependencies()

	// Log memory stats after initialization
	diagnostics.LogMemoryStats("After initialization")

	// Graceful shutdown
	gracefulShutdown(cancel)

	// Start the server
	var err error
	if isSSL == "TRUE" {
		//Generated using sh generate-certificate.sh
		err = r.RunTLS(":"+port, sslCert, sslKeys)

	} else {
		err = r.Run(":" + port)
	}

	if err != nil {
		logger.Error(fmt.Sprintf("Failed to start server with SSL: %v", err), log.Fields{
			"environment": env,
			"port":        port,
			"is_ssl":      isSSL,
			"api_version": apiVersion,
			"ssl_cert":    sslCert,
			"ssl_keys":    sslKeys,
		})
		return
	}
}

func initializeDependencies() {
	if !validatorInitialized {
		if !validation.Init() {
			logger.Fatal("Failed to initialize validator", nil)
		} else {
			validatorInitialized = true
		}
	}

	if !dbInitialized {
		if !database.InitPostgres() {
			logger.Fatal("Failed to initialize Postgres database", nil)
		} else {
			dbInitialized = true
		}
	}
}

func gracefulShutdown(cancel context.CancelFunc) {
	// Handle graceful shutdown signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-quit
		logger.Info(fmt.Sprintf("Received signal: %s. Initiating graceful shutdown...", sig), nil)

		// Cancel context
		cancel()

		if dbInitialized {
			logger.Info("Closing Postgres connection...", nil)
			database.ClosePostgres()
		}
		if validatorInitialized {
			logger.Info("Clearing validator instance...", nil)
			validation.ClearValidator()
		}

		diagnostics.LogMemoryStats("After shutdown cleanup")

		logger.Info("Shutdown complete. Bye 👋", nil)
		logger.Exit()
		os.Exit(0)
	}()
}
