package http

import (
	"github.com/dorsium/dorsium-rpc-gateway/internal/config"
	dapphttp "github.com/dorsium/dorsium-rpc-gateway/internal/http/dapp"
	mininghttp "github.com/dorsium/dorsium-rpc-gateway/internal/http/mining"
	nfthttp "github.com/dorsium/dorsium-rpc-gateway/internal/http/nft"
	nodehttp "github.com/dorsium/dorsium-rpc-gateway/internal/http/node"
	validatorhttp "github.com/dorsium/dorsium-rpc-gateway/internal/http/validator"
	wallethttp "github.com/dorsium/dorsium-rpc-gateway/internal/http/wallet"
	nftrepo "github.com/dorsium/dorsium-rpc-gateway/internal/repository/nft"
	noderepo "github.com/dorsium/dorsium-rpc-gateway/internal/repository/node"
	validatorrepo "github.com/dorsium/dorsium-rpc-gateway/internal/repository/validator"
	walletrepo "github.com/dorsium/dorsium-rpc-gateway/internal/repository/wallet"
	"github.com/dorsium/dorsium-rpc-gateway/internal/service"
	dappservice "github.com/dorsium/dorsium-rpc-gateway/internal/service/dapp"
	miningservice "github.com/dorsium/dorsium-rpc-gateway/internal/service/mining"
	nftservice "github.com/dorsium/dorsium-rpc-gateway/internal/service/nft"
	nodeservice "github.com/dorsium/dorsium-rpc-gateway/internal/service/node"
	validatorservice "github.com/dorsium/dorsium-rpc-gateway/internal/service/validator"
	walletservice "github.com/dorsium/dorsium-rpc-gateway/internal/service/wallet"
	"github.com/dorsium/dorsium-rpc-gateway/pkg/model"
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

	// wallet routes
	repo := walletrepo.New()
	svc := walletservice.New(repo)
	handler := wallethttp.NewHandler(svc)
	walletGroup := api.Group("/wallet")
	handler.RegisterRoutes(walletGroup)

	// nft routes
	nRepo := nftrepo.New()
	nSvc := nftservice.New(nRepo, nftservice.NewDummyMintHandler())
	nHandler := nfthttp.NewHandler(nSvc)
	nftGroup := api.Group("/nft")
	nHandler.RegisterRoutes(nftGroup)

	// dapp routes
	dSvc := dappservice.New(nRepo)
	dHandler := dapphttp.NewHandler(dSvc)
	dappGroup := api.Group("/dapp")
	dHandler.RegisterRoutes(dappGroup)

	// validator routes
	vRepo := validatorrepo.New()
	vSvc := validatorservice.New(vRepo)
	vHandler := validatorhttp.NewHandler(vSvc)
	vGroup := api.Group("/validator")
	vHandler.RegisterRoutes(vGroup)

	// node routes
	nodeRepo := noderepo.New()
	nodeSvc := nodeservice.New(nodeRepo)
	nodeHandler := nodehttp.NewHandler(nodeSvc)
	nodeGroup := api.Group("/node")
	nodeHandler.RegisterRoutes(nodeGroup)

	// mining routes
	mVerifier := miningservice.NewDummyVerifier()
	mSvc := miningservice.New(mVerifier, model.MiningStatus{
		Mode:       "pow",
		Difficulty: 1,
		Challenge:  "0000",
	})
	mHandler := mininghttp.NewHandler(mSvc)
	miningGroup := api.Group("/mining")
	mHandler.RegisterRoutes(miningGroup)

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
