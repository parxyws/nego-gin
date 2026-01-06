package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/parxyws/nego-gin/config"
	"github.com/parxyws/nego-gin/pkg/database/psql"
	"github.com/parxyws/nego-gin/pkg/database/rd"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// SetupTestDB creates a database connection for testing using the application config
func SetupTestDB(t *testing.T) *gorm.DB {
	cfg, err := InitEnv()
	if err != nil {
		t.Fatalf("failed to initialize environment: %v", err)
	}

	db, err := psql.NewDB(cfg)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	return db
}

func SetupTestRedis(t *testing.T) *redis.Client {
	cfg, err := InitEnv()
	if err != nil {
		t.Fatalf("failed to initialize environment: %v", err)
	}

	db := rds.NewRedis(cfg)

	return db
}

// CleanupTestRedis flushes the current Redis database
func CleanupTestRedis(rdb *redis.Client) error {
	if rdb == nil {
		return fmt.Errorf("redis connection is nil")
	}
	return rdb.FlushDB(context.Background()).Err()
}

// CleanupTestDB cleans up all test data from the database
// This function truncates all tables to ensure a clean state between tests
func CleanupTestDB(db *gorm.DB, tables ...string) error {
	// Guard against nil db
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	ctx := context.Background()

	// Disable foreign key checks temporarily
	if err := db.WithContext(ctx).Exec("SET session_replication_role = 'replica'").Error; err != nil {
		return fmt.Errorf("failed to disable foreign key checks: %w", err)
	}

	// Truncate specified tables
	for _, table := range tables {
		if err := db.WithContext(ctx).Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error; err != nil {
			return fmt.Errorf("failed to truncate table %s: %w", table, err)
		}
	}

	// Re-enable foreign key checks
	if err := db.WithContext(ctx).Exec("SET session_replication_role = 'origin'").Error; err != nil {
		return fmt.Errorf("failed to re-enable foreign key checks: %w", err)
	}

	return nil
}

// CleanupAllTables removes all data from common test tables
func CleanupAllTables(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	tables := []string{
		"users",
		// Add other table names as needed
	}

	return CleanupTestDB(db, tables...)
}

// InitTestDB initializes a test database connection with the given config
func InitTestDB(cfg *config.Config) (*gorm.DB, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}
	return psql.NewDB(cfg)
}
