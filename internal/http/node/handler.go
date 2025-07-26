package node

import (
	"strconv"

	nodesvc "github.com/dorsium/dorsium-rpc-gateway/internal/service/node"
	"github.com/dorsium/dorsium-rpc-gateway/pkg/model"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Handler provides node HTTP handlers.
type Handler struct {
	service  nodesvc.Service
	validate *validator.Validate
}

// NewHandler creates a node handler.
func NewHandler(svc nodesvc.Service) *Handler {
	return &Handler{service: svc, validate: validator.New()}
}

// RegisterRoutes registers node routes.
func (h *Handler) RegisterRoutes(r fiber.Router) {
	r.Get("/list", h.list)
	r.Get("/:id/status", h.status)
	r.Post("/ping", h.ping)
	r.Get("/:id/profile", h.profile)
	r.Get("/:id/metrics", h.metrics)
}

func (h *Handler) status(c *fiber.Ctx) error {
	id := c.Params("id")
	st, err := h.service.GetStatus(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(st)
}

func (h *Handler) ping(c *fiber.Ctx) error {
	var p model.NodePing
	if err := c.BodyParser(&p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}
	if err := h.validate.Struct(p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "validation failed"})
	}
	if err := h.service.Ping(p); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(fiber.Map{"success": true})
}

func (h *Handler) profile(c *fiber.Ctx) error {
	id := c.Params("id")
	prof, err := h.service.GetProfile(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(prof)
}

func (h *Handler) metrics(c *fiber.Ctx) error {
	id := c.Params("id")
	m, err := h.service.GetMetrics(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(m)
}

func (h *Handler) list(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	res, err := h.service.List(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}
