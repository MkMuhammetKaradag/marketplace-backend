package domain

import (
	"context"
	"time"
)

type Session struct {
	UserID string    `json:"id"`     // Kullanıcı ID'si
	Device string    `json:"device"` // Kullanıcı adı
	Ip     string    `json:"ip"`     // E-posta adresi
	Expiry time.Time `json:"expiry"` // Oturum sona erme zamanı

}
type SessionData struct {
	UserID    string    `json:"userID"`
	Username  string    `json:"username"`
	Role      UserRole  `json:"role"`
	Device    string    `json:"device"`
	Ip        string    `json:"ip"`
	CreatedAt time.Time `json:"createdAt"`
}

type SessionRepository interface {
	CreateSession(ctx context.Context, token string, duration time.Duration, data *SessionData) error
	DeleteSession(ctx context.Context, token string) error
	DeleteUserAllSession(ctx context.Context, token string) error
	GetSessionData(ctx context.Context, token string) (*SessionData, error)
	RefreshSession(ctx context.Context, token string, duration time.Duration) error
	GetTTL(ctx context.Context, token string) (time.Duration, error)
}
