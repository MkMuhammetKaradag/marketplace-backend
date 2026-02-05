package domain

import "time"

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	UpdatedAt time.Time `json:"updated_at"`
}
