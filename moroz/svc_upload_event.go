package moroz

import (
	"compress/zlib"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/groob/moroz/santa"
)

func (svc *SantaService) UploadEvent(ctx context.Context, machineID string, body io.ReadCloser) error {
	_, err := io.Copy(svc.eventWriter, body)
	defer body.Close()
	return err
}

type eventRequest struct {
	MachineID string
	eventBody io.ReadCloser
}

type eventResponse struct {
	Err error `json:"error,omitempty"`
}

func (r eventResponse) Failed() error { return r.Err }

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

func decodeEventUpload(ctx context.Context, r *http.Request) (interface{}, error) {
	zr, err := zlib.NewReader(r.Body)
	if err != nil {
		return nil, err
	}
	id, err := machineIDFromRequest(r)
	if err != nil {
		return nil, err
	}
	req := eventRequest{MachineID: id, eventBody: zr}
	return req, nil
}

func (mw logmw) UploadEvent(ctx context.Context, machineID string, body io.ReadCloser) (err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "UploadEvent",
			"machine_id", machineID,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.next.UploadEvent(ctx, machineID, body)
	return
}
