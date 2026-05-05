package dependencytracktest

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/nais/dependencytrack/pkg/dependencytrack"

	"github.com/nais/dependencytrack/pkg/dependencytrack/client"
	"github.com/stretchr/testify/assert"
)

func TestParseFinding(t *testing.T) {
	b, err := os.ReadFile("testdata/finding.json")
	assert.NoError(t, err)
	var f client.Finding
	err = json.Unmarshal(b, &f)
	assert.NoError(t, err)

	v, err := dependencytrack.ParseFinding(f)
	assert.NoError(t, err)
	assert.Equal(t, "pkg:pypi/cryptography@43.0.1", v.Package)
	assert.Equal(t, "GHSA-79v4-65xg-pq4g", v.Cve.Id)
	assert.Equal(t, dependencytrack.SeverityLow, v.Cve.Severity)
	assert.Equal(t, "17170e88-cfcb-4900-b3fb-5b0be0a071a5", v.Metadata.ProjectId)
	assert.Equal(t, "5b009251-5efd-4703-8579-49af6cd3d0c6", v.Metadata.ComponentId)
	assert.Equal(t, "6fa86367-6014-427e-8300-69269c16025b", v.Metadata.VulnerabilityUuid)
	assert.Equal(t, fmt.Sprintf("https://github.com/advisories/%s", "GHSA-79v4-65xg-pq4g"), v.Cve.Link)
	assert.Equal(t, true, v.Suppressed)
	assert.Equal(t, "Vulnerable OpenSSL included in cryptography wheels", v.Cve.Title)
	assert.Equal(t, "a loooong description", v.Cve.Description)
	assert.Equal(t, "44.0.1", v.LatestVersion)
	assert.Equal(t, 1, len(v.Cve.References))
	assert.NotNil(t, v.Cvss)
	assert.Equal(t, 2.5, *v.Cvss)
	assert.NotNil(t, v.EpssScore)
	assert.Equal(t, 0.00527, *v.EpssScore)
	assert.NotNil(t, v.EpssPercentile)
	assert.Equal(t, 0.66622, *v.EpssPercentile)
}

func TestParseFinding_Aliases(t *testing.T) {
	makeFinding := func(vulnId string, aliases []any) client.Finding {
		return client.Finding{
			Component:     map[string]any{"uuid": "c1", "project": "p1", "purl": "pkg:pypi/test@1.0"},
			Vulnerability: map[string]any{"vulnId": vulnId, "severity": "HIGH", "source": "GITHUB", "uuid": "u1", "aliases": aliases},
			Analysis:      map[string]any{"isSuppressed": false},
		}
	}

	tests := []struct {
		name     string
		vulnId   string
		aliases  []any
		wantRefs map[string]string
	}{
		{
			name:   "both cveId and ghsaId present",
			vulnId: "GHSA-79v4-65xg-pq4g",
			aliases: []any{
				map[string]any{"cveId": "CVE-2024-12797", "ghsaId": "GHSA-79v4-65xg-pq4g"},
			},
			wantRefs: map[string]string{"CVE-2024-12797": "GHSA-79v4-65xg-pq4g"},
		},
		{
			name:   "cveId present but ghsaId absent — falls back to vulnId",
			vulnId: "GHSA-79v4-65xg-pq4g",
			aliases: []any{
				map[string]any{"cveId": "CVE-2024-12797"},
			},
			wantRefs: map[string]string{"CVE-2024-12797": "GHSA-79v4-65xg-pq4g"},
		},
		{
			name:   "ghsaId present but cveId absent — no entry emitted",
			vulnId: "GHSA-79v4-65xg-pq4g",
			aliases: []any{
				map[string]any{"ghsaId": "GHSA-79v4-65xg-pq4g"},
			},
			wantRefs: map[string]string{},
		},
		{
			name:   "cveId equals vulnId — self-reference skipped",
			vulnId: "CVE-2024-12797",
			aliases: []any{
				map[string]any{"cveId": "CVE-2024-12797", "ghsaId": "GHSA-79v4-65xg-pq4g"},
			},
			wantRefs: map[string]string{},
		},
		{
			name:   "empty cveId — skipped",
			vulnId: "GHSA-79v4-65xg-pq4g",
			aliases: []any{
				map[string]any{"cveId": "", "ghsaId": "GHSA-79v4-65xg-pq4g"},
			},
			wantRefs: map[string]string{},
		},
		{
			name:     "no aliases — empty references",
			vulnId:   "GHSA-79v4-65xg-pq4g",
			aliases:  []any{},
			wantRefs: map[string]string{},
		},
		{
			name:   "empty vulnId with no ghsaId — skipped",
			vulnId: "",
			aliases: []any{
				map[string]any{"cveId": "CVE-2024-12797"},
			},
			wantRefs: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := dependencytrack.ParseFinding(makeFinding(tt.vulnId, tt.aliases))
			assert.NoError(t, err)
			assert.Equal(t, tt.wantRefs, v.Cve.References)
		})
	}
}

func TestParseFinding_CvssScore(t *testing.T) {
	ptrF := func(f float64) *float64 { return &f }

	tests := []struct {
		name     string
		vuln     map[string]any
		wantCvss *float64
	}{
		{
			name: "v4 takes precedence over v3",
			vuln: map[string]any{
				"vulnId":          "CVE-2025-10492",
				"severity":        "HIGH",
				"cvssV3BaseScore": float64(9.8),
				"cvssV4Score":     float64(8.7),
				"source":          "NVD",
				"uuid":            "00000000-0000-0000-0000-000000000001",
			},
			wantCvss: ptrF(8.7),
		},
		{
			name: "v4 only",
			vuln: map[string]any{
				"vulnId":      "GHSA-gg57-587f-h5v6",
				"severity":    "MEDIUM",
				"cvssV4Score": float64(5.1),
				"source":      "GITHUB",
				"uuid":        "00000000-0000-0000-0000-000000000002",
			},
			wantCvss: ptrF(5.1),
		},
		{
			name: "v3 only",
			vuln: map[string]any{
				"vulnId":          "GHSA-3vqj-43w4-2q58",
				"severity":        "HIGH",
				"cvssV3BaseScore": float64(7.5),
				"source":          "GITHUB",
				"uuid":            "00000000-0000-0000-0000-000000000003",
			},
			wantCvss: ptrF(7.5),
		},
		{
			name: "neither present",
			vuln: map[string]any{
				"vulnId":   "GHSA-0000-0000-0000",
				"severity": "UNASSIGNED",
				"source":   "GITHUB",
				"uuid":     "00000000-0000-0000-0000-000000000004",
			},
			wantCvss: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			finding := client.Finding{
				Component:     map[string]any{"uuid": "c1", "project": "p1", "purl": "pkg:maven/test@1.0"},
				Vulnerability: tt.vuln,
				Analysis:      map[string]any{"isSuppressed": false},
			}
			v, err := dependencytrack.ParseFinding(finding)
			assert.NoError(t, err)
			if tt.wantCvss == nil {
				assert.Nil(t, v.Cvss)
			} else {
				assert.NotNil(t, v.Cvss)
				assert.Equal(t, *tt.wantCvss, *v.Cvss)
			}
		})
	}
}
