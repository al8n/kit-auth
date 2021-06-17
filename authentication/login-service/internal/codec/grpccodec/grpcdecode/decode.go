package grpcdecode

import (
	"context"
	"github.com/al8n/kit-auth/authentication/login-service/internal/codec/grpccodec"
	"github.com/al8n/kit-auth/authentication/login-service/pb"
	"github.com/al8n/kit-auth/models/requests"
)

func AuthenticationResponse(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.AuthenticationResponse)
	return grpccodec.AuthpbResp2Resp(*req), nil
}

func LoginByEmailRequest(_ context.Context, grpcReq interface{}) (interface{}, error)  {
	req := grpcReq.(*pb.LoginByEmailRequest)
	return requests.LoginByEmailRequest{
		Email: req.Email,
		Password: req.Password,
	}, nil
}