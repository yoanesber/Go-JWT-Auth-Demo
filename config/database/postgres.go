package database

import (
	"fmt"
	"os"
	"sync"

	"gorm.io/driver/postgres"        // Import the PostgreSQL driver for GORM
	"gorm.io/gorm"                   // Import GORM for ORM functionalities
	gormLogger "gorm.io/gorm/logger" // Import GORM logger for logging SQL queries
	"gorm.io/gorm/schema"

	"github.com/yoanesber/go-consumer-api-with-jwt/internal/entity"
	"github.com/yoanesber/go-consumer-api-with-jwt/pkg/logger"
)

var (
	once       sync.Once
	db         *gorm.DB
	DBHost     string
	DBPort     string
	DBUser     string
	DBPass     string
	DBName     string
	DBSchema   string
	DBSSLMode  string
	DBTimeZone string
	DBMigrate  string
	DBSeed     string
	DBSeedFile string
	DBLog      string
)

// LoadPostgresEnv loads environment variables from the .env file
// It sets the database connection parameters such as host, port, user, password, etc.
func LoadPostgresEnv() bool {
	DBHost = os.Getenv("DB_HOST")
	DBPort = os.Getenv("DB_PORT")
	DBUser = os.Getenv("DB_USER")
	DBPass = os.Getenv("DB_PASS")
	DBName = os.Getenv("DB_NAME")
	DBSchema = os.Getenv("DB_SCHEMA")
	DBSSLMode = os.Getenv("DB_SSL_MODE")
	DBTimeZone = os.Getenv("DB_TIMEZONE")
	DBMigrate = os.Getenv("DB_MIGRATE")
	DBSeed = os.Getenv("DB_SEED")
	DBSeedFile = os.Getenv("DB_SEED_FILE")
	DBLog = os.Getenv("DB_LOG")

	if DBHost == "" || DBPort == "" || DBUser == "" || DBPass == "" || DBName == "" || DBSchema == "" {
		logger.Panic("One or more required environment variables are not set", nil)
		return false
	}

	return true
}

// InitPostgres initializes the GORM database connection
func InitPostgres() bool {
	isSuccess := true
	once.Do(func() {
		if !LoadPostgresEnv() {
			isSuccess = false
			return
		}

		// Create the connection string
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s search_path=%s",
			DBHost,
			DBPort,
			DBUser,
			DBPass,
			DBName,
			DBSSLMode,
			DBTimeZone,
			DBSchema,
		)

		// Set the log level based on the environment variable
		var logLevel gormLogger.LogLevel
		if DBLog == "INFO" {
			logLevel = gormLogger.Info
		} else if DBLog == "ERROR" {
			logLevel = gormLogger.Error
		} else if DBLog == "SILENT" {
			logLevel = gormLogger.Silent
		} else {
			logLevel = gormLogger.Warn
		}

		// Open the connection using GORM and PostgreSQL driver
		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   DBSchema + ".",
				SingularTable: false,
			},
			Logger: gormLogger.Default.LogMode(logLevel),
		})
		if err != nil {
			logger.Fatal(fmt.Sprintf("Failed to connect to PostgreSQL: %v", err), nil)
			isSuccess = false
			return
		}

		logger.Info("Connected to PostgreSQL database", nil)

		// Migrate the database schema and all tables
		if DBMigrate == "TRUE" {
			if err = MigratePostgres(); err != nil {
				logger.Fatal(fmt.Sprintf("Failed to migrate PostgreSQL database: %v", err), nil)
				isSuccess = false
				return
			}
		}
	})

	return isSuccess
}

// MigratePostgres migrates the PostgreSQL database schema
// It creates the schema if it does not exist, sets the search path, and migrates the tables.
func MigratePostgres() error {
	// Create the schema in the database
	if DBSchema != "" {
		if err := db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", DBSchema)).Error; err != nil {
			return fmt.Errorf("failed to create schema %s: %v", DBSchema, err)
		}
		logger.Info(fmt.Sprintf("Schema %s created successfully", DBSchema), nil)

		// Set the schema for the database connection
		if err := db.Exec(fmt.Sprintf("SET search_path TO %s", DBSchema)).Error; err != nil {
			return fmt.Errorf("failed to set search path to schema %s: %v", DBSchema, err)
		}
		logger.Info(fmt.Sprintf("Search path set to schema %s", DBSchema), nil)
	} else {
		return fmt.Errorf("DB_SCHEMA environment variable is not set")
	}

	// Perform database migration within a transaction
	err := db.Transaction(func(tx *gorm.DB) error {
		// Check if the transaction is valid
		if tx == nil {
			return fmt.Errorf("transaction is nil")
		}

		// Drop and recreate tables if they exist
		err := tx.Migrator().DropTable(
			&entity.Consumer{},
			&entity.User{},
			&entity.Role{},
			&entity.UserRole{},
			&entity.RefreshToken{})
		if err != nil {
			return fmt.Errorf("failed to drop tables: %v", err)
		}

		// Migrate the database schema
		err = tx.AutoMigrate(
			&entity.Role{},
			&entity.User{},
			&entity.RefreshToken{},
			&entity.Consumer{})
		if err != nil {
			return fmt.Errorf("failed to migrate database: %v", err)
		}

		if DBSeed == "TRUE" {
			// Import initial data from the seed file
			if DBSeedFile == "" {
				return fmt.Errorf("DB_SEED_FILE environment variable is not set")
			}

			// Read the seed file
			seedData, err := os.ReadFile(DBSeedFile)
			if err != nil {
				return fmt.Errorf("failed to read seed file: %v", err)
			}

			// Execute the seed data
			if err := tx.Exec(string(seedData)).Error; err != nil {
				return fmt.Errorf("failed to execute seed data: %v", err)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("database migration failed: %v", err)
	}

	logger.Info("Database migrated successfully", nil)

	return nil
}

// GetPostgres returns the GORM database instance
func GetPostgres() *gorm.DB {
	if db == nil {
		if !InitPostgres() {
			logger.Panic("Failed to initialize PostgreSQL database", nil)
			return nil
		}
	}
	return db
}

// ClosePostgres closes the database connection (optional, for when needed)
func ClosePostgres() {
	sqlDB, err := db.DB()
	if err != nil || sqlDB == nil {
		logger.Error(fmt.Sprintf("Failed to get SQL DB from GORM: %v", err), nil)
		return
	}

	if err := sqlDB.Close(); err != nil {
		logger.Error(fmt.Sprintf("Failed to close database connection: %v", err), nil)
	}

	once = sync.Once{} // Reset the once to allow re-initialization
	db = nil           // Clear the db variable to prevent further use
	logger.Info("Database connection closed successfully", nil)
}
