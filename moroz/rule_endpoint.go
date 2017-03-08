package moroz

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/groob/moroz/santa"
)

type ruleRequest struct {
	MachineID string
}

type rulesResponse struct {
	Rules []santa.Rule `json:"rules"`
	Err   error        `json:"error,omitempty"`
}

func (r rulesResponse) error() error { return r.Err }

func makeRuleDownloadEndpoint(svc santa.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ruleRequest)
		rules, err := svc.RuleDownload(ctx, req.MachineID)
		if err != nil {
			return rulesResponse{Err: err}, nil
		}
		return rulesResponse{Rules: rules}, nil
	}
}
