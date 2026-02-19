package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Database connection
	connStr := "host=localhost port=5432 user=devops password=devops123 dbname=devops_db sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer db.Close()

	// Create admin user
	adminPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	_, err = db.Exec(`
		INSERT INTO accounts_user (email, password, name, role, is_active, is_staff, is_superuser)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (email) DO UPDATE SET
			password = EXCLUDED.password,
			name = EXCLUDED.name,
			role = EXCLUDED.role,
			is_active = EXCLUDED.is_active,
			is_staff = EXCLUDED.is_staff,
			is_superuser = EXCLUDED.is_superuser
	`, "admin@example.com", string(adminPassword), "Admin User", "admin", true, true, false)
	
	if err != nil {
		log.Println("Error creating admin:", err)
	} else {
		fmt.Println("✓ Admin user created (email: admin@example.com, password: admin123)")
	}

	// Create regular user
	userPassword, _ := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)
	_, err = db.Exec(`
		INSERT INTO accounts_user (email, password, name, role, is_active, is_staff, is_superuser)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (email) DO UPDATE SET
			password = EXCLUDED.password,
			name = EXCLUDED.name,
			role = EXCLUDED.role,
			is_active = EXCLUDED.is_active,
			is_staff = EXCLUDED.is_staff,
			is_superuser = EXCLUDED.is_superuser
	`, "user@example.com", string(userPassword), "Regular User", "user", true, false, false)
	
	if err != nil {
		log.Println("Error creating user:", err)
	} else {
		fmt.Println("✓ Regular user created (email: user@example.com, password: user123)")
	}

	fmt.Println("\nTest users are ready!")
}
