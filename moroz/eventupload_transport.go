package moroz

import (
	"compress/zlib"
	"context"
	"net/http"
)

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
