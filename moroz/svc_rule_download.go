package moroz

import (
	"context"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/groob/moroz/santa"
)

func (svc *SantaService) RuleDownload(ctx context.Context, machineID string) ([]santa.Rule, error) {
	config := svc.config(machineID)
	return config.Rules, nil
}

func (svc *SantaService) config(machineID string) *santa.Config {
	var config *santa.Config
	var err error
	config, err = svc.repo.Config(machineID)
	if err != nil {
		config = svc.global
	}
	return config
}

type ruleRequest struct {
	MachineID string
}

type rulesResponse struct {
	Rules []santa.Rule `json:"rules"`
	Err   error        `json:"error,omitempty"`
}

func (r rulesResponse) Failed() error { return r.Err }

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

func decodeRuleRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := machineIDFromRequest(r)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	req := ruleRequest{MachineID: id}
	return req, nil
}

func (mw logmw) RuleDownload(ctx context.Context, machineID string) (rules []santa.Rule, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "RuleDownload",
			"machine_id", machineID,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	rules, err = mw.next.RuleDownload(ctx, machineID)
	return
}
