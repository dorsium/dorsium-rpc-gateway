package proxy

import (
	svc "github.com/dorsium/dorsium-rpc-gateway/internal/service/proxy"
	"github.com/gofiber/fiber/v2"
)

// Handler provides proxy HTTP handlers.
type Handler struct {
	service svc.Service
}

// NewHandler creates a proxy handler.
func NewHandler(svc svc.Service) *Handler {
	return &Handler{service: svc}
}

// RegisterRoutes registers proxy routes.
func (h *Handler) RegisterRoutes(r fiber.Router) {
	r.Post("/tx/send", h.sendTx)
	r.Get("/*", h.proxyGet)
}

func (h *Handler) sendTx(c *fiber.Ctx) error {
	data := c.Body()
	if len(data) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "empty body"})
	}
	resp, err := h.service.SendTx(data)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Send(resp)
}

func (h *Handler) proxyGet(c *fiber.Ctx) error {
	path := c.Params("*")
	resp, err := h.service.ProxyGet("/"+path, c.Context().QueryArgs().String())
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Send(resp)
}
