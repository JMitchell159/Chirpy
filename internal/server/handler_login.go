package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JMitchell159/chirpy/internal/auth"
)

func (cfg *ApiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
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

	user, err := cfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err == sql.ErrNoRows {
		w.Header().Set("Content-Type", "plain/text; charset=utf-8")
		w.WriteHeader(401)
		w.Write([]byte("Incorrect email or password"))
		return
	} else if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		message := fmt.Sprintf("{\"error\":\"%s\"}", err)
		w.Write([]byte(message))
		return
	} else if err = auth.CheckPasswordHash(params.Password, user.HashedPassword); err != nil {
		w.Header().Set("Content-Type", "plain/text; charset=utf-8")
		w.WriteHeader(401)
		w.Write([]byte("Incorrect email or password"))
		return
	}

	result := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
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
