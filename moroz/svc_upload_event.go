package moroz

import (
	"compress/zlib"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"

	"github.com/groob/moroz/santa"
)

func (svc *SantaService) UploadEvent(ctx context.Context, machineID string, events []santa.EventPayload) error {
	if !svc.flPersistEvents {
		return nil
	}
	for _, ev := range events {
		eventDir := filepath.Join(svc.eventDir, ev.FileSHA, machineID)
		if err := os.MkdirAll(eventDir, 0700); err != nil {
			return errors.Wrapf(err, "create event directory %s", eventDir)
		}

		eventPath := filepath.Join(eventDir, fmt.Sprintf("%f.json", ev.UnixTime))

		eventInfoJSON, err := json.Marshal(ev.EventInfo)
		if err != nil {
			return errors.Wrap(err, "marshal event info to json")
		}

		if err := os.WriteFile(eventPath, eventInfoJSON, 0644); err != nil {
			return errors.Wrapf(err, "write event to path %s", eventPath)
		}
	}
	return nil
}

// UploadFileAccessEvents persists File Access Authorization (FAA) events sent in
// an /eventupload request. They have no file SHA, so they are stored under a
// dedicated file_access/{machineID} directory keyed by access time.
func (svc *SantaService) UploadFileAccessEvents(ctx context.Context, machineID string, events []santa.FileAccessEvent) error {
	if !svc.flPersistEvents {
		return nil
	}
	for _, ev := range events {
		eventDir := filepath.Join(svc.eventDir, "file_access", machineID)
		if err := os.MkdirAll(eventDir, 0700); err != nil {
			return errors.Wrapf(err, "create file access event directory %s", eventDir)
		}

		eventPath := filepath.Join(eventDir, fmt.Sprintf("%f.json", ev.AccessTime))

		eventJSON, err := json.Marshal(ev)
		if err != nil {
			return errors.Wrap(err, "marshal file access event to json")
		}

		if err := os.WriteFile(eventPath, eventJSON, 0644); err != nil {
			return errors.Wrapf(err, "write file access event to path %s", eventPath)
		}
	}
	return nil
}

type eventRequest struct {
	MachineID        string
	events           []santa.EventPayload
	fileAccessEvents []santa.FileAccessEvent
}

type eventResponse struct {
	Err error
}

func (r eventResponse) Failed() error { return r.Err }

func makeEventUploadEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(eventRequest)
		if err := svc.UploadEvent(ctx, req.MachineID, req.events); err != nil {
			return eventResponse{Err: err}, nil
		}
		err := svc.UploadFileAccessEvents(ctx, req.MachineID, req.fileAccessEvents)
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
	var eventPayload santa.EventUploadRequest

	if err := json.NewDecoder(zr).Decode(&eventPayload); err != nil {
		return nil, errors.Wrap(err, "decoding event upload request json")
	}

	var events []santa.EventPayload
	for _, ev := range eventPayload.Events {
		var payload santa.EventPayload
		payload.EventInfo = ev
		payload.FileSHA = ev.FileSHA256
		payload.UnixTime = ev.ExecutionTime
		events = append(events, payload)
	}

	req := eventRequest{
		MachineID:        id,
		events:           events,
		fileAccessEvents: eventPayload.FileAccessEvents,
	}
	return req, nil
}

func (mw logmw) UploadEvent(ctx context.Context, machineID string, events []santa.EventPayload) (err error) {
	defer func(begin time.Time) {
		for _, ev := range events {
			_ = mw.logger.Log(
				"method", "UploadEvent",
				"machine_id", machineID,
				"event", ev.EventInfo,
				"err", err,
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	err = mw.next.UploadEvent(ctx, machineID, events)
	return
}

func (mw logmw) UploadFileAccessEvents(ctx context.Context, machineID string, events []santa.FileAccessEvent) (err error) {
	defer func(begin time.Time) {
		for _, ev := range events {
			_ = mw.logger.Log(
				"method", "UploadFileAccessEvents",
				"machine_id", machineID,
				"rule_name", ev.RuleName,
				"target", ev.Target,
				"decision", ev.Decision,
				"err", err,
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	err = mw.next.UploadFileAccessEvents(ctx, machineID, events)
	return
}
