// package database

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"

// 	beego "github.com/beego/beego/v2/server/web"
// 	"github.com/golang-migrate/migrate/v4"
// 	"github.com/golang-migrate/migrate/v4/database/postgres"
// 	_ "github.com/golang-migrate/migrate/v4/source/file"
// 	_ "github.com/lib/pq"
// )

// var DB *sql.DB

// func InitDB() {
// 	dbtype, _ := beego.AppConfig.String("dbtype")
// 	dbhost, _ := beego.AppConfig.String("dbhost")
// 	dbport, _ := beego.AppConfig.String("dbport")
// 	dbuser, _ := beego.AppConfig.String("dbuser")
// 	dbpass, _ := beego.AppConfig.String("dbpass")
// 	dbname, _ := beego.AppConfig.String("dbname")
// 	dbsslmode, _ := beego.AppConfig.String("dbsslmode")

// 	connStr := fmt.Sprintf(
// 		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
// 		dbhost, dbport, dbuser, dbpass, dbname, dbsslmode,
// 	)

// 	var err error
// 	DB, err = sql.Open(dbtype, connStr)
// 	if err != nil {
// 		log.Fatal("Failed to connect to database:", err)
// 	}

// 	// Test the connection
// 	err = DB.Ping()
// 	if err != nil {
// 		log.Fatal("Failed to ping database:", err)
// 	}

// 	log.Println("Database connected successfully!")

// 	// Run migrations
// 	RunMigrations()
// }

// // RunMigrations applies all pending database migrations
// func RunMigrations() {
// 	driver, err := postgres.WithInstance(DB, &postgres.Config{})
// 	if err != nil {
// 		log.Fatal("Failed to create migration driver:", err)
// 	}

// 	m, err := migrate.NewWithDatabaseInstance(
// 		"file://migrations",
// 		"postgres", driver)
// 	if err != nil {
// 		log.Fatal("Failed to initialize migrations:", err)
// 	}

// 	// Apply all pending migrations
// 	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
// 		log.Fatal("Failed to run migrations:", err)
// 	}

// 	version, dirty, err := m.Version()
// 	if err != nil && err != migrate.ErrNilVersion {
// 		log.Fatal("Failed to get migration version:", err)
// 	}

// 	if err == migrate.ErrNilVersion {
// 		log.Println("No migrations applied yet")
// 	} else {
// 		log.Printf("Migrations applied successfully! Current version: %d (dirty: %v)", version, dirty)
// 	}
// }

// func CloseDB() {
// 	if DB != nil {
// 		DB.Close()
// 	}
// }

package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var connStr string
	
	// Try DATABASE_URL environment variable first (for Docker)
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		connStr = dbURL
	} else {
		// Fallback to app.conf (for local development)
		// dbtype, _ := beego.AppConfig.String("dbtype")
		dbhost, _ := beego.AppConfig.String("dbhost")
		dbport, _ := beego.AppConfig.String("dbport")
		dbuser, _ := beego.AppConfig.String("dbuser")
		dbpass, _ := beego.AppConfig.String("dbpass")
		dbname, _ := beego.AppConfig.String("dbname")
		dbsslmode, _ := beego.AppConfig.String("dbsslmode")

		connStr = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			dbhost, dbport, dbuser, dbpass, dbname, dbsslmode,
		)
	}

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Database connected successfully!")
	RunMigrations()
}

func RunMigrations() {
	driver, err := postgres.WithInstance(DB, &postgres.Config{})
	if err != nil {
		log.Fatal("Failed to create migration driver:", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal("Failed to initialize migrations:", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Failed to run migrations:", err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Fatal("Failed to get migration version:", err)
	}

	if err == migrate.ErrNilVersion {
		log.Println("No migrations applied yet")
	} else {
		log.Printf("Migrations applied successfully! Current version: %d (dirty: %v)", version, dirty)
	}
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}