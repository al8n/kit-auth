package httpdecode

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/al8n/kit-auth/models/requests"
	"github.com/al8n/kit-auth/models/responses"
	"net/http"
)

func LoginByEmailRequest(_ context.Context, r *http.Request) (interface{}, error)  {
	var req requests.LoginByEmailRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func AuthenticationResponse(_ context.Context, r *http.Response) (interface{}, error)  {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp responses.AuthenticationResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}




