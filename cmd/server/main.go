package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/lukelinux79/casaos-native-panel/internal/api"
	"github.com/lukelinux79/casaos-native-panel/internal/config"
	"github.com/lukelinux79/casaos-native-panel/internal/runtime"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	rt, err := runtime.New(cfg.Mode)
	if err != nil {
		log.Fatalf("failed to initialize runtime: %v", err)
	}

	srv := api.NewServer(rt)
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	log.Printf("Starting CasaOS Native Panel on %s in %s mode", addr, cfg.Mode)

	if err := http.ListenAndServe(addr, srv.Handler()); err != nil && err != http.ErrServerClosed {
		log.Printf("server error: %v", err)
		os.Exit(1)
	}
}
