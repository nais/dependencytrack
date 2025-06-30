package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/nais/dependencytrack/cmd/common"
	"github.com/nais/dependencytrack/pkg/dependencytrack"
	"github.com/nais/dependencytrack/pkg/dependencytrack/client"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var cfg = &Config{
	LogLevel: "debug",
}

type Config struct {
	AdminPassword        string `json:"admin-password"`
	FrontendBaseUrl      string `json:"frontend-base-url"`
	BaseUrl              string `json:"base-url"`
	DefaultAdminPassword string `json:"default-admin-password"`
	LogLevel             string `json:"log-level"`
	UsersFile            string `json:"users-file"`
	GithubAdvisoryToken  string `json:"github-advisory-token"`
	GoogleOSVEnabled     bool   `json:"google-osv-enabled"`
	NVDApiKey            string `json:"nvd-api-key"`
	TrivyApiToken        string `json:"trivy-token"`
	TrivyBaseURL         string `json:"trivy-base-url"`
	TrivyIgnoreUnfixed   bool   `json:"trivy-ignore-unfixed"`
}

func init() {
	flag.StringVar(&cfg.LogLevel, "log-level", cfg.LogLevel, "which log level to use, default 'info'")
	flag.StringVar(&cfg.FrontendBaseUrl, "frontend-base-url", "http://localhost:9000", "frontend base url")
	flag.StringVar(&cfg.BaseUrl, "base-url", "http://localhost:9001", "base url of dependencytrack")
	flag.StringVar(&cfg.DefaultAdminPassword, "default-admin-password", "admin", "default admin password")
	flag.StringVar(&cfg.AdminPassword, "admin-password", cfg.AdminPassword, "new admin password")
	flag.StringVar(&cfg.GithubAdvisoryToken, "github-advisory-token", cfg.GithubAdvisoryToken, "github advisory mirroring token")
	flag.StringVar(&cfg.UsersFile, "users-file", "/bootstrap/users.yaml", "file with users to create")
	flag.BoolVar(&cfg.GoogleOSVEnabled, "google-osv-enabled", cfg.GoogleOSVEnabled, "enable google osv integration")
	flag.StringVar(&cfg.NVDApiKey, "nvd-api-key", cfg.NVDApiKey, "nvd api key")
	flag.StringVar(&cfg.TrivyApiToken, "trivy-api-token", cfg.TrivyApiToken, "trivy api token to use for scanning")
	flag.StringVar(&cfg.TrivyBaseURL, "trivy-base-url", cfg.TrivyBaseURL, "trivy base url")
	flag.BoolVar(&cfg.TrivyIgnoreUnfixed, "trivy-ignore-unfixed", cfg.TrivyIgnoreUnfixed, "ignore unfixed vulnerabilities")
}

