package main

import "github.com/JMitchell159/chirpy/internal/server"

func main() {
	apiCfg := &server.ApiConfig{}
	apiCfg.FileServerHits.Store(0)
	server := apiCfg.CreateServer()
	server.ListenAndServe()
}
