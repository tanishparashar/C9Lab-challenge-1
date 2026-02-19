package models

import (
	"database/sql"
	"errors"
	"time"

	"devops-api/database"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id           int        `json:"id"`
	Password     string     `json:"-"`
	LastLogin    *time.Time `json:"last_login,omitempty"`
	IsSuperuser  bool       `json:"is_superuser"`
	Email        string     `json:"email"`
	Name         string     `json:"name"`
	Role         string     `json:"role"`
	IsActive     bool       `json:"is_active"`
	IsStaff      bool       `json:"is_staff"`
}

type UserDTO struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

// GetUserByEmail retrieves a user by email
func GetUserByEmail(email string) (*User, error) {
	user := &User{}
	query := `SELECT id, password, last_login, is_superuser, email, name, role, is_active, is_staff 
	          FROM accounts_user WHERE email = $1`
	
	err := database.DB.QueryRow(query, email).Scan(
		&user.Id, &user.Password, &user.LastLogin, &user.IsSuperuser,
		&user.Email, &user.Name, &user.Role, &user.IsActive, &user.IsStaff,
	)
	
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

// GetUserById retrieves a user by ID
func GetUserById(id int) (*User, error) {
	user := &User{}
	query := `SELECT id, password, last_login, is_superuser, email, name, role, is_active, is_staff 
	          FROM accounts_user WHERE id = $1`
	
	err := database.DB.QueryRow(query, id).Scan(
		&user.Id, &user.Password, &user.LastLogin, &user.IsSuperuser,
		&user.Email, &user.Name, &user.Role, &user.IsActive, &user.IsStaff,
	)
	
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

// CreateUser creates a new user
func CreateUser(email, password, name, role string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO accounts_user (email, password, name, role, is_active, is_staff, is_superuser) 
	          VALUES ($1, $2, $3, $4, true, false, false) RETURNING id`
	
	user := &User{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
		Role:     role,
		IsActive: true,
		IsStaff:  false,
		IsSuperuser: false,
	}
	
	err = database.DB.QueryRow(query, email, string(hashedPassword), name, role).Scan(&user.Id)
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

// UpdateLastLogin updates the last login timestamp
func UpdateLastLogin(userId int) error {
	query := `UPDATE accounts_user SET last_login = $1 WHERE id = $2`
	_, err := database.DB.Exec(query, time.Now(), userId)
	return err
}

// CheckPassword verifies the password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// GetTotalUsers returns the total number of users
func GetTotalUsers() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM accounts_user`
	err := database.DB.QueryRow(query).Scan(&count)
	return count, err
}

// GetActiveUsers returns the number of active users
func GetActiveUsers() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM accounts_user WHERE is_active = true`
	err := database.DB.QueryRow(query).Scan(&count)
	return count, err
}

// ToDTO converts User to UserDTO (safe for API responses)
func (u *User) ToDTO() UserDTO {
	return UserDTO{
		Id:    u.Id,
		Email: u.Email,
		Name:  u.Name,
		Role:  u.Role,
	}
}

