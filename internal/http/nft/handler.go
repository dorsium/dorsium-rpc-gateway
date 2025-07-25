package nft

import (
	"github.com/dorsium/dorsium-rpc-gateway/internal/service/nft"
	"github.com/dorsium/dorsium-rpc-gateway/pkg/model"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Handler provides NFT HTTP handlers.
type Handler struct {
	service  nft.Service
	validate *validator.Validate
}

// NewHandler creates an NFT handler.
func NewHandler(svc nft.Service) *Handler {
	return &Handler{service: svc, validate: validator.New()}
}

// RegisterRoutes registers NFT routes.
func (h *Handler) RegisterRoutes(r fiber.Router) {
	r.Get("/:id", h.getMetadata)
	r.Post("/mint", h.mint)
	r.Get("/:id/image", h.getImage)
}

func (h *Handler) getMetadata(c *fiber.Ctx) error {
	id := c.Params("id")
	meta, err := h.service.GetMetadata(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(meta)
}

func (h *Handler) mint(c *fiber.Ctx) error {
	var req model.MintRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "validation failed"})
	}
	meta, err := h.service.MintNFT(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(meta)
}

func (h *Handler) getImage(c *fiber.Ctx) error {
	id := c.Params("id")
	data, ct, err := h.service.GetImage(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	c.Set("Content-Type", ct)
	return c.Send(data)
}
