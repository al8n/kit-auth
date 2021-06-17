package grpcencode

import (
	"context"
	"github.com/al8n/kit-auth/authentication/login-service/internal/codec/grpccodec"
	"github.com/al8n/kit-auth/models/requests"
	"github.com/al8n/kit-auth/models/responses"
	"github.com/al8n/kit-auth/utils"
	utilstransport "github.com/al8n/kit-auth/utils/transport"
)

func AuthenticationResponse(_ context.Context, resp interface{}) (interface{}, error)  {
	res, ok := resp.(responses.AuthenticationResponse)
	if !ok {
		return nil, utilstransport.ErrorCodecCasting("AuthenticationResponse", utilstransport.Response,utilstransport.GRPC)
	}
	if res.Error != "" {
		return nil, utils.Str2Err(res.Error)
	}

	return grpccodec.AuthResp2pbResp(res), nil
}

func LoginByEmailRequest(_ context.Context, request interface{}) (interface{}, error)  {
	req, ok := request.(requests.LoginByEmailRequest)
	if !ok {
		return nil, utilstransport.ErrorCodecCasting("LoginByEmailRequest", utilstransport.Request, utilstransport.GRPC)
	}
	return grpccodec.LoginByEmailReq2pbReq(req), nil
}