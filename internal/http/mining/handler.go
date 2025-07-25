package mining

import (
	"github.com/dorsium/dorsium-rpc-gateway/internal/service/mining"
	"github.com/dorsium/dorsium-rpc-gateway/pkg/model"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Handler provides mining HTTP handlers.
type Handler struct {
	service  mining.Service
	validate *validator.Validate
}

// NewHandler creates a mining handler.
func NewHandler(svc mining.Service) *Handler {
	return &Handler{service: svc, validate: validator.New()}
}

// RegisterRoutes registers mining routes.
func (h *Handler) RegisterRoutes(r fiber.Router) {
	r.Post("/proof", h.submitProof)
	r.Get("/status", h.status)
}

func (h *Handler) submitProof(c *fiber.Ctx) error {
	var p model.Proof
	if err := c.BodyParser(&p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}
	if err := h.validate.Struct(p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "validation failed"})
	}
	if err := h.service.SubmitProof(p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true})
}

func (h *Handler) status(c *fiber.Ctx) error {
	st := h.service.GetStatus()
	return c.JSON(st)
}
