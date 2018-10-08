package moroz

import (
	"compress/zlib"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"

	"github.com/groob/moroz/santa"
)

func (svc *SantaService) UploadEvent(ctx context.Context, machineID string, events []santa.EventPayload) error {
	for _, ev := range events {
		eventDir := filepath.Join(svc.eventDir, ev.FileSHA, machineID)
		if err := os.MkdirAll(eventDir, 0700); err != nil {
			return errors.Wrapf(err, "create event directory %s", eventDir)
		}

		eventPath := filepath.Join(eventDir, fmt.Sprintf("%f.json", ev.UnixTime))
		if err := ioutil.WriteFile(eventPath, ev.Content, 0644); err != nil {
			return errors.Wrapf(err, "write event to path %s", eventPath)
		}
	}
	return nil
}

type eventRequest struct {
	MachineID string
	events    []santa.EventPayload
}

type eventResponse struct {
	Err error
}

func (r eventResponse) Failed() error { return r.Err }

func makeEventUploadEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(eventRequest)
		err := svc.UploadEvent(ctx, req.MachineID, req.events)
		return eventResponse{Err: err}, nil
	}
}

func decodeEventUpload(ctx context.Context, r *http.Request) (interface{}, error) {
	// santa sends zlib compressed payloads
	zr, err := zlib.NewReader(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "create zlib reader to decode event upload")
	}
	defer zr.Close()

	id, err := machineIDFromRequest(r)
	if err != nil {
		return nil, errors.Wrap(err, "get machine ID from event upload URL")
	}

	// decode the JSON into individual log events.
	var eventPayload = struct {
		Events []json.RawMessage `json:"events"`
	}{}

	if err := json.NewDecoder(zr).Decode(&eventPayload); err != nil {
		return nil, errors.Wrap(err, "decoding event upload request json")
	}

	var events []santa.EventPayload
	for _, ev := range eventPayload.Events {
		var payload santa.EventPayload
		if err := json.Unmarshal(ev, &payload); err != nil {
			return nil, errors.Wrap(err, "decoding file sha from event upload json")
		}
		payload.Content = ev
		events = append(events, payload)
	}

	req := eventRequest{MachineID: id, events: events}
	return req, nil
}

func (mw logmw) UploadEvent(ctx context.Context, machineID string, events []santa.EventPayload) (err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "UploadEvent",
			"machine_id", machineID,
			"event_count", len(events),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.next.UploadEvent(ctx, machineID, events)
	return
}
