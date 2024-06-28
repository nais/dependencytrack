package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/nais/dependencytrack/internal/dependencytrack"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

const collectMetricsInterval = 3 * time.Minute

type Config struct {
	Port               string `envconfig:"PORT" default:"8000"`
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

	// TODO: how to do graceful shutdown for gin
	router := gin.Default()
	router.GET("/ping", pong())
	router.GET("/metrics", metrics())

	srv := &http.Server{
		Addr:    "localhost:" + cfg.Port,
		Handler: router.Handler(),
	}

	wg.Go(func() error {
		log.Infof("HTTP server accepting requests on %q", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.WithError(err).Infof("unexpected error from HTTP server")
			return err
		}
		log.Infof("HTTP server finished, terminating...")
		return nil
	})

	wg.Go(func() error {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		log.Infof("HTTP server shutting down...")
		if err := srv.Shutdown(ctx); err != nil && !errors.Is(err, context.Canceled) {
			log.WithError(err).Infof("HTTP server shutdown failed")
			return err
		}
		return nil
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

func pong() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	}
}

func metrics() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
