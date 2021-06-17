package grpcdecode

import (
	"context"
	"github.com/al8n/kit-auth/models/requests"
	"github.com/al8n/kit-auth/authentication/register-service/internal/codec/grpccodec"
	"github.com/al8n/kit-auth/authentication/register-service/pb"
)

func RegisterByPhoneRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.RegisterByPhoneRequest)
	return requests.RegisterByPhoneRequest{
		Prefix: req.Prefix,
		Phone: req.Phone,
		OTP: req.Otp,
	}, nil
}

func RegisterByEmailRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.RegisterByEmailRequest)
	return requests.RegisterByEmailRequest{
		Email: req.Email,
		OTP: req.Otp,
	}, nil
}

func AuthenticationResponse(_ context.Context, grpcResp interface{}) (interface{}, error) {
	resp := grpcResp.(*pb.AuthenticationResponse)
	return grpccodec.AuthpbResp2Resp(*resp), nil
}

func SendOTPResponse(_ context.Context, grpcResp interface{}) (interface{}, error) {
	resp := grpcResp.(*pb.SendOTPResponse)
	return grpccodec.SendOTPpbResp2Resp(*resp), nil
}

func SendOTPToEmailRequest(_ context.Context, grpcReq interface{}) (interface{}, error)  {
	req := grpcReq.(*pb.SendOTPToEmailRequest)
	return requests.SendOTPToEmailRequest{
		Email: req.Email,
	}, nil
}

func SendOTPToPhoneRequest(_ context.Context, grpcReq interface{}) (interface{}, error)  {
	req := grpcReq.(*pb.SendOTPToPhoneRequest)
	return requests.SendOTPToPhoneRequest{
		Prefix: req.Prefix,
		Phone: req.Phone,
	}, nil
}