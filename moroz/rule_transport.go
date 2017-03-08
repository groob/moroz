package moroz

import (
	"context"
	"net/http"
)

func decodeRuleRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := machineIDFromRequest(r)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	req := ruleRequest{MachineID: id}
	return req, nil
}
