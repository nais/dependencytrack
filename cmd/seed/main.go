package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/in-toto/in-toto-golang/in_toto"
	"github.com/nais/dependencytrack/pkg/dependencytrack"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	seedDependencyTrack(context.Background())
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

	att := &in_toto.CycloneDXStatement{}

	var predicate any
	err = json.Unmarshal(b, &predicate)
	if err != nil {
		return err
	}

	att.Predicate = predicate

	imgPrefix := "ghcr.io/nais/nais-deploy-chicken-%d"
	for i := 1; i < 9; i++ {
		_, err = c.CreateProjectWithSbom(ctx, att, fmt.Sprintf(imgPrefix, i), "1")
		if err != nil {
			return err
		}
	}
	return nil
}
