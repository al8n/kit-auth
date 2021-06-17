package httpdecode

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/al8n/kit-auth/models/requests"
	"github.com/al8n/kit-auth/models/responses"
	"net/http"
)

func RegisterByEmailRequest(_ context.Context, r *http.Request) (interface{}, error)  {
	var req requests.RegisterByEmailRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func RegisterByPhoneRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req requests.RegisterByPhoneRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func SendOTPToPhoneRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req requests.SendOTPToPhoneRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func SendOTPToEmailRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req requests.SendOTPToEmailRequest
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

func SendOTPResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp responses.SendOTPResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}





