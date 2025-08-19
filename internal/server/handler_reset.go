package server

import "net/http"

func (cfg *ApiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.FileServerHits.Store(0)
	w.WriteHeader(200)
}
