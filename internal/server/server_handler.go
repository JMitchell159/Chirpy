package server

import "net/http"

func (cfg *ApiConfig) CreateServer() *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", handlerHealth)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	mux.HandleFunc("POST /api/chirps", cfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", cfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirp_id}", cfg.handlerGetChirp)
	mux.HandleFunc("POST /api/users", cfg.handlerUsers)
	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}
