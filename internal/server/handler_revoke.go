package server

import (
	"fmt"
	"net/http"

	"github.com/JMitchell159/chirpy/internal/auth"
)

func (cfg *ApiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	bearer_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		message := fmt.Sprintf("{\"error\":\"%s\"}", err)
		w.Write([]byte(message))
		return
	}

	err = cfg.DB.RevokeRefreshToken(r.Context(), bearer_token)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		message := fmt.Sprintf("{\"error\":\"%s\"}", err)
		w.Write([]byte(message))
		return
	}

	w.WriteHeader(204)
}
