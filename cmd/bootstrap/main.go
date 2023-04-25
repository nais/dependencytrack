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
	"github.com/nais/dependencytrack/pkg/dependencytrack"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
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
	flag.StringVar(&cfg.BaseUrl, "base-url", "http://localhost:9001/api/v1", "base url of dependencytrack")
	flag.StringVar(&cfg.DefaultAdminPassword, "default-admin-password", "admin", "default admin password")
	flag.StringVar(&cfg.AdminPassword, "admin-password", cfg.AdminPassword, "new admin password")
	flag.StringVar(&cfg.UsersFile, "users-file", "/bootstrap/users.yaml", "file with users to create")
}

// TODO: add timer and retry logic to wait for dependencytrack to be ready
func main() {
	parseFlags()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	log, err := setupLogger()
	if err != nil {
		log.Fatalf("setup logger: %v", err)
	}

	c := dependencytrack.NewClient(cfg.BaseUrl, "admin", cfg.AdminPassword, nil, log.WithField("system", "dependencytrack-bootstrap"))

	//if auth doesn't fail, it means we have already changed the password
	_, err = c.Headers(ctx)
	if err != nil {
		err := c.ChangeAdminPassword(ctx, cfg.DefaultAdminPassword, cfg.AdminPassword)
		if err != nil {
			log.Fatalf("change admin password: %v", err)
		}
	}

	file, err := os.ReadFile(cfg.UsersFile)
	if err != nil {
		log.Fatalf("read users file: %v", err)
	}
	users := &dependencytrack.AdminUsers{}
	err = yaml.Unmarshal(file, users)
	if err != nil {
		log.Fatalf("unmarshal users file: %v", err)
	}

	teamUuid, err := c.GetTeamUuid(ctx, "Administrators")
	if err != nil {
		log.Fatalf("get team uuid: %v", err)
	}

	//remove users before adding to ensure passwords in sync
	err = c.RemoveAdminUsers(ctx, users)
	if err != nil {
		log.Fatalf("remove users: %v", err)
	}

	err = c.CreateAdminUsers(ctx, users, teamUuid)
	if err != nil {
		log.Fatalf("create users: %v", err)
	}

	log.Infof("done: created users and added to Administrators team")

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
