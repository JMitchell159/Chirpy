package server

import (
	"sync/atomic"
	"time"

	"github.com/JMitchell159/chirpy/internal/database"
	"github.com/google/uuid"
)

type ApiConfig struct {
	FileServerHits atomic.Int32
	DB             *database.Queries
	Platform       string
	SecretToken    string
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}
