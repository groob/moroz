package moroz

import (
	"compress/zlib"
	"context"
	"net/http"
	"time"

	"github.com/goccy/go-yaml"

	"github.com/go-kit/kit/endpoint"
	"github.com/groob/moroz/santa"
)

func (svc *SantaService) Postflight(ctx context.Context, machineID string, p santa.PostflightPayload) (*santa.Postflight, error) {
	return &santa.Postflight{}, nil
}

type postflightRequest struct {
	MachineID string
	payload   santa.PostflightPayload
}

type postflightResponse struct {
	*santa.Postflight
	Err error `json:"error,omitempty"`
}

func (r postflightResponse) Failed() error { return r.Err }

func makePostflightEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(postflightRequest)
		postflight, err := svc.Postflight(ctx, req.MachineID, req.payload)
		if err != nil {
			return postflightResponse{Err: err}, nil
		}
		return postflightResponse{Postflight: postflight}, nil
	}
}

func decodePostflightRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	zr, err := zlib.NewReader(r.Body)
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	defer r.Body.Close()
	id, err := machineIDFromRequest(r)
	if err != nil {
		return nil, err
	}
	req := postflightRequest{MachineID: id}
	if err := yaml.NewDecoder(zr).Decode(&req.payload); err != nil {
		return nil, err
	}
	return req, nil
}

func (mw logmw) Postflight(ctx context.Context, machineID string, p santa.PostflightPayload) (pf *santa.Postflight, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "Postflight",
			"machine_id", machineID,
			"postflight_payload", p,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	pf, err = mw.next.Postflight(ctx, machineID, p)
	return
}
