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
	LogLevel             string `json:"log-level"`
	BaseUrl              string `json:"base-url"`
	DefaultAdminPassword string `json:"default-admin-password"`
	AdminPassword        string `json:"admin-password"`
	NewUser              string `json:"new-user"`
	NewUserPassword      string `json:"new-user-password"`
}

func init() {
	flag.StringVar(&cfg.LogLevel, "log-level", cfg.LogLevel, "which log level to use, default 'info'")
	flag.StringVar(&cfg.BaseUrl, "base-url", "http://localhost:9001", "base url of dependencytrack")
	flag.StringVar(&cfg.DefaultAdminPassword, "default-admin-password", "admin", "default admin password")
	flag.StringVar(&cfg.AdminPassword, "admin-password", cfg.AdminPassword, "new admin password")
	flag.StringVar(&cfg.NewUser, "new-user", cfg.NewUser, "new admin user")
	flag.StringVar(&cfg.NewUserPassword, "new-user-password", cfg.NewUserPassword, "new admin user password")
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

	err = c.NewManagedUser(ctx, cfg.NewUser, cfg.NewUserPassword)
	if err != nil {
		e, ok := err.(*client.RequestError)
		if ok {
			if e.StatusCode == 409 {
				log.Infof("user %s already exists", cfg.NewUser)
				return
			}
		}
		log.Fatalf("creating new managed user: %v", err)
	}

	err = c.AddToTeam(ctx, cfg.NewUser, "Administrators")
	if err != nil {
		log.Fatalf("adding user to team: %v", err)
	}
	log.Infof("done: created %s and added to Administrators team", "user1")

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
