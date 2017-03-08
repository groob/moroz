package moroz

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/groob/moroz/santa"
)

type preflightRequest struct {
	MachineID string
	payload   santa.PreflightPayload
}

type preflightResponse struct {
	*santa.Preflight
	Err error `json:"error,omitempty"`
}

func (r preflightResponse) error() error { return r.Err }

func makePreflightEndpoint(svc santa.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(preflightRequest)
		preflight, err := svc.Preflight(ctx, req.MachineID, req.payload)
		if err != nil {
			return preflightResponse{Err: err}, nil
		}
		return preflightResponse{Preflight: preflight}, nil
	}
}
