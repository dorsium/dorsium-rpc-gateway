package validator

import (
	"regexp"
	"strconv"

	validatorsvc "github.com/dorsium/dorsium-rpc-gateway/internal/service/validator"
	"github.com/gofiber/fiber/v2"
)

// Handler provides validator HTTP handlers.
type Handler struct {
	service validatorsvc.Service
}

// NewHandler creates a validator handler.
func NewHandler(svc validatorsvc.Service) *Handler { return &Handler{service: svc} }

// RegisterRoutes registers validator routes.
func (h *Handler) RegisterRoutes(r fiber.Router) {
	r.Get("/list", h.list)
	r.Get("/:address/status", h.status)
	r.Get("/:address/profile", h.profile)
}

func (h *Handler) status(c *fiber.Ctx) error {
	addr := c.Params("address")
	if !isValidAddress(addr) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid address"})
	}
	st, err := h.service.GetStatus(addr)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(st)
}

func (h *Handler) profile(c *fiber.Ctx) error {
	addr := c.Params("address")
	if !isValidAddress(addr) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid address"})
	}
	prof, err := h.service.GetProfile(addr)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(prof)
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

var (
	hexRegex    = regexp.MustCompile(`^(0x)?[0-9a-fA-F]{40}$`)
	bech32Regex = regexp.MustCompile(`^[a-z0-9]{1,83}1[qpzry9x8gf2tvdw0s3jn54khce6mua7l]{38}$`)
)

func isValidAddress(addr string) bool {
	return hexRegex.MatchString(addr) || bech32Regex.MatchString(addr)
}
