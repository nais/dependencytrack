package dependencytrack

import (
	"context"
	"fmt"
	"strings"

	"github.com/nais/dependencytrack/pkg/dependencytrack/client"
)

type Severity string

const SeverityCritical = Severity("CRITICAL")
const SeverityHigh = Severity("HIGH")
const SeverityMedium = Severity("MEDIUM")
const SeverityLow = Severity("LOW")
const SeverityUnassigned = Severity("UNASSIGNED")

type Vulnerability struct {
	Package       string
	Suppressed    bool
	Cve           *Cve
	LatestVersion string
	Metadata      *VulnMetadata
}

type Cve struct {
	Id          string
	Description string
	Title       string
	Link        string
	Severity    Severity
	References  map[string]string
}

type VulnMetadata struct {
	ProjectId         string
	ComponentId       string
	VulnerabilityUuid string
}

// GetVulnerabilities Is this function lacking pagination for all findings in a project or do we not need it?
// https://github.com/DependencyTrack/dependency-track/issues/3811
// https://github.com/DependencyTrack/dependency-track/issues/4677
func (c *dependencyTrackClient) GetVulnerabilities(ctx context.Context, uuid, vulnId string, suppressed bool) ([]*Vulnerability, error) {
	return withAuthContextValue(c, ctx, func(tokenCtx context.Context) ([]*Vulnerability, error) {
		req := c.client.FindingAPI.GetFindingsByProject(tokenCtx, uuid).Suppressed(suppressed)

		switch {
		case strings.Contains(vulnId, "CVE-"):
			req.Source("NVD")
		case strings.Contains(vulnId, "GHSA-"):
			req.Source("GITHUB")
		case strings.Contains(vulnId, "TRIVY-"):
			req.Source("TRIVY")
		case strings.Contains(vulnId, "NPM-"):
			req.Source("NPM")
		case strings.Contains(vulnId, "UBUNTU-"):
			req.Source("UBUNTU")
		case strings.Contains(vulnId, "OSSINDEX-"):
			req.Source("OSSINDEX")
		}

		findings, resp, err := req.Execute()
		if err != nil {
			return nil, convertError(err, "GetVulnerabilities", resp)
		}

		vulns := make([]*Vulnerability, 0)
		for _, f := range findings {
			v, err := parseVulnerability(f)
			if err != nil {
				return nil, err
			}
			vulns = append(vulns, v)
		}

		return vulns, nil
	})
}

// TriggerAnalysis triggers an analysis for a project in Dependency-Track.
func (c *dependencyTrackClient) TriggerAnalysis(ctx context.Context, uuid string) error {
	// Fire and forget
	return c.withAuthContext(ctx, func(tokenCtx context.Context) error {
		_, _, err := c.client.FindingAPI.AnalyzeProject(tokenCtx, uuid).Execute()
		if err != nil {
			return fmt.Errorf("failed to trigger analysis: %w", err)
		}
		return nil
	})
}

func parseVulnerability(finding client.Finding) (*Vulnerability, error) {
	component, componentOk := finding.GetComponentOk()
	if !componentOk {
		return nil, fmt.Errorf("missing component for finding")
	}

	analysis, analysisOk := finding.GetAnalysisOk()
	if !analysisOk {
		return nil, fmt.Errorf("missing analysis for finding")
	}

	vulnData, vulnOk := finding.GetVulnerabilityOk()
	if !vulnOk {
		return nil, fmt.Errorf("missing vulnerability data for finding")
	}

	var severity Severity
	if severityStr, ok := vulnData["severity"].(string); ok {
		switch severityStr {
		case "CRITICAL":
			severity = SeverityCritical
		case "HIGH":
			severity = SeverityHigh
		case "MEDIUM":
			severity = SeverityMedium
		case "LOW":
			severity = SeverityLow
		default:
			severity = SeverityUnassigned
		}
	} else {
		// default to unassigned if severity is missing, or it is not a known value
		severity = SeverityUnassigned
	}

	var vulnId string
	if v, ok := vulnData["vulnId"].(string); ok {
		vulnId = v
	}

	var projectId string
	if p, ok := component["project"].(string); ok {
		projectId = p
	}
	var componentId string
	if c, ok := component["uuid"].(string); ok {
		componentId = c
	}
	var vulnerabilityUuid string
	if v, ok := vulnData["uuid"].(string); ok {
		vulnerabilityUuid = v
	}

	var link string
	if source, ok := vulnData["source"].(string); ok {
		switch source {
		case "NVD":
			link = fmt.Sprintf("https://nvd.nist.gov/vuln/detail/%s", vulnId)
		case "GITHUB":
			link = fmt.Sprintf("https://github.com/advisories/%s", vulnId)
		case "UBUNTU":
			link = fmt.Sprintf("https://ubuntu.com/security/CVE-%s", vulnId)
		case "OSSINDEX":
			link = fmt.Sprintf("https://ossindex.sonatype.org/vuln/%s", vulnId)
		case "DEBIAN":
			link = fmt.Sprintf("https://security-tracker.debian.org/tracker/%s", vulnId)
		}
	}

	suppressed := false
	if s, ok := analysis["isSuppressed"].(bool); ok {
		suppressed = s
	}

	title := ""
	if t, ok := vulnData["title"].(string); ok && t != "" {
		title = t
	} else if cwe, ok := vulnData["cweName"].(string); ok {
		title = cwe
	}

	desc := "unknown"
	if d, ok := vulnData["description"].(string); ok {
		desc = d
	}

	purl := ""
	if p, ok := component["purl"].(string); ok {
		purl = p
	}

	componentLatestVersion := ""
	if lv, ok := component["latestVersion"].(string); ok {
		componentLatestVersion = lv
	}

	references := map[string]string{}
	if aliases, ok := vulnData["aliases"].([]interface{}); ok {
		for _, a := range aliases {
			if alias, ok := a.(map[string]interface{}); ok {
				if cveId, ok := alias["cveId"].(string); ok {
					if ghsaId, ok := alias["ghsaId"].(string); ok {
						references[cveId] = ghsaId
					}
				}
			}
		}
	}

	return &Vulnerability{
		Package:       purl,
		Suppressed:    suppressed,
		LatestVersion: componentLatestVersion,
		Cve: &Cve{
			Id:          vulnId,
			Description: desc,
			Title:       title,
			Link:        link,
			Severity:    severity,
			References:  references,
		},
		Metadata: &VulnMetadata{
			ProjectId:         projectId,
			ComponentId:       componentId,
			VulnerabilityUuid: vulnerabilityUuid,
		},
	}, nil
}
