package test

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	log "github.com/sirupsen/logrus"
)

func DependencyTrackPool() (string, func()) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	opts := &dockertest.RunOptions{
		Repository:   "dependencytrack/apiserver",
		Tag:          "4.8.0",
		ExposedPorts: []string{"8080"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"8080": {
				{HostIP: "0.0.0.0", HostPort: "8080"},
			},
		},
	}

	resource, err := pool.RunWithOptions(opts)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	err = resource.Expire(60)
	if err != nil {
		log.Fatalf("Could not set expiration: %s", err)
	}
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://localhost:8080", nil)
		if err != nil {
			return fmt.Errorf("creating request: %w", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode > 299 {
			return fmt.Errorf("got status code %d", resp.StatusCode)
		}
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to dtrack: %s", err)
	}

	return "http://localhost:8080", func() {
		// You can't defer this because os.Exit doesn't care for defer
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}
}
