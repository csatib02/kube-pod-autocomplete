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
	// TODO: Add Auth middleware
	router := gin.New()
	router.Use(gin.Recovery(), cors.Default())

	if config.LogServerAddress != "" {
		writer, err := net.Dial("udp", config.LogServerAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to log server: %w", err)
		}

		router.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: writer}))
	} else {
		router.Use(gin.Logger())
	}

	// Set the mode of the gin router, default is debug
	if config.Mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	// Add trusted proxies e.g. when running KPA behind a reverse proxy or a load balancer
	if err := router.SetTrustedProxies(config.TrustedProxies); err != nil {
		return nil, fmt.Errorf("failed to set trusted proxies: %w", err)
	}

	handlers.SetupRouter(router)

	server := &Server{
		router: router,
		config: config,
	}

	return server, nil
}

func (s *Server) Run() error {
	return s.router.Run(s.config.ListenAddress)
}
