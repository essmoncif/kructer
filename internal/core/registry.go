package core

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"kructer.com/config"
)

type ModelRegistry struct {
	*gorm.DB
	models map[string]reflect.Value
	isOpen bool
}

func NewRegistry() *ModelRegistry {
	return &ModelRegistry{
		models: make(map[string]reflect.Value),
		isOpen: false,
	}
}

func (m *ModelRegistry) IsOpen() bool {
	return m.isOpen
}

func (m *ModelRegistry) OpenWithConfig(config *config.Configuration) error {
	db, err := gorm.Open(config.Dialect, config.ConnectionString)
	if err != nil {
		return err
	}

	db.DB().SetConnMaxLifetime(time.Minute * 5)
	db.DB().SetMaxIdleConns(0)
	db.DB().SetMaxOpenConns(20)

	m.DB = db
	m.isOpen = true

	// Initialize changelog table
	m.InitChangelogTable()

	// Apply migrations after opening the DB
	m.ApplyMigrations()

	return nil
}

// InitChangelogTable initializes the changelog table if it doesn't exist
func (m *ModelRegistry) InitChangelogTable() {
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS changelog (
        id SERIAL PRIMARY KEY,
        script_id VARCHAR(255) UNIQUE NOT NULL,
        developer VARCHAR(255) NOT NULL,
        checksum VARCHAR(64) NOT NULL,
        applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

	if err := m.Exec(createTableQuery).Error; err != nil {
		log.Fatal("Failed to create changelog table:", err)
	}
}

// ApplyMigrations applies SQL migrations and updates the changelog
func (m *ModelRegistry) ApplyMigrations() {
	migrationDir := "./migrations"

	files, err := ioutil.ReadDir(migrationDir)
	if err != nil {
		log.Fatal("Failed to read migrations directory:", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			m.applyMigration(filepath.Join(migrationDir, file.Name()), file.Name())
		}
	}
}

type ChangelogEntry struct {
	ScriptID  string    `gorm:"column:script_id"`
	Developer string    `gorm:"column:developer"`
	Checksum  string    `gorm:"column:checksum"`
	AppliedAt time.Time `gorm:"column:applied_at;default:current_timestamp"`
}

func (m *ModelRegistry) applyMigration(filePath string, fileName string) {
	// Open and read SQL file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Failed to read migration file:", err)
	}
	defer file.Close()

	// Compute checksum
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.Fatal("Failed to compute checksum:", err)
	}
	checksum := fmt.Sprintf("%x", hash.Sum(nil))

	// Reset file pointer to beginning
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		log.Fatal("Failed to reset file pointer:", err)
	}

	scanner := bufio.NewScanner(file)
	var sql string
	var developer, scriptID string

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "--") {
			parts := strings.SplitN(strings.TrimPrefix(line, "--"), ":", 2)
			if len(parts) == 2 {
				developer = strings.TrimSpace(parts[0])
				scriptID = strings.TrimSpace(parts[1])
				log.Printf("Parsed developer: %s, scriptID: %s\n", developer, scriptID) // Logging
			}
		} else {
			sql += line + "\n"
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Failed to scan migration file:", err)
	}

	if scriptID == "" || developer == "" {
		log.Fatalf("No script ID or developer found in migration file: %s (scriptID: %s, developer: %s)", fileName, scriptID, developer)
	}

	// Check if migration has already been applied
	var count int64
	m.Table("changelog").Where("script_id = ?", scriptID).Count(&count)
	if count > 0 {
		log.Println("Migration already applied:", scriptID)
		return
	}

	// Execute SQL
	if err := m.Exec(sql).Error; err != nil {
		log.Fatal("Failed to execute migration:", err)
	}

	// Create a ChangelogEntry struct
	changelogEntry := ChangelogEntry{
		ScriptID:  scriptID,
		Developer: developer,
		Checksum:  checksum,
	}

	// Debugging: print the struct being passed to Create
	log.Printf("Changelog entry: %+v\n", changelogEntry)

	// Update changelog
	if err := m.Table("changelog").Create(&changelogEntry).Error; err != nil {
		log.Fatalf("Failed to update changelog (scriptID: %s, developer: %s): %v", scriptID, developer, err)
	}
	log.Println("Migration applied:", scriptID)
}
