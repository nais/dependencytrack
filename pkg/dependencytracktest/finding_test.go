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
	// GITHUB finding with a CVE alias: Cve.Id is promoted to the CVE canonical
	assert.Equal(t, "CVE-2024-12797", v.Cve.Id)
	assert.Equal(t, dependencytrack.SeverityLow, v.Cve.Severity)
	assert.Equal(t, "17170e88-cfcb-4900-b3fb-5b0be0a071a5", v.Metadata.ProjectId)
	assert.Equal(t, "5b009251-5efd-4703-8579-49af6cd3d0c6", v.Metadata.ComponentId)
	assert.Equal(t, "6fa86367-6014-427e-8300-69269c16025b", v.Metadata.VulnerabilityUuid)
	assert.Equal(t, fmt.Sprintf("https://nvd.nist.gov/vuln/detail/%s", "CVE-2024-12797"), v.Cve.Link)
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
	makeFinding := func(source, vulnId string, aliases []any) client.Finding {
		return client.Finding{
			Component:     map[string]any{"uuid": "c1", "project": "p1", "purl": "pkg:pypi/test@1.0"},
			Vulnerability: map[string]any{"vulnId": vulnId, "severity": "HIGH", "source": source, "uuid": "u1", "aliases": aliases},
			Analysis:      map[string]any{"isSuppressed": false},
		}
	}

	tests := []struct {
		name     string
		source   string
		vulnId   string
		aliases  []any
		wantId   string
		wantLink string
		wantRefs map[string]string
	}{
		{
			// References trimmed to promoted canonical only — remaining CVE keys
			// would have no cve row and would violate cve_alias_canonical_fkey.
			name:   "GITHUB finding with multiple CVE aliases — lexicographically first CVE wins, refs trimmed to promoted only",
			source: "GITHUB",
			vulnId: "GHSA-79v4-65xg-pq4g",
			aliases: []any{
				map[string]any{"cveId": "CVE-2024-99999", "ghsaId": "GHSA-79v4-65xg-pq4g"},
				map[string]any{"cveId": "CVE-2024-00001", "ghsaId": "GHSA-79v4-65xg-pq4g"},
			},
			wantId:   "CVE-2024-00001",
			wantLink: "https://nvd.nist.gov/vuln/detail/CVE-2024-00001",
			wantRefs: map[string]string{
				"CVE-2024-00001": "GHSA-79v4-65xg-pq4g",
			},
		},
		{
			// GITHUB finding where vulnId is already a CVE — promotion guard
			// (strings.HasPrefix(vulnId, "GHSA-")) prevents overwriting Cve.Id.
			name:   "GITHUB finding with CVE vulnId and ghsaId alias — no promotion, no ref trim",
			source: "GITHUB",
			vulnId: "CVE-2024-12797",
			aliases: []any{
				map[string]any{"cveId": "CVE-2024-12797", "ghsaId": "GHSA-79v4-65xg-pq4g"},
			},
			wantId:   "CVE-2024-12797",
			wantLink: "https://github.com/advisories/CVE-2024-12797",
			wantRefs: map[string]string{"CVE-2024-12797": "GHSA-79v4-65xg-pq4g"},
		},
		{
			name:   "GITHUB finding with cveId and ghsaId — promoted to CVE canonical",
			source: "GITHUB",
			vulnId: "GHSA-79v4-65xg-pq4g",
			aliases: []any{
				map[string]any{"cveId": "CVE-2024-12797", "ghsaId": "GHSA-79v4-65xg-pq4g"},
			},
			wantId:   "CVE-2024-12797",
			wantLink: "https://nvd.nist.gov/vuln/detail/CVE-2024-12797",
			wantRefs: map[string]string{"CVE-2024-12797": "GHSA-79v4-65xg-pq4g"},
		},
		{
			// ghsaId in the alias entry differs from vulnId — the trimmed
			// references map must preserve the original ghsaId value, not vulnId.
			name:   "GITHUB finding where ghsaId differs from vulnId — original ghsaId preserved in refs",
			source: "GITHUB",
			vulnId: "GHSA-79v4-65xg-pq4g",
			aliases: []any{
				map[string]any{"cveId": "CVE-2024-12797", "ghsaId": "GHSA-different-alias"},
			},
			wantId:   "CVE-2024-12797",
			wantLink: "https://nvd.nist.gov/vuln/detail/CVE-2024-12797",
			wantRefs: map[string]string{"CVE-2024-12797": "GHSA-different-alias"},
		},
		{
			name:   "GITHUB finding with cveId only (no ghsaId) — falls back to vulnId as alias, promoted to CVE canonical",
			source: "GITHUB",
			vulnId: "GHSA-79v4-65xg-pq4g",
			aliases: []any{
				map[string]any{"cveId": "CVE-2024-12797"},
			},
			wantId:   "CVE-2024-12797",
			wantLink: "https://nvd.nist.gov/vuln/detail/CVE-2024-12797",
			wantRefs: map[string]string{"CVE-2024-12797": "GHSA-79v4-65xg-pq4g"},
		},
		{
			name:     "GITHUB finding with no aliases — Cve.Id stays as GHSA",
			source:   "GITHUB",
			vulnId:   "GHSA-79v4-65xg-pq4g",
			aliases:  []any{},
			wantId:   "GHSA-79v4-65xg-pq4g",
			wantLink: "https://github.com/advisories/GHSA-79v4-65xg-pq4g",
			wantRefs: map[string]string{},
		},
		{
			name:   "NVD finding with GHSA alias — Cve.Id stays as CVE (no promotion needed)",
			source: "NVD",
			vulnId: "CVE-2024-12797",
			aliases: []any{
				map[string]any{"cveId": "CVE-2024-12797", "ghsaId": "GHSA-79v4-65xg-pq4g"},
			},
			wantId:   "CVE-2024-12797",
			wantLink: "https://nvd.nist.gov/vuln/detail/CVE-2024-12797",
			wantRefs: map[string]string{"CVE-2024-12797": "GHSA-79v4-65xg-pq4g"},
		},
		{
			name:   "cveId present but ghsaId absent, non-GITHUB source — skipped",
			source: "NVD",
			vulnId: "CVE-2024-12797",
			aliases: []any{
				map[string]any{"cveId": "CVE-2024-99999"},
			},
			wantId:   "CVE-2024-12797",
			wantLink: "https://nvd.nist.gov/vuln/detail/CVE-2024-12797",
			wantRefs: map[string]string{},
		},
		{
			name:   "cveId present but ghsaId absent, GITHUB source with CVE vulnId — skipped",
			source: "GITHUB",
			vulnId: "CVE-2024-12797",
			aliases: []any{
				map[string]any{"cveId": "CVE-2024-99999"},
			},
			wantId:   "CVE-2024-12797",
			wantLink: "https://github.com/advisories/CVE-2024-12797",
			wantRefs: map[string]string{},
		},
		{
			name:   "ghsaId present but cveId absent — no alias entry, Cve.Id stays as GHSA",
			source: "GITHUB",
			vulnId: "GHSA-79v4-65xg-pq4g",
			aliases: []any{
				map[string]any{"ghsaId": "GHSA-79v4-65xg-pq4g"},
			},
			wantId:   "GHSA-79v4-65xg-pq4g",
			wantLink: "https://github.com/advisories/GHSA-79v4-65xg-pq4g",
			wantRefs: map[string]string{},
		},
		{
			name:   "cveId equals ghsaId — self-reference skipped, Cve.Id stays as GHSA",
			source: "GITHUB",
			vulnId: "GHSA-79v4-65xg-pq4g",
			aliases: []any{
				map[string]any{"cveId": "GHSA-79v4-65xg-pq4g"},
			},
			wantId:   "GHSA-79v4-65xg-pq4g",
			wantLink: "https://github.com/advisories/GHSA-79v4-65xg-pq4g",
			wantRefs: map[string]string{},
		},
		{
			name:   "empty cveId — skipped",
			source: "GITHUB",
			vulnId: "GHSA-79v4-65xg-pq4g",
			aliases: []any{
				map[string]any{"cveId": "", "ghsaId": "GHSA-79v4-65xg-pq4g"},
			},
			wantId:   "GHSA-79v4-65xg-pq4g",
			wantLink: "https://github.com/advisories/GHSA-79v4-65xg-pq4g",
			wantRefs: map[string]string{},
		},
		{
			name:     "no aliases — empty references, Cve.Id stays as GHSA",
			source:   "GITHUB",
			vulnId:   "GHSA-79v4-65xg-pq4g",
			aliases:  []any{},
			wantId:   "GHSA-79v4-65xg-pq4g",
			wantLink: "https://github.com/advisories/GHSA-79v4-65xg-pq4g",
			wantRefs: map[string]string{},
		},
		{
			name:   "empty vulnId with no ghsaId — skipped",
			source: "GITHUB",
			vulnId: "",
			aliases: []any{
				map[string]any{"cveId": "CVE-2024-12797"},
			},
			wantId:   "",
			wantLink: "https://github.com/advisories/",
			wantRefs: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := dependencytrack.ParseFinding(makeFinding(tt.source, tt.vulnId, tt.aliases))
			assert.NoError(t, err)
			assert.Equal(t, tt.wantId, v.Cve.Id)
			assert.Equal(t, tt.wantLink, v.Cve.Link)
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
