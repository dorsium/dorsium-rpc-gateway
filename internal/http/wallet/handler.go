package wallet

import (
	"regexp"

	"github.com/dorsium/dorsium-rpc-gateway/internal/service/wallet"
	"github.com/gofiber/fiber/v2"
)

// Handler provides wallet HTTP handlers.
type Handler struct {
	service wallet.Service
}

// NewHandler creates a wallet handler.
func NewHandler(svc wallet.Service) *Handler {
	return &Handler{service: svc}
}

// RegisterRoutes registers wallet routes.
func (h *Handler) RegisterRoutes(r fiber.Router) {
	r.Get("/:address", h.getInfo)
	r.Get("/:address/transactions", h.getTransactions)
	r.Get("/:address/nfts", h.getNFTs)
}

func (h *Handler) getInfo(c *fiber.Ctx) error {
	addr := c.Params("address")
	if !isValidAddress(addr) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid address"})
	}
	info, err := h.service.GetInfo(addr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(info)
}

func (h *Handler) getTransactions(c *fiber.Ctx) error {
	addr := c.Params("address")
	if !isValidAddress(addr) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid address"})
	}
	txs, err := h.service.GetTransactions(addr, 50)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(txs)
}

func (h *Handler) getNFTs(c *fiber.Ctx) error {
	addr := c.Params("address")
	if !isValidAddress(addr) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid address"})
	}
	nfts, err := h.service.GetNFTs(addr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(nfts)
}

var (
	hexRegex    = regexp.MustCompile(`^(0x)?[0-9a-fA-F]{40}$`)
	bech32Regex = regexp.MustCompile(`^[a-z0-9]{1,83}1[qpzry9x8gf2tvdw0s3jn54khce6mua7l]{38}$`)
)

func isValidAddress(addr string) bool {
	return hexRegex.MatchString(addr) || bech32Regex.MatchString(addr)
}
