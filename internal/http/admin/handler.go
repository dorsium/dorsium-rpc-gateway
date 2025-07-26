package admin

import (
	adminsvc "github.com/dorsium/dorsium-rpc-gateway/internal/service/admin"
	"github.com/gofiber/fiber/v2"
)

// Handler provides admin HTTP handlers.
type Handler struct {
	service adminsvc.Service
	token   string
}

// NewHandler creates an admin handler.
func NewHandler(svc adminsvc.Service, token string) *Handler {
	return &Handler{service: svc, token: token}
}

// RegisterRoutes registers admin routes.
func (h *Handler) RegisterRoutes(r fiber.Router) {
	r.Use(h.auth)
	r.Post("/broadcast", h.broadcast)
	r.Get("/logs", h.logs)
}

func (h *Handler) auth(c *fiber.Ctx) error {
	key := c.Get("X-Admin-Token")
	if key == "" {
		key = c.Get("X-Admin-Key")
	}
	if key != h.token {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	return c.Next()
}

func (h *Handler) broadcast(c *fiber.Ctx) error {
	var req struct {
		Message string `json:"message"`
	}
	if err := c.BodyParser(&req); err != nil || req.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "message required"})
	}
	if err := h.service.Broadcast(req.Message); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true})
}

func (h *Handler) logs(c *fiber.Ctx) error {
	logs := h.service.GetLogs()
	return c.JSON(fiber.Map{"logs": logs})
}
