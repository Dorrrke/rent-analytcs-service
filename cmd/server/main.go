package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Dorrrke/rent-analytcs-service/config"
	natshandler "github.com/Dorrrke/rent-analytcs-service/internal/handlers/nats"
	"github.com/Dorrrke/rent-analytcs-service/internal/infrastructure/logger"
	pgrepo "github.com/Dorrrke/rent-analytcs-service/internal/repository/postgres"
	"github.com/Dorrrke/rent-analytcs-service/internal/service"
	"go.uber.org/zap"
)

func main() {
	// --- Config -------------------------------------------------------------
	cfg, err := config.Load()
	if err != nil {
		panic("load config: " + err.Error())
	}

	// --- Logger -------------------------------------------------------------
	log, err := logger.New(cfg.Log.Level, cfg.App.Env)
	if err != nil {
		panic("init logger: " + err.Error())
	}
	defer log.Sync()

	log.Info("starting", zap.String("service", cfg.App.Name), zap.String("env", cfg.App.Env))

	// --- Postgres ------------------------------------------------------------
	// TODO: Init Postgres connection pool and pass it to repository
	// ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	// pool, err := pgifra.NewPool(ctx, cfg.Postgres, log)
	// cancel()
	// if err != nil {
	// 	log.Fatal("postgres init", zap.Error(err))
	// }
	// defer pool.Close()

	// --- NATS ----------------------------------------------------------------
	// TODO: Init NATS connection and pass it to repository
	// nc, err := natsinfra.NewConn(cfg.NATS, log)
	// if err != nil {
	// 	log.Fatal("nats init", zap.Error(err))
	// }
	// defer nc.Drain()

	// --- Wire up layers ------------------------------------------------------------
	repo := pgrepo.NewRepository(nil, log)
	svc := service.NewAnalyticsService(repo, log)
	handler := natshandler.New(svc, log)

	if err := handler.Subscribe(nil); err != nil {
		log.Fatal("nats subscribe", zap.Error(err))
	}

	defer handler.Unsubscribe()

	log.Info("analytics service is ready, waiting for events")

	// --- Graceful shutdown ------------------------------------------------------------
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Info("shutting down", zap.String("signal", sig.String()))
}
