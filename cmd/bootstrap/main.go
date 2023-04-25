package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"bootstrap/internal/client"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var cfg = &Config{
	LogLevel: "debug",
}

type Config struct {
	AdminPassword        string `json:"admin-password"`
	BaseUrl              string `json:"base-url"`
	DefaultAdminPassword string `json:"default-admin-password"`
	LogLevel             string `json:"log-level"`
	UsersFile            string `json:"users-file"`
}

func init() {
	flag.StringVar(&cfg.LogLevel, "log-level", cfg.LogLevel, "which log level to use, default 'info'")
	flag.StringVar(&cfg.BaseUrl, "base-url", "http://localhost:9001", "base url of dependencytrack")
	flag.StringVar(&cfg.DefaultAdminPassword, "default-admin-password", "admin", "default admin password")
	flag.StringVar(&cfg.AdminPassword, "admin-password", cfg.AdminPassword, "new admin password")
	flag.StringVar(&cfg.UsersFile, "users-file", "/bootstrap/users.yaml", "file with users to create")
}

// TODO: add timer and retry logic to wait for dependencytrack to be ready
func main() {
	parseFlags()
	setupLogger()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	c := client.New(cfg.BaseUrl, "admin", cfg.AdminPassword)

	_, err := c.Token(ctx)
	if err != nil {
		err := c.ChangeAdminPassword(ctx, cfg.DefaultAdminPassword, cfg.AdminPassword)
		if err != nil {
			log.Fatalf("change admin password: %v", err)
		}
	}

	teamUuid, err := c.GetTeamUuid(ctx, "Administrators")
	if err != nil {
		log.Fatalf("get team uuid: %v", err)
	}
	//remove users before adding to ensure passwords in sync
	err = c.RemoveUsers(ctx, cfg.UsersFile)
	if err != nil {
		log.Fatalf("remove users: %v", err)
	}

	err = c.CreateUsers(ctx, cfg.UsersFile, teamUuid)
	if err != nil {
		log.Fatalf("create users: %v", err)
	}

	log.Infof("done: created users and added to Administrators team")

}

func parseFlags() {
	err := godotenv.Load()
	if err != nil {
		log.Debugf("loading .env file %v", err)
	}

	flag.VisitAll(func(f *flag.Flag) {
		name := strings.ToUpper(strings.Replace(f.Name, "-", "_", -1))
		if value, ok := os.LookupEnv(name); ok {
			err = flag.Set(f.Name, value)
			if err != nil {
				log.Fatalf("failed setting flag from environment: %v", err)
				return
			}
		}
	})

	flag.Parse()
}

func setupLogger() {
	log.SetFormatter(&log.JSONFormatter{})
	l, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(l)
}
