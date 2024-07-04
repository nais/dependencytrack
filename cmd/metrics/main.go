package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/nais/dependencytrack/internal/dependencytrack"
	"github.com/nais/dependencytrack/internal/http"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"os/signal"
	"syscall"
	"time"
)

const (
	collectMetricsInterval     = 3 * time.Minute
	dependencyTrackSyncTimeout = time.Minute
)

type Config struct {
	DependencytrackUrl string `envconfig:"DEPENDENCYTRACK_URL" default:"http://localhost:8080"`
	ListenAddress      string `envconfig:"LISTEN_ADDRESS" default:":8000"`
	Password           string `envconfig:"PASSWORD"`
	Tenant             string `envconfig:"TENANT" default:"test"`
	Username           string `envconfig:"USERNAME" default:"admin"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx, signalStop := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer signalStop()

	err := godotenv.Load()
	if err != nil {
		log.Debugf("not loading .env file: %s", err.Error())
	}

	var cfg Config
	err = envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	c := dependencytrack.NewClient(cfg.DependencytrackUrl, cfg.Username, cfg.Password)

	wg, ctx := errgroup.WithContext(ctx)
	wg.Go(func() error {
		return collectMetrics(ctx, collectMetricsInterval, c, cfg.Tenant)
	})

	wg.Go(func() error {
		return http.RunHttpServer(ctx, cfg.ListenAddress, log.WithField("component", "httpserver"))
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

func collectMetrics(ctx context.Context, d time.Duration, c *dependencytrack.Client, tenant string) error {
	ticker := time.NewTicker(time.Second) // initial run
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			ticker.Reset(d) // regular schedule

			ctx, cancel := context.WithTimeout(ctx, dependencyTrackSyncTimeout)
			defer cancel()

			start := time.Now()
			log.Infof("start scheduled collectMetrics run")
			err := c.UpdateTotalProjects(ctx, tenant)
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
