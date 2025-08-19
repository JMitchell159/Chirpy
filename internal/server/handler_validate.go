package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func handlerValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVal struct {
		CleanedBody string `json:"cleaned_body"`
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

	respBody := returnVal{
		CleanedBody: params.Body,
	}

	if strings.Contains(respBody.CleanedBody, "kerfuffle") || strings.Contains(respBody.CleanedBody, "sharbert") || strings.Contains(respBody.CleanedBody, "fornax") {
		splitBody := strings.Split(respBody.CleanedBody, " ")
		for i, word := range splitBody {
			lower := strings.ToLower(word)
			if lower == "kerfuffle" || lower == "sharbert" || lower == "fornax" {
				splitBody[i] = "****"
			}
		}
		respBody.CleanedBody = strings.Join(splitBody, " ")
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
	w.WriteHeader(200)
	w.Write(dat)
}
