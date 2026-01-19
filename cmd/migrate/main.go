package main

import (
	"fmt"
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"OpsGo/internal/domain/entity/devops"
	"OpsGo/internal/infrastructure/config"
)

// Blind import to ensure modernc is registered if gorm driver doesn't automatically do it efficiently
// But gorm.io/driver/sqlite usually uses CGO sqlite3 by default, we need to make sure we use the one compatible.
// Actually standard gorm.io/driver/sqlite requires CGO.
// To use modernc with gorm, we might need a specific dialector or tag.
// Ref: https://github.com/go-gorm/sqlite usually uses github.com/mattn/go-sqlite3
// To use modernc, we usually use "gorm.io/driver/sqlite" but with a build tag or specific initializers,
// OR we use a custom driver wrapper.
// However, the user said "modernc.org/sqlite".
// Let's use the standard `gorm.io/driver/sqlite` but we might need to be careful if it tries to link CGO.
// Setting CGO_ENABLED=0 usually forces pure Go, but if `go-sqlite3` is dependency it might fail.
// A common trick is `import _ "modernc.org/sqlite"` and then open with a specific DSN format?
// No, GORM's sqlite driver is hardcoded to use `github.com/mattn/go-sqlite3`.
//
// Wait, if I want to use modernc/sqlite with GORM, I check `gorm.io/driver/sqlite`...
// The official driver imports `github.com/mattn/go-sqlite3`.
// To use modernc, there used to be `github.com/glebarez/sqlite` which is a GORM driver for pure go sqlite (modernc or others).
//
// Let me verify if I added `github.com/glebarez/sqlite`? No I added `modernc.org/sqlite`.
// If I use `gorm.io/driver/sqlite`, it WILL require CGO.
//
// CORRECTION: I should use `github.com/glebarez/sqlite` for modernc support with GORM.
// OR I should use the `modernc.org/sqlite` directly if I wasn't using GORM, but OpsGo uses GORM.
//
// I will check if I can use `github.com/glebarez/sqlite` instead. This is the standard "GORM for modernc/sqlite" driver.
//
// I will change the plan slightly to get `github.com/glebarez/sqlite` which wraps modernc.
//
// I will auto-correct my dependency installation in the next step.
// For now, I'll write the code assuming `gitlab.com/cznic/sqlite` or similar via glebarez.
//
// Reference: https://github.com/glebarez/sqlite
// import "github.com/glebarez/sqlite"
// db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
//
// So the import should be `github.com/glebarez/sqlite` NOT `gorm.io/driver/sqlite`.

func main() {
	// 1. Load Config
	if err := config.LoadConfig("config.yaml"); err != nil {
		// As fallback try local default or env
		log.Println("Could not load config.yaml, assuming environment variables or defaults")
		// We need the MySQL credentials though.
		// If fails, we probably cant migrate.
	}

	// MySQL DSN
	cfg := config.AppConfig.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	log.Printf("Connecting to MySQL source: %s...", cfg.Host)
	mysqlDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	log.Println("Connected to MySQL.")

	// SQLite Destination
	sqliteFile := "opsgo.db"
	log.Printf("Connecting to SQLite destination: %s...", sqliteFile)
	// We need to use the glebarez driver for pure go
	// I will define the opening logic here assuming I fix the import below.
	sqliteDB, err := gorm.Open(sqlite.Open(sqliteFile), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to SQLite: %v", err)
	}
	log.Println("Connected to SQLite.")

	// 2. AutoMigrate Schema
	log.Println("Migrating schema...")
	err = sqliteDB.AutoMigrate(
		&devops.RepoConfig{},
		&devops.PipelineRecord{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate schema: %v", err)
	}

	// 3. Migrate Data
	// RepoConfigs
	var configs []devops.RepoConfig
	if err := mysqlDB.Find(&configs).Error; err != nil {
		log.Fatalf("Failed to list configs: %v", err)
	}
	log.Printf("Found %d RepoConfigs", len(configs))

	if len(configs) > 0 {
		if err := sqliteDB.Create(&configs).Error; err != nil {
			log.Printf("Error creating configs: %v", err)
		} else {
			log.Printf("Migrated %d RepoConfigs", len(configs))
		}
	}

	// PipelineRecords
	// Process in batches to be safe
	var records []devops.PipelineRecord
	batchSize := 100
	var count int64
	mysqlDB.Model(&devops.PipelineRecord{}).Count(&count)
	log.Printf("Found %d PipelineRecords", count)

	var offset int
	for {
		if err := mysqlDB.Offset(offset).Limit(batchSize).Find(&records).Error; err != nil {
			log.Fatalf("Failed to query records: %v", err)
		}
		if len(records) == 0 {
			break
		}

		if err := sqliteDB.Create(&records).Error; err != nil {
			log.Printf("Error creating batch %d: %v", offset, err)
		} else {
			log.Printf("Migrated batch of %d records", len(records))
		}

		offset += len(records)
	}

	log.Println("Migration complete!")
}
