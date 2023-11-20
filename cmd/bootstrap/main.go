package main

import (
	"context"
	"flag"
	"github.com/nais/dependencytrack/cmd/common"
	"github.com/nais/dependencytrack/pkg/client"
	"gopkg.in/yaml.v3"
	"os"
	"os/signal"
	"strings"
	"syscall"
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
	GithubAdvisoryToken  string `json:"github-advisory-token"`
	GoogleOSVEnabled     bool   `json:"google-osv-enabled"`
}

func init() {
	flag.StringVar(&cfg.LogLevel, "log-level", cfg.LogLevel, "which log level to use, default 'info'")
	flag.StringVar(&cfg.BaseUrl, "base-url", "http://localhost:9001", "base url of dependencytrack")
	flag.StringVar(&cfg.DefaultAdminPassword, "default-admin-password", "admin", "default admin password")
	flag.StringVar(&cfg.AdminPassword, "admin-password", cfg.AdminPassword, "new admin password")
	flag.StringVar(&cfg.GithubAdvisoryToken, "github-advisory-token", cfg.GithubAdvisoryToken, "github advisory mirroring token")
	flag.StringVar(&cfg.UsersFile, "users-file", "/bootstrap/users.yaml", "file with users to create")
	flag.BoolVar(&cfg.GoogleOSVEnabled, "google-osv-enabled", cfg.GoogleOSVEnabled, "enable google osv integration")
}

// TODO: add timer and retry logic to wait for dependencytrack to be ready
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
		cfg.AdminPassword,
		client.WithLogger(log.WithField("system", "dependencytrack-bootstrap")),
	)

	// if auth doesn't fail, it means we have already changed the password
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
	users := &client.AdminUsers{}
	err = yaml.Unmarshal(file, users)
	if err != nil {
		log.Fatalf("unmarshal users file: %v", err)
	}

	team, err := c.GetTeam(ctx, "Administrators")
	if err != nil {
		log.Fatalf("get team uuid: %v", err)
	}

	if len(team.ApiKeys) == 0 {
		_, err = c.GenerateApiKey(ctx, team.Uuid)
		if err != nil {
			log.Fatalf("generate api key: %v", err)
		}
	}

	// remove users before adding to ensure passwords in sync
	err = c.RemoveAdminUsers(ctx, users)
	if err != nil {
		log.Fatalf("remove users: %v", err)
	}

	err = c.CreateAdminUsers(ctx, users, team.Uuid)
	if err != nil {
		log.Fatalf("create users: %v", err)
	}

	log.Infof("done: created users and added to Administrators team")

	props, err := c.GetConfigProperties(ctx)
	if err != nil {
		log.Fatalf("get config properties: %v", err)
	}

	var cp []client.ConfigProperty
	for _, prop := range props {
		if cfg.GithubAdvisoryToken != "" {
			switch prop.PropertyName {
			case "github.advisories.enabled":
				prop.PropertyValue = "true"
				cp = append(cp, prop)
				log.Infof("done: added github advisory mirroring")
			case "github.advisories.access.token":
				prop.PropertyValue = cfg.GithubAdvisoryToken
				cp = append(cp, prop)
				log.Infof("done: added github advisory mirroring token")
			}
		}

		if cfg.GoogleOSVEnabled {
			switch prop.PropertyName {
			case "google.osv.enabled":
				eco, err := c.GetEcosystems(ctx)
				if err != nil {
					log.Fatalf("get ecosystems: %v", err)
				}
				// if the list is empty we activated all ecosystems
				if len(eco) > 0 {
					prop.PropertyValue = strings.Join(eco, ";")
					cp = append(cp, prop)
					log.Infof("done: added github osv integration")
				}
			}
		}
	}

	// only update if we have new properties
	if cp != nil {
		if _, err := c.ConfigPropertyAggregate(ctx, cp); err != nil {
			log.Fatalf("config property aggregate: %v", err)
		}
		log.Infof("done: added config properties")
	}
}
