package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/csatib02/kube-pod-autocomplete/internal/config"
	"github.com/csatib02/kube-pod-autocomplete/internal/server"
	"github.com/csatib02/kube-pod-autocomplete/pkg/log"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		slog.Error(fmt.Errorf("failed to load configuration: %w", err).Error())
		os.Exit(1)
	}

	log.InitLogger(config)

	server, err := server.New(config)
	if err != nil {
		slog.Error(fmt.Errorf("failed to create server: %w", err).Error())
		os.Exit(1)
	}

	if err := server.Run(); err != nil {
		slog.Error(fmt.Errorf("failed to start server: %w", err).Error())
		os.Exit(1)
	}
}
