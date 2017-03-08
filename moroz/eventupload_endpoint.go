package moroz

import (
	"context"
	"io"

	"github.com/go-kit/kit/endpoint"
	"github.com/groob/moroz/santa"
)

type eventRequest struct {
	MachineID string
	eventBody io.ReadCloser
}

type eventResponse struct {
	Err error `json:"error,omitempty"`
}

func (r eventResponse) error() error { return r.Err }

func makeEventUploadEndpoint(svc santa.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(eventRequest)
		err := svc.UploadEvent(ctx, req.MachineID, req.eventBody)
		if err != nil {
			return eventResponse{Err: err}, nil
		}
		return eventResponse{}, nil
	}
}
