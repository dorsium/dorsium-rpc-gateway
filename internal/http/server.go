package http

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/dorsium/dorsium-rpc-gateway/internal/config"
	adminhttp "github.com/dorsium/dorsium-rpc-gateway/internal/http/admin"
	dapphttp "github.com/dorsium/dorsium-rpc-gateway/internal/http/dapp"
	mininghttp "github.com/dorsium/dorsium-rpc-gateway/internal/http/mining"
	nfthttp "github.com/dorsium/dorsium-rpc-gateway/internal/http/nft"
	nodehttp "github.com/dorsium/dorsium-rpc-gateway/internal/http/node"
	proxyhttp "github.com/dorsium/dorsium-rpc-gateway/internal/http/proxy"
	validatorhttp "github.com/dorsium/dorsium-rpc-gateway/internal/http/validator"
	wallethttp "github.com/dorsium/dorsium-rpc-gateway/internal/http/wallet"
	nftrepo "github.com/dorsium/dorsium-rpc-gateway/internal/repository/nft"
	noderepo "github.com/dorsium/dorsium-rpc-gateway/internal/repository/node"
	proxyrepo "github.com/dorsium/dorsium-rpc-gateway/internal/repository/proxy"
	validatorrepo "github.com/dorsium/dorsium-rpc-gateway/internal/repository/validator"
	walletrepo "github.com/dorsium/dorsium-rpc-gateway/internal/repository/wallet"
	"github.com/dorsium/dorsium-rpc-gateway/internal/service"
	adminservice "github.com/dorsium/dorsium-rpc-gateway/internal/service/admin"
	dappservice "github.com/dorsium/dorsium-rpc-gateway/internal/service/dapp"
	miningservice "github.com/dorsium/dorsium-rpc-gateway/internal/service/mining"
	nftservice "github.com/dorsium/dorsium-rpc-gateway/internal/service/nft"
	nodeservice "github.com/dorsium/dorsium-rpc-gateway/internal/service/node"
	proxyservice "github.com/dorsium/dorsium-rpc-gateway/internal/service/proxy"
	validatorservice "github.com/dorsium/dorsium-rpc-gateway/internal/service/validator"
	walletservice "github.com/dorsium/dorsium-rpc-gateway/internal/service/wallet"
	"github.com/dorsium/dorsium-rpc-gateway/pkg/model"
	"github.com/gofiber/fiber/v2"
	"github.com/shirou/gopsutil/v3/cpu"
)

const maxCallsEntries = 100

// Server holds dependencies for HTTP handlers and routes.
type Server struct {
	cfg     *config.Config
	app     *fiber.App
	service service.Service
	start   time.Time
	mu      sync.Mutex
	calls   map[string]int
}

// NewServer creates a Server with all dependencies wired.
func NewServer(cfg *config.Config, svc service.Service) *Server {
	app := fiber.New()
	s := &Server{cfg: cfg, app: app, service: svc, start: time.Now(), calls: make(map[string]int)}
	app.Use(func(c *fiber.Ctx) error {
		err := c.Next()
		route := c.Route()
		ep := "unknown"
		if route != nil && route.Path != "" && route.Path != "/" {
			ep = route.Path
		}
		s.mu.Lock()
		if _, ok := s.calls[ep]; !ok && len(s.calls) >= maxCallsEntries {
			ep = "unknown"
		}
		s.calls[ep]++
		s.mu.Unlock()
		return err
	})
	return s
}

// RegisterRoutes sets up HTTP routes.
func (s *Server) RegisterRoutes() {
	s.app.Get("/status", s.status)
	s.app.Get("/metrics", s.metrics)
	s.app.Get("/mode", s.mode)

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
	nSvc := nftservice.New(nRepo, nftservice.NewDummyMintHandler(), s.cfg.MaxResponseSize)
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

	// proxy routes
	pRepo := proxyrepo.New(s.cfg.NodeRPC, s.cfg.MaxResponseSize)
	pSvc := proxyservice.New(pRepo)
	pHandler := proxyhttp.NewHandler(pSvc)
	proxyGroup := api.Group("/proxy")
	pHandler.RegisterRoutes(proxyGroup)

	// admin routes
	aSvc := adminservice.New(nodeRepo, vRepo)
	aHandler := adminhttp.NewHandler(aSvc, s.cfg.AdminToken)
	adminGroup := s.app.Group("/admin")
	aHandler.RegisterRoutes(adminGroup)

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

func (s *Server) status(c *fiber.Ctx) error {
	uptime := time.Since(s.start).Seconds()
	return c.JSON(fiber.Map{
		"version": s.cfg.Version,
		"uptime":  uptime,
		"health":  "ok",
	})
}

func (s *Server) metrics(c *fiber.Ctx) error {
	if s.cfg.DisableMetrics {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "metrics disabled"})
	}
	cpuPerc, _ := cpu.Percent(0, false)
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	metrics := fmt.Sprintf("# HELP cpu_usage_percent CPU usage percent\n# TYPE cpu_usage_percent gauge\ncpu_usage_percent %.2f\n# HELP ram_usage_bytes RAM usage bytes\n# TYPE ram_usage_bytes gauge\nram_usage_bytes %d\n", cpuPerc[0], m.Alloc)
	s.mu.Lock()
	for ep, count := range s.calls {
		metrics += fmt.Sprintf("endpoint_calls_total{endpoint=\"%s\"} %d\n", ep, count)
	}
	s.mu.Unlock()
	return c.Type("text/plain").SendString(metrics)
}

func (s *Server) mode(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"mode": s.cfg.Mode})
}

// Start runs the HTTP server.
func (s *Server) Start() error {
	s.RegisterRoutes()
	return s.app.Listen(s.cfg.Address)
}

// App exposes the underlying Fiber app for testing.
func (s *Server) App() *fiber.App {
	return s.app
}

// Calls returns a snapshot of recorded endpoint calls.
func (s *Server) Calls() map[string]int {
	s.mu.Lock()
	defer s.mu.Unlock()
	cp := make(map[string]int, len(s.calls))
	for k, v := range s.calls {
		cp[k] = v
	}
	return cp
}
