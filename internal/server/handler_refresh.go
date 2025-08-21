package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/JMitchell159/chirpy/internal/auth"
)

func (cfg *ApiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	bearer_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		message := fmt.Sprintf("{\"error\":\"%s\"}", err)
		w.Write([]byte(message))
		return
	}

	refresh_token, err := cfg.DB.GetRefreshToken(r.Context(), bearer_token)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		message := fmt.Sprintf("{\"error\":\"%s\"}", err)
		w.Write([]byte(message))
		return
	} else if refresh_token.RevokedAt.Valid {
		w.WriteHeader(401)
		return
	}

	user, err := cfg.DB.GetUserFromRefreshToken(r.Context(), bearer_token)
	if err == sql.ErrNoRows {
		w.WriteHeader(401)
		return
	} else if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		message := fmt.Sprintf("{\"error\":\"%s\"}", err)
		w.Write([]byte(message))
		return
	}

	jwt, err := auth.MakeJWT(user.ID, cfg.SecretToken, time.Hour)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		message := fmt.Sprintf("{\"error\":\"%s\"}", err)
		w.Write([]byte(message))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	message := fmt.Sprintf("{\"token\": \"%s\"}", jwt)
	w.Write([]byte(message))
}
