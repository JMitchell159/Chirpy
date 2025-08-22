package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/JMitchell159/chirpy/internal/auth"
	"github.com/JMitchell159/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
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

	bearer_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	uuid, err := auth.ValidateJWT(bearer_token, cfg.SecretToken)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	if len(params.Body) > 140 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte("{\"error\":\"Chirp is too long\"}"))
		return
	}

	if strings.Contains(params.Body, "kerfuffle") || strings.Contains(params.Body, "sharbert") || strings.Contains(params.Body, "fornax") {
		splitBody := strings.Split(params.Body, " ")
		for i, word := range splitBody {
			lower := strings.ToLower(word)
			if lower == "kerfuffle" || lower == "sharbert" || lower == "fornax" {
				splitBody[i] = "****"
			}
		}
		params.Body = strings.Join(splitBody, " ")
	}

	chirp, err := cfg.DB.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   params.Body,
		UserID: uuid,
	})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		message := fmt.Sprintf("{\"error\":\"%s\"}", err)
		w.Write([]byte(message))
		return
	}

	respBody := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
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

func (cfg *ApiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author_id")
	direction := r.URL.Query().Get("sort")

	var chirps []database.Chirp

	if author == "" {
		temp, err := cfg.DB.GetChirps(r.Context())
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			message := fmt.Sprintf("{\"error\":\"%s\"}", err)
			w.Write([]byte(message))
			return
		}

		chirps = temp
	} else {
		author_id, err := uuid.Parse(author)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			message := fmt.Sprintf("{\"error\":\"%s\"}", err)
			w.Write([]byte(message))
			return
		}

		chirps, err = cfg.DB.GetChirpsByAuthor(r.Context(), author_id)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			message := fmt.Sprintf("{\"error\":\"%s\"}", err)
			w.Write([]byte(message))
			return
		}
	}

	result := make([]Chirp, len(chirps))
	for i, chirp := range chirps {
		result[i] = Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}

	if direction == "desc" {
		sort.Slice(result, func(i, j int) bool { return result[i].CreatedAt.Compare(result[j].CreatedAt) > 0 })
	}

	dat, err := json.Marshal(result)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		message := fmt.Sprintf("{\"error\":\"%s\"}", err)
		w.Write([]byte(message))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}

func (cfg *ApiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	arg := r.PathValue("chirp_id")
	id, err := uuid.Parse(arg)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		message := fmt.Sprintf("{\"error\":\"%s\"}", err)
		w.Write([]byte(message))
		return
	}

	chirp, err := cfg.DB.GetChirp(r.Context(), id)
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

	result := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	dat, err := json.Marshal(result)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		message := fmt.Sprintf("{\"error\":\"%s\"}", err)
		w.Write([]byte(message))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}

func (cfg *ApiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	arg := r.PathValue("chirp_id")
	id, err := uuid.Parse(arg)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		message := fmt.Sprintf("{\"error\":\"%s\"}", err)
		w.Write([]byte(message))
		return
	}

	user, err := cfg.DB.GetUserByChirp(r.Context(), id)
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

	jwt, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	uuid, err := auth.ValidateJWT(jwt, cfg.SecretToken)
	if err != nil {
		w.WriteHeader(403)
		return
	} else if uuid != user.ID {
		w.WriteHeader(403)
		return
	}

	err = cfg.DB.DeleteChirp(r.Context(), id)
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
