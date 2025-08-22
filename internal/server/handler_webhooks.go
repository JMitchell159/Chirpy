package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JMitchell159/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) handlerPolka(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		message := fmt.Sprintf("{\"error\":\"%s\"}", err)
		w.Write([]byte(message))
		return
	} else if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		w.WriteHeader(401)
		return
	} else if apiKey != cfg.PolkaKey {
		w.WriteHeader(401)
		return
	}

	err = cfg.DB.GetMembership(r.Context(), params.Data.UserID)
	if err == sql.ErrNoRows {
		w.WriteHeader(404)
		return
	} else if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		message := fmt.Sprintf("{\"error\":\"%s\"}", err)
		w.Write([]byte(message))
		return
	}

	w.WriteHeader(204)
}
