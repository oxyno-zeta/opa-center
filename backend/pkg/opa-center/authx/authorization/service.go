package authorization

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/oxyno-zeta/opa-center/pkg/opa-center/authx/authentication"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/authx/models"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/common/errors"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/config"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/log"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/tracing"
)

type service struct {
	cfgManager config.Manager
}

type generalInputOPA struct {
	Input *generalInputDataOPA `json:"input"`
}

type generalInputDataOPA struct {
	User    *models.OIDCUser  `json:"user"`
	Tags    map[string]string `json:"tags"`
	Request *generalDataOPA   `json:"data"`
}

type generalDataOPA struct {
	Action   string `json:"action"`
	Resource string `json:"resource"`
}

type restInputOPA struct {
	Input *restInputDataOPA `json:"input"`
}

type restInputDataOPA struct {
	User    *models.OIDCUser    `json:"user"`
	Request *httpRequestDataOPA `json:"request"`
	Tags    map[string]string   `json:"tags"`
}

type httpRequestDataOPA struct {
	Method     string            `json:"method"`
	Protocol   string            `json:"protocol"`
	Headers    map[string]string `json:"headers"`
	RemoteAddr string            `json:"remoteAddr"`
	Scheme     string            `json:"scheme"`
	Host       string            `json:"host"`
	ParsedPath []string          `json:"parsed_path"`
	Path       string            `json:"path"`
}

type opaAnswer struct {
	Result bool `json:"result"`
}

func (s *service) IsAuthorized(ctx context.Context, action, resource string) (bool, error) {
	// Get logger
	logger := log.GetLoggerFromContext(ctx)
	// Get configuration to check that authorization can be calculated
	cfg := s.cfgManager.GetConfig().OPAServerAuthorization
	if cfg == nil {
		// Configuration doesn't exists, authorization is given
		return true, nil
	}

	// Get user from context
	user := authentication.GetAuthenticatedUserFromContext(ctx)

	// Create opa input
	input := &generalInputOPA{
		Input: &generalInputDataOPA{
			User: user,
			Tags: cfg.Tags,
			Request: &generalDataOPA{
				Action:   action,
				Resource: resource,
			},
		},
	}
	// Json encode body
	bb, err := json.Marshal(input)
	if err != nil {
		return false, err
	}

	authorized, err := s.requestOPAServer(ctx, cfg, bb)
	// Check error
	if err != nil {
		return false, err
	}

	// Check if user isn't authorized
	if !authorized {
		logger.Infof("User %s not authorized for action %s on resource %s", user.GetIdentifier(), action, resource)

		return false, nil
	}

	logger.Infof("User %s authorized for action %s on resource %s", user.GetIdentifier(), action, resource)

	return true, nil
}

func (s *service) CheckAuthorized(ctx context.Context, action, resource string) error {
	// Call is authorized
	res, err := s.IsAuthorized(ctx, action, resource)
	// Check error
	if err != nil {
		return err
	}

	// Check not authorized
	if !res {
		return errors.NewForbiddenError("forbidden")
	}

	return nil
}

func (s *service) requestOPAServer(ctx context.Context, opaCfg *config.OPAServerAuthorization, body []byte) (bool, error) {
	// Get trace from context
	trace := tracing.GetTraceFromContext(ctx)
	// Generate child trace
	childTrace := trace.GetChildTrace("opa-server.request")
	defer childTrace.Finish()
	// Add data
	childTrace.SetTag("opa.uri", opaCfg.URL)

	// Change NewRequest to NewRequestWithContext and pass context it
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, opaCfg.URL, bytes.NewBuffer(body))
	if err != nil {
		return false, err
	}
	// Add content type
	req.Header.Add("Content-Type", "application/json")
	// Making request to OPA server
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	// Defer closing body
	defer resp.Body.Close()

	// Prepare answer
	var answer opaAnswer
	// Decode answer
	err = json.NewDecoder(resp.Body).Decode(&answer)
	if err != nil {
		return false, err
	}

	return answer.Result, nil
}

func (s *service) isRequestAuthorized(req *http.Request, oidcUser *models.OIDCUser) (bool, error) {
	// Get configuration
	opaServerCfg := s.cfgManager.GetConfig().OPAServerAuthorization

	// Transform headers into map
	headers := make(map[string]string)
	for k, v := range req.Header {
		headers[strings.ToLower(k)] = v[0]
	}
	// Parse path
	parsedPath := deleteEmpty(strings.Split(req.RequestURI, "/"))
	// Calculate scheme
	scheme := "http"
	if req.TLS != nil {
		scheme = "https"
	}
	// Generate OPA Server input data
	input := &restInputOPA{
		Input: &restInputDataOPA{
			User: oidcUser,
			Tags: opaServerCfg.Tags,
			Request: &httpRequestDataOPA{
				Method:     req.Method,
				Protocol:   req.Proto,
				Headers:    headers,
				RemoteAddr: req.RemoteAddr,
				Scheme:     scheme,
				Host:       req.Host,
				ParsedPath: parsedPath,
				Path:       req.RequestURI,
			},
		},
	}
	// Json encode body
	bb, err := json.Marshal(input)
	if err != nil {
		return false, err
	}

	return s.requestOPAServer(req.Context(), opaServerCfg, bb)
}

func deleteEmpty(s []string) []string {
	var r []string

	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}

	return r
}
