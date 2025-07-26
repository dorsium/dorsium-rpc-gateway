package dapp

import (
	dappsvc "github.com/dorsium/dorsium-rpc-gateway/internal/service/dapp"
	"github.com/dorsium/dorsium-rpc-gateway/pkg/model"
	"github.com/dorsium/dorsium-rpc-gateway/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Handler provides DAPP HTTP handlers.
type Handler struct {
	service  dappsvc.Service
	validate *validator.Validate
}

// NewHandler creates a DAPP handler.
func NewHandler(svc dappsvc.Service) *Handler {
	return &Handler{service: svc, validate: validator.New()}
}

// RegisterRoutes registers DAPP routes.
func (h *Handler) RegisterRoutes(r fiber.Router) {
	r.Get("/config", h.getConfig)
	r.Get("/nft/:id/verify", h.verifyNFT)
	r.Post("/connect/verify-wallet", h.verifyWallet)
	r.Get("/:address/permissions", h.permissions)
}

func (h *Handler) getConfig(c *fiber.Ctx) error {
	cfg := h.service.GetConfig()
	return c.JSON(cfg)
}

func (h *Handler) verifyNFT(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.VerifyNFT(id); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"valid": true})
}

func (h *Handler) verifyWallet(c *fiber.Ctx) error {
	var req model.WalletVerifyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "validation failed"})
	}
	ok, err := h.service.VerifyWallet(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"verified": ok})
}

func (h *Handler) permissions(c *fiber.Ctx) error {
	addr := c.Params("address")
	if !utils.IsValidAddress(addr) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid address"})
	}
	p, err := h.service.GetPermissions(addr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(p)
}
