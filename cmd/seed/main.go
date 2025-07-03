package main

import (
	"context"
	"fmt"
	"github.com/nais/dependencytrack/pkg/dependencytrack"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	err := seedDependencyTrack(context.Background())
	if err != nil {
		log.Fatalf("Failed to seed DependencyTrack: %v", err)
	}
}

func seedDependencyTrack(ctx context.Context) error {
	c, err := dependencytrack.NewClient(
		"http://localhost:9010/api",
		"admin",
		"yolo",
		log.WithField("subsystem", "dp-client"),
	)
	if err != nil {
		return err
	}

	b, err := os.ReadFile("local/vuln-nginx.json")
	if err != nil {
		return err
	}

	imgPrefix := "ghcr.io/nais/nais-deploy-chicken-%d"
	for i := 1; i < 9; i++ {
		_, err = c.CreateProjectWithSbom(ctx, fmt.Sprintf(imgPrefix, i), "1", b)
		if err != nil {
			return err
		}
	}
	return nil
}
