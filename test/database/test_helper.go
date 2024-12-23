package database

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"NAME/model"
)

var testDB *gorm.DB

// setupTestDB creates a test database and returns a cleanup function
// Required environment variables:
// - TEST_POSTGRES_HOST: PostgreSQL host (default: localhost)
// - TEST_POSTGRES_PORT: PostgreSQL port (default: 5432)
// - TEST_POSTGRES_USER: PostgreSQL user (default: postgres)
// - TEST_POSTGRES_PASSWORD: PostgreSQL password (required)
func setupTestDB(t *testing.T) func() {
	// Test database configuration
	dbName := "name_unit_test"
	dbHost := getEnvOrDefault("TEST_POSTGRES_HOST", "localhost")
	dbPort := getEnvOrDefault("TEST_POSTGRES_PORT", "5432")
	dbUser := getEnvOrDefault("TEST_POSTGRES_USER", "postgres")
	dbPass := os.Getenv("TEST_POSTGRES_PASSWORD")

	if dbPass == "" {
		t.Skip("TEST_POSTGRES_PASSWORD environment variable not set")
	}

	// Create test database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		dbHost, dbPort, dbUser, dbPass)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer sqlDB.Close()

	// Drop test database if it exists
	db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))

	// Create test database
	err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)).Error
	require.NoError(t, err)

	// Connect to test database
	dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)
	testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	// Enable required extensions
	err = testDB.Exec("CREATE EXTENSION IF NOT EXISTS ltree").Error
	require.NoError(t, err)

	// Auto migrate tables
	err = testDB.AutoMigrate(
		&model.Attachment{},
		&model.Comment{},
		&model.Content{},
		&model.User{},
		&model.Tag{},
		&model.Setting{},
		&model.Category{},
	)
	require.NoError(t, err)

	// Return cleanup function
	return func() {
		// Close test database connection
		if testDB != nil {
			sqlDB, err := testDB.DB()
			if err == nil {
				sqlDB.Close()
			}
		}

		// Connect to default database to drop test database
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
			dbHost, dbPort, dbUser, dbPass)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			// Terminate all connections to the test database
			err = db.Exec(fmt.Sprintf(`
				SELECT pg_terminate_backend(pg_stat_activity.pid)
				FROM pg_stat_activity
				WHERE pg_stat_activity.datname = '%s'
				AND pid <> pg_backend_pid()`, dbName)).Error
			if err == nil {
				// Now we can safely drop the database
				db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
			}
			sqlDB, err := db.DB()
			if err == nil {
				sqlDB.Close()
			}
		}
	}
}

// getTestDB returns the test database instance
func getTestDB(t *testing.T) *gorm.DB {
	if testDB == nil {
		t.Fatal("Test database not initialized")
	}
	return testDB
}

// getEnvOrDefault returns the value of the environment variable or the default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
