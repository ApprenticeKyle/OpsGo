package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"OpsGo/internal/domain/entity/devops"
)

func main() {
	// 1. Ensure data directory exists
	dbFile := "data/opsgo.db"
	if err := os.MkdirAll(filepath.Dir(dbFile), 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// 2. Connect to SQLite
	log.Printf("Connecting to SQLite database: %s...", dbFile)
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to SQLite: %v", err)
	}
	log.Println("Connected to SQLite.")

	// 3. AutoMigrate Schema
	log.Println("Migrating schema...")
	err = db.AutoMigrate(
		&devops.RepoConfig{},
		&devops.PipelineRecord{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate schema: %v", err)
	}

	log.Println("Migration complete! Tables 'devops_repo_configs' and 'devops_pipeline_records' are ready.")
}
