package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

func (c Config) address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

type Server struct {
	config Config
	engine *gin.Engine
}

func New(config Config, engine *gin.Engine) *Server {
	return &Server{
		config: config,
		engine: engine,
	}
}

func (s *Server) Start() error {
	err := http.ListenAndServe(s.config.address(), s.engine)
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	return nil
}
