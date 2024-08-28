package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/nais/dependencytrack/pkg/client"

	"github.com/sirupsen/logrus"
)

var cfg = &Config{
	LogLevel: "debug",
}

type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
	BaseUrl  string `json:"base-url"`
	LogLevel string `json:"log-level"`
}

func init() {
	flag.StringVar(&cfg.LogLevel, "log-level", cfg.LogLevel, "which log level to use, default 'debug'")
	flag.StringVar(&cfg.BaseUrl, "base-url", "http://localhost:9001", "base url of dependencytrack")
	flag.StringVar(&cfg.Password, "password", cfg.Password, "password for application")
	flag.StringVar(&cfg.Username, "username", "dtrack-resync", "username for application")
}

func main() {
	parseFlags()
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	log, err := setupLogger()
	if err != nil {
		log.Fatalf("setup logger: %v", err)
	}

	c := client.New(
		cfg.BaseUrl,
		cfg.Username,
		cfg.Password,
		client.WithLogger(log.WithField("system", "dtrack-resync")),
	)

	projects, err := c.GetProjects(ctx)
	log.Infof("number of projects to re-analyze: %v", len(projects))
	if err != nil {
		log.Fatalf("get projects: %v", err)
	}

	var projectlog string
	now := time.Now()
	for _, project := range projects {
		if err = c.TriggerAnalysis(ctx, project.Uuid); err != nil {
			log.Errorf("trigger analysis: %v", err)
		}
		projectlog += "Updated project: " + project.Name + "\n"
	}

	log.Infof(projectlog)
	log.Infof("Syncing time: %v\n", time.Since(now))
	log.Info("sync complete")
}

func parseFlags() {
	err := godotenv.Load()
	if err != nil {
		logrus.Debugf("loading .env file %v", err)
	}

	flag.VisitAll(func(f *flag.Flag) {
		name := strings.ToUpper(strings.Replace(f.Name, "-", "_", -1))
		if value, ok := os.LookupEnv(name); ok {
			err = flag.Set(f.Name, value)
			if err != nil {
				logrus.Fatalf("failed setting flag from environment: %v", err)
				return
			}
		}
	})

	flag.Parse()
}

func setupLogger() (*logrus.Logger, error) {
	log := logrus.StandardLogger()
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})
	lvl, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		return nil, err
	}
	log.SetLevel(lvl)
	return log, err
}
