package http

import (
	"github.com/dorsium/dorsium-rpc-gateway/internal/config"
	"github.com/dorsium/dorsium-rpc-gateway/internal/service"
	"github.com/gofiber/fiber/v2"
)

// Server holds dependencies for HTTP handlers and routes.
type Server struct {
	cfg     *config.Config
	app     *fiber.App
	service service.Service
}

// NewServer creates a Server with all dependencies wired.
func NewServer(cfg *config.Config, svc service.Service) *Server {
	app := fiber.New()
	return &Server{cfg: cfg, app: app, service: svc}
}

// RegisterRoutes sets up HTTP routes.
func (s *Server) RegisterRoutes() {
	api := s.app.Group("/api")
	api.Get("/ping", s.ping)

	// Placeholders for future endpoints
	for i := 0; i < 25; i++ {
		path := "/placeholder" + string(rune('A'+i))
		api.Get(path, s.placeholder)
	}
}

func (s *Server) ping(c *fiber.Ctx) error {
	msg := s.service.Ping()
	return c.JSON(fiber.Map{"message": msg})
}

func (s *Server) placeholder(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "not implemented"})
}

// Start runs the HTTP server.
func (s *Server) Start() error {
	s.RegisterRoutes()
	return s.app.Listen(s.cfg.Address)
}
