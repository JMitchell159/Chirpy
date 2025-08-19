package server

import (
	"fmt"
	"net/http"
)

func (cfg *ApiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	message := fmt.Sprintf("<html>\n\t<body>\n\t\t<h1>Welcome, Chirpy Admin</h1>\n\t\t<p>Chirpy has been visited %d times!</p>\n\t</body>\n</html>", cfg.FileServerHits.Load())
	w.Write([]byte(message))
}
