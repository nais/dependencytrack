package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/nais/dependencytrack/internal/dependencytrack"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"os/signal"
	"syscall"
	"time"
)

const collectMetricsInterval = 3 * time.Minute

type Config struct {
	ListenAddress      string `envconfig:"LISTEN_ADDRESS" default:"localhost:8000"`
	DependencytrackUrl string `envconfig:"DEPENDENCYTRACK_URL" default:"http://localhost:8080"`
	Username           string `envconfig:"USERNAME" default:"admin"`
	Password           string `envconfig:"PASSWORD"`
}

// TODO: review waitgroups, context and signal handling
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx, signalStop := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer signalStop()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	var cfg Config
	err = envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	c := dependencytrack.NewClient(cfg.DependencytrackUrl, cfg.Username, cfg.Password)

	wg, ctx := errgroup.WithContext(ctx)
	wg.Go(func() error {
		return collectMetrics(ctx, collectMetricsInterval, c)
	})

	wg.Go(func() error {
		return runHttpServer(ctx, cfg.ListenAddress, log.WithField("component", "httpserver"))
	})

	<-ctx.Done()
	signalStop()
	log.Infof("shutting down...")

	ch := make(chan error)
	go func() {
		ch <- wg.Wait()
	}()

	select {
	case <-time.After(10 * time.Second):
		log.Warn("timed out waiting for graceful shutdown")
	case err := <-ch:
		log.Fatalf("error during shutdown: %s", err)
	}
	log.Infof("shutdown complete")
}

func collectMetrics(ctx context.Context, d time.Duration, c *dependencytrack.Client) error {
	ticker := time.NewTicker(time.Second) // initial run
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			ticker.Reset(d) // regular schedule
			start := time.Now()
			log.Infof("start scheduled collectMetrics run")
			err := c.UpdateTotalProjects(ctx)
			if err != nil {
				log.Errorf("UpdateTotalProjects: %s", err)
			}
			log.
				WithFields(log.Fields{
					"duration": time.Since(start),
				}).
				Infof("scheduled collectMetrics run finished")
		}
	}
}
