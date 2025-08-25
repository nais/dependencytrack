package dependencytrack

import (
	"context"
)

// IsTaskInProgress Determines if there are any tasks associated with the token that are being processed, or in the queue to be processed.
func (c *dependencyTrackClient) IsTaskInProgress(ctx context.Context, uuid string) (bool, error) {
	return withAuthContextValue(c.auth, ctx, func(tokenCtx context.Context) (bool, error) {
		tokenProcess, resp, err := c.client.EventAPI.IsTokenBeingProcessed1(tokenCtx, uuid).Execute()
		if err != nil {
			return false, convertError(err, "IsTaskInProgress", resp)
		}
		return tokenProcess.GetProcessing(), nil
	})
}
