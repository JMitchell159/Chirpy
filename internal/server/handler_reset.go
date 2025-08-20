package server

import (
	"fmt"
	"net/http"
)

func (cfg *ApiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.FileServerHits.Store(0)

	if cfg.Platform != "dev" {
		w.WriteHeader(403)
		return
	}

	err := cfg.DB.ResetUsers(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		message := fmt.Sprintf("{\"error\":\"%s\"}", err)
		w.Write([]byte(message))
		return
	}

	w.WriteHeader(200)
}
