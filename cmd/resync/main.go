package main

import (
	"context"
	"flag"
	"github.com/nais/dependencytrack/cmd/common"
	"github.com/nais/dependencytrack/pkg/client"
	"os/signal"
	"syscall"
)

var cfg = &Config{
	LogLevel: "debug",
}

type Config struct {
	Password string `json:"password"`
	BaseUrl  string `json:"base-url"`
	LogLevel string `json:"log-level"`
}

func init() {
	flag.StringVar(&cfg.LogLevel, "log-level", cfg.LogLevel, "which log level to use, default 'info'")
	flag.StringVar(&cfg.BaseUrl, "base-url", "http://localhost:9001", "base url of dependencytrack")
	flag.StringVar(&cfg.Password, "password", cfg.Password, "password for application")
}

func main() {
	common.ParseFlags()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	log, err := common.SetupLogger(cfg.LogLevel)
	if err != nil {
		log.Fatalf("setup logger: %v", err)
	}

	c := client.New(
		cfg.BaseUrl,
		"admin",
		cfg.Password,
		client.WithLogger(log.WithField("system", "dtrack-resync")),
	)

	projects, err := c.GetProjects(ctx)
	total := len(projects)
	log.Infof("number of projects to re-analyze: %v", total)
	if err != nil {
		log.Fatalf("get projects: %v", err)
	}

	for _, project := range projects {
		if err = c.TriggerAnalysis(ctx, project.Uuid); err != nil {
			log.Errorf("trigger analysis: %v", err)
		}
		total--
	}

	log.Info("sync complete, projects left to analyze: ", total)
}
