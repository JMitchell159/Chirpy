package server

import (
	"sync/atomic"

	"github.com/JMitchell159/chirpy/internal/database"
)

type ApiConfig struct {
	FileServerHits atomic.Int32
	DB             *database.Queries
}