func main() {
	common.ParseFlags()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	log, err := common.SetupLogger(cfg.LogLevel)
	if err != nil {
		log.Fatalf("setup logger: %v", err)
	}

	c, err := dependencytrack.NewClient(cfg.BaseUrl, "admin", cfg.AdminPassword, log.WithField("subsystem", "client"))
	if err != nil {
		log.Fatalf("create dependencytrack client: %v", err)
	}

	err = c.ChangeAdminPassword(ctx, cfg.DefaultAdminPassword, cfg.AdminPassword)
	if err != nil {
		log.Fatalf("change admin password: %v", err)
	}

	ctx, err = c.AuthContext(ctx)
	if err != nil {
		log.Fatalf("login failed: %v", err)
	}

	file, err := os.ReadFile(cfg.UsersFile)
	if err != nil {
		log.Fatalf("read users file: %v", err)
	}
	type Users struct {
		Users []*dependencytrack.AdminUser `yaml:"users"`
	}
	users := &Users{}
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
	err = c.RemoveAdminUsers(ctx, users.Users)
	if err != nil {
		log.Fatalf("remove users: %v", err)
	}

	err = c.CreateAdminUsers(ctx, users.Users, team.Uuid)
	if err != nil {
		log.Fatalf("create users: %v", err)
	}

	log.Info("created: users and added to Administrators team")

	props, err := c.GetConfigProperties(ctx)
	if err != nil {
		log.Fatalf("get config properties: %v", err)
	}

	var cp []client.ConfigProperty
	for _, prop := range props {
		if cfg.GithubAdvisoryToken != "" {
			switch *prop.PropertyName {
			case "github.advisories.enabled":
				if isAlreadySet(prop.PropertyValue, "true") {
					log.Info("github advisory mirroring already enabled")
					continue
				}
				token := "true"
				prop.PropertyValue = &token
				cp = append(cp, prop)
				log.Info("added: github advisory mirroring")
			case "github.advisories.access.token":
				if isAlreadySet(prop.PropertyValue, cfg.GithubAdvisoryToken) {
					log.Info("github advisory mirroring token already set")
					continue
				}
				prop.PropertyValue = &cfg.GithubAdvisoryToken
				cp = append(cp, prop)
				log.Info("added: github advisory mirroring token")
			}
		}

		if cfg.NVDApiKey != "" {
			switch *prop.PropertyName {
			case "nvd.api.enabled":
				if isAlreadySet(prop.PropertyValue, "true") {
					log.Info("nvd api already enabled")
					continue
				}
				enabled := "true"
				prop.PropertyValue = &enabled
				cp = append(cp, prop)
				log.Info("added: nvd api")
			case "nvd.api.download.feeds":
				if isAlreadySet(prop.PropertyValue, "false") {
					log.Info("nvd api download feeds already disabled")
					continue
				}
				download := "false"
				prop.PropertyValue = &download
				cp = append(cp, prop)
				log.Info("added: nvd api download feeds")
			case "nvd.api.key":
				if isAlreadySet(prop.PropertyValue, cfg.NVDApiKey) {
					log.Info("nvd api key already set")
					continue
				}
				prop.PropertyValue = &cfg.NVDApiKey
				cp = append(cp, prop)
				log.Info("added: nvd api key")
			}
		}

		if cfg.GoogleOSVEnabled {
			switch *prop.PropertyName {
			case "google.osv.enabled":
				eco, err := c.GetEcosystems(ctx)
				if err != nil {
					log.Fatalf("get ecosystems: %v", err)
				}
				// if the list is empty we activated all ecosystems
				if len(eco) == 0 {
					log.Info("google osv integration already enabled")
					continue
				}

				if err = updateEcosystems(ctx, c, eco, prop, log); err != nil {
					log.Fatalf("update ecosystems: %v", err)
				}
			}
		}

		if cfg.TrivyApiToken != "" {
			switch *prop.PropertyName {
			case "trivy.enabled":
				if isAlreadySet(prop.PropertyValue, "true") {
					log.Info("trivy integration already enabled")
					continue
				}
				enabled := "true"
				prop.PropertyValue = &enabled
				cp = append(cp, prop)
				log.Info("added: trivy integration")
			case "trivy.api.token":
				// we cant check if the token is already set, so we just set it
				prop.PropertyValue = &cfg.TrivyApiToken
				cp = append(cp, prop)
				log.Info("added: trivy token")
			case "trivy.base.url":
				if isAlreadySet(prop.PropertyValue, cfg.TrivyBaseURL) {
					log.Info("trivy base url already set")
					continue
				}
				prop.PropertyValue = &cfg.TrivyBaseURL
				cp = append(cp, prop)
				log.Info("added: trivy base url")
			case "trivy.ignore.unfixed":
				if isAlreadySet(prop.PropertyValue, "true") {
					log.Info("trivy ignore unfixed already enabled")
					continue
				}
				unfixed := "true"
				prop.PropertyValue = &unfixed
				cp = append(cp, prop)
				log.Info("added: trivy ignore unfixed")
			}
		}

		if cfg.FrontendBaseUrl != "" {
			switch *prop.PropertyName {
			case "base.url":
				if isAlreadySet(prop.PropertyValue, cfg.FrontendBaseUrl) {
					log.Info("general base url already set")
					continue
				}
				prop.PropertyValue = &cfg.FrontendBaseUrl
				cp = append(cp, prop)
				log.Info("added: general base url")
			}
		}
	}

	// only update if we have new properties
	if len(cp) > 0 {
		for _, p := range cp {
			if err = c.ConfigPropertyAggregate(ctx, p); err != nil {
				log.Fatalf("config property aggregate: %v", err)
			}
		}
		log.Info("done: config properties updated")
	}
}

func updateEcosystems(ctx context.Context, c dependencytrack.Client, eco []string, prop client.ConfigProperty, log *logrus.Logger) error {
	chunkSize := 10

	for len(eco) > 0 {
		log.Info("Processing chunk of ecosystems: ", len(eco))
		// Get the next chunk of eco
		end := chunkSize
		if len(eco) < chunkSize {
			end = len(eco)
		}

		chunk := eco[:end]
		eco = eco[end:] // Remove the processed chunk from eco

		propVal := strings.Join(chunk, ";")
		p := client.ConfigProperty{
			GroupName:     prop.GroupName,
			PropertyName:  prop.PropertyName,
			PropertyType:  prop.PropertyType,
			PropertyValue: &propVal,
			Description:   prop.Description,
		}

		if err := c.ConfigPropertyAggregate(ctx, p); err != nil {
			return err
		}

		log.Info("Chunk processed and sent. Remaining items:", len(eco))
	}
	return nil
}

func isAlreadySet(config *string, inputValue string) bool {
	if config == nil {
		return false
	}
	return strings.EqualFold(*config, inputValue)
}
