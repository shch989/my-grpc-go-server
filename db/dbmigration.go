package dbmigration

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// 마이그레이션 경로 가져오기
func getMigrationPath() string {
	// 현재 실행 중인 디렉터리 확인
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("❌ Failed to get current directory: %v", err)
	}

	// 절대 경로로 변환 (Windows 호환)
	migrationPath := filepath.Join(dir, "..", "db", "migrations")

	// ✅ 경로를 Unix 스타일로 변환
	migrationPath = filepath.ToSlash(migrationPath)

	// 경로 확인 로그
	log.Println("✅ Using migration path:", migrationPath)

	// 파일이 실제로 존재하는지 확인
	if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
		log.Fatalf("❌ Migration path does not exist: %s", migrationPath)
	}

	return "file://" + migrationPath
}

func Migrate(conn *sql.DB) {
	log.Println("Database migration start")

	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to create database driver: %v", err)
	}

	// ✅ 절대 경로 사용
	migrationPath := getMigrationPath()

	m, err := migrate.NewWithDatabaseInstance(migrationPath, "postgres", driver)
	if err != nil {
		log.Fatalf("❌ Database migration failed: %v", err)
	}

	// 마이그레이션 실행
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		log.Println("⚠️ Database migration (down) failed:", err)
	}

	if err := m.Up(); err != nil {
		log.Println("⚠️ Database migration (up) failed:", err)
	}

	log.Println("✅ Database migration completed successfully")
}
