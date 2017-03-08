package moroz

import (
	"compress/zlib"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

func decodePreflightRequest(ctx context.Context, r *http.Request) (interface{}, error) {
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
	req := preflightRequest{MachineID: id}
	if err := json.NewDecoder(zr).Decode(&req.payload); err != nil {
		return nil, err
	}
	return req, nil
}

// errBadRoute is used for mux errors
var errBadRoute = errors.New("bad route")

func machineIDFromRequest(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return "", errBadRoute
	}
	return id, nil
}
