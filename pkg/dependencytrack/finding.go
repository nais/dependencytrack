package dependencytrack

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/nais/dependencytrack/pkg/dependencytrack/client"
)

// GetFindings Is this function lacking pagination for all findings in a project or do we not need it?
// https://github.com/DependencyTrack/dependency-track/issues/3811
// https://github.com/DependencyTrack/dependency-track/issues/4677
func (c *dependencyTrackClient) GetFindings(ctx context.Context, uuid, vulnId string, suppressed bool) ([]client.Finding, error) {
	return withAuthContextValue(c, ctx, func(tokenCtx context.Context) ([]client.Finding, error) {
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
			return nil, convertError(err, "GetFindings", resp)
		}

		return findings, nil
	})
}

// UpdateFinding updates a finding in Dependency-Track.
func (c *dependencyTrackClient) UpdateFinding(
	ctx context.Context,
	suppressedBy, reason string,
	projectId, componentId, vulnerabilityId, state string,
	suppressed bool,
) error {
	return c.withAuthContext(ctx, func(tokenCtx context.Context) error {
		comment := fmt.Sprintf("on-behalf-of:%s|suppressed:%t|state:%s|comment:%s", suppressedBy, suppressed, state, reason)
		analysisJustification := "NOT_SET"
		analysisResponse := "NOT_SET"
		analysisRequest := client.AnalysisRequest{
			Vulnerability:         vulnerabilityId,
			Component:             componentId,
			Project:               &projectId,
			AnalysisState:         &state,
			AnalysisJustification: &analysisJustification,
			AnalysisResponse:      &analysisResponse,
			AnalysisDetails:       &reason,
			Comment:               &comment,
			Suppressed:            &suppressed,
		}

		_, resp, err := c.client.AnalysisAPI.UpdateAnalysis(tokenCtx).
			AnalysisRequest(analysisRequest).
			Execute()

		if err != nil {
			return fmt.Errorf("failed to update finding: %v details: %s", err, parseErrorResponseBody(resp))
		}

		return nil
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

// GetAnalysisTrailForImage retrieves the analysis trail for a specific image component and vulnerability.
func (c *dependencyTrackClient) GetAnalysisTrailForImage(
	ctx context.Context,
	projectId, componentID, vulnerabilityId string,
) (*client.Analysis, error) {
	return withAuthContextValue(c, ctx, func(tokenCtx context.Context) (*client.Analysis, error) {
		trail, resp, err := c.client.AnalysisAPI.RetrieveAnalysis(tokenCtx).
			Project(projectId).
			Component(componentID).
			Vulnerability(vulnerabilityId).
			Execute()

		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				return nil, nil
			}
			return nil, convertError(err, "GetAnalysisTrailForImage", resp)
		}
		return trail, nil
	})
}
