package moroz

import (
	"compress/zlib"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/groob/moroz/santa"
)

func (svc *SantaService) UploadEvent(ctx context.Context, machineID string, events santa.EventsList) error {
	when := time.Now().UTC().Format(time.RFC3339)
	for _, event := range events.Events {
		event_line := santa.EventLine{machineID, when, event}
		b, err := json.Marshal(event_line)
		if err != nil {
			// this can't happen?
			return err
		}
		svc.eventWriter.Write(append(b, '\n'))
	}
	return nil
}

type eventRequest struct {
	MachineID string
	Events    santa.EventsList
	// eventBody io.ReadCloser
}

type eventResponse struct {
	Err error `json:"error,omitempty"`
}

func (r eventResponse) Failed() error { return r.Err }

func makeEventUploadEndpoint(svc santa.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(eventRequest)
		err := svc.UploadEvent(ctx, req.MachineID, req.Events)
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
	req := eventRequest{MachineID: id}
	if err := json.NewDecoder(zr).Decode(&req.Events); err != nil {
		return nil, err
	}
	return req, nil
}

func (mw logmw) UploadEvent(ctx context.Context, machineID string, body santa.EventsList) (err error) {
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
