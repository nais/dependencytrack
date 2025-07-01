package dependencytrack

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nais/dependencytrack/pkg/dependencytrack/client"
)

type AnalysisRequest struct {
	SuppressedBy    string
	Reason          string
	ProjectId       string
	ComponentId     string
	VulnerabilityId string
	State           string
	Suppressed      bool
}

type Analysis struct {
	AnalysisState    string            `json:"analysisState"`
	AnalysisResponse string            `json:"analysisResponse"`
	AnalysisDetails  string            `json:"analysisDetails"`
	AnalysisComments []AnalysisComment `json:"analysisComments,omitempty"`
	IsSuppressed     *bool             `json:"isSuppressed,omitempty"`
}

type AnalysisComment struct {
	Comment   string  `json:"comment"`
	Commenter *string `json:"commenter,omitempty"`
}

func (c *dependencyTrackClient) UpdateFinding(ctx context.Context, request AnalysisRequest) error {
	return c.withAuthContext(ctx, func(authCtx context.Context) error {
		comment := fmt.Sprintf("on-behalf-of:%s|suppressed:%t|state:%s|comment:%s", request.SuppressedBy, request.Suppressed, request.State, request.Reason)
		analysisJustification := "NOT_SET"
		analysisResponse := "NOT_SET"
		analysisRequest := client.AnalysisRequest{
			Vulnerability:         request.VulnerabilityId,
			Component:             request.ComponentId,
			Project:               &request.ProjectId,
			AnalysisState:         &request.State,
			AnalysisJustification: &analysisJustification,
			AnalysisResponse:      &analysisResponse,
			AnalysisDetails:       &request.Reason,
			Comment:               &comment,
			Suppressed:            &request.Suppressed,
		}

		_, resp, err := c.client.AnalysisAPI.UpdateAnalysis(authCtx).
			AnalysisRequest(analysisRequest).
			Execute()

		if err != nil {
			return fmt.Errorf("failed to update finding: %v details: %s", err, parseErrorResponseBody(resp))
		}

		return nil
	})
}

// GetAnalysisTrailForImage retrieves the analysis trail for a specific image component and vulnerability.
func (c *dependencyTrackClient) GetAnalysisTrailForImage(
	ctx context.Context,
	projectId, componentID, vulnerabilityId string,
) (*Analysis, error) {
	return withAuthContextValue(c, ctx, func(tokenCtx context.Context) (*Analysis, error) {
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
		return parseAnalysisResponse(trail)
	})
}

func parseAnalysisResponse(resp *client.Analysis) (*Analysis, error) {
	if resp == nil {
		return nil, nil
	}

	analysis := &Analysis{
		AnalysisState:    resp.AnalysisState,
		AnalysisResponse: resp.AnalysisResponse,
		AnalysisDetails:  resp.AnalysisDetails,
		IsSuppressed:     resp.IsSuppressed,
	}

	if resp.AnalysisComments != nil {
		for _, comment := range resp.AnalysisComments {
			analysisComment := AnalysisComment{
				Comment:   comment.Comment,
				Commenter: comment.Commenter,
			}
			analysis.AnalysisComments = append(analysis.AnalysisComments, analysisComment)
		}
	}

	return analysis, nil
}
