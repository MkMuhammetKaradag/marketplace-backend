package domain

import "time"

type UserRole string

const (
	Buyer  UserRole = "Buyer"
	Seller UserRole = "Seller"
	Admin  UserRole = "Admin"
)

type User struct {
	ID                  string    `json:"id"`
	Username            string    `json:"username"`
	Email               string    `json:"email"`
	Password            string    `json:"password"`
	Permissions         int64     `json:"permissions"`
	ActivationCode      string    `json:"activationCode"`
	ActivationExpiry    time.Time `json:"activationExpiry"`
	FailedLoginAttempts int       `json:"failedLoginAttempts"`
	AccountLocked       bool      `json:"accountLocked"`
	LockUntil           time.Time `json:"lockUntil"`
}
