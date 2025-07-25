package main

import (
	"log"

	"github.com/dorsium/dorsium-rpc-gateway/internal/config"
	"github.com/dorsium/dorsium-rpc-gateway/internal/http"
	"github.com/dorsium/dorsium-rpc-gateway/internal/service"
)

func main() {
	cfg := config.New()
	svc := service.New()
	srv := http.NewServer(cfg, svc)

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
