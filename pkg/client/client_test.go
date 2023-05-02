package client

import (
	"context"
	"github.com/in-toto/in-toto-golang/in_toto"
	"testing"
)

func TestClient_UploadProject(t *testing.T) {
	ctx := context.Background()
	c := New("http://localhost:9001/api/v1", "admin", "yolo", WithApiKeySource("Administrators"))

	bom := in_toto.CycloneDXStatement{}

	err := c.UploadProject(ctx, "test", "1.0.0", &bom)
	if err != nil {
		t.Fatal(err)
	}

}
