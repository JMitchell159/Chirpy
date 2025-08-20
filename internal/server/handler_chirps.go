package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/JMitchell159/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	type returnVal struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
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
	}

	if len(params.Body) > 140 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte("{\"error\":\"Chirp is too long\"}"))
		return
	}

	chirp, err := cfg.DB.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   params.Body,
		UserID: params.UserID,
	})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		message := fmt.Sprintf("{\"error\":\"%s\"}", err)
		w.Write([]byte(message))
		return
	}

	respBody := returnVal{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	if strings.Contains(respBody.Body, "kerfuffle") || strings.Contains(respBody.Body, "sharbert") || strings.Contains(respBody.Body, "fornax") {
		splitBody := strings.Split(respBody.Body, " ")
		for i, word := range splitBody {
			lower := strings.ToLower(word)
			if lower == "kerfuffle" || lower == "sharbert" || lower == "fornax" {
				splitBody[i] = "****"
			}
		}
		respBody.Body = strings.Join(splitBody, " ")
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		message := fmt.Sprintf("{\"error\":\"%s\"}", err)
		w.Write([]byte(message))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(dat)
}
