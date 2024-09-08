package server

import (
	"fmt"
	"net"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/csatib02/kube-pod-autocomplete/internal/config"
	"github.com/csatib02/kube-pod-autocomplete/internal/handlers"
)

type Server struct {
	router *gin.Engine
	config *config.Config
}

func New(config *config.Config) (*Server, error) {
	// Init gin router and set up middlewares
	// NOTE: Add Auth middleware if required
	router := gin.New()
	router.Use(gin.Recovery(), cors.Default())
	gin.SetMode(config.Mode)

	if config.LogServerAddress != "" {
		writer, err := net.Dial("udp", config.LogServerAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to log server: %w", err)
		}

		router.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: writer}))
	} else {
		router.Use(gin.Logger())
	}

	// Add trusted proxies e.g. when running KPA behind a reverse proxy or a load balancer
	if err := router.SetTrustedProxies(config.TrustedProxies); err != nil {
		return nil, fmt.Errorf("failed to set trusted proxies: %w", err)
	}

	handlers.SetupRoutes(router)

	return &Server{
		router: router,
		config: config,
	}, nil
}

func (s *Server) Run() error {
	return s.router.Run(s.config.ListenAddress)
}
