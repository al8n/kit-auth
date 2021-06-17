package decode

import (
	"context"
	"encoding/json"
	"github.com/al8n/kit-auth/models/requests"
	"github.com/al8n/kit-auth/models/responses"
	"net/http"
)

func AuthenticationResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response responses.AuthenticationResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

func RegisterByEmailRequest(_ context.Context, req *http.Request) (interface{}, error)  {
	var request requests.RegisterByEmailRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}


func LoginByEmailRequest(_ context.Context, req *http.Request) (interface{}, error)  {
	var request requests.LoginByEmailRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}


