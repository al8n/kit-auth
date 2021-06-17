package grpcencode

import (
	"context"
	"github.com/al8n/kit-auth/authentication/register-service/internal/codec/grpccodec"
	"github.com/al8n/kit-auth/authentication/register-service/pb"
	"github.com/al8n/kit-auth/models/requests"
	"github.com/al8n/kit-auth/models/responses"
	"github.com/al8n/kit-auth/utils"
	utilstransport "github.com/al8n/kit-auth/utils/transport"
)

func RegisterByPhoneRequest(_ context.Context, request interface{}) ( interface{}, error)  {
	req, ok := request.(requests.RegisterByPhoneRequest)
	if !ok {
		return nil, utilstransport.ErrorCodecCasting("RegisterByPhone", utilstransport.Request,utilstransport.GRPC)
	}
	return grpccodec.RegisterByPhoneReq2pbReq(req), nil
}

func RegisterByEmailRequest(_ context.Context, request interface{}) ( interface{}, error)  {
	req, ok := request.(requests.RegisterByEmailRequest)
	if !ok {
		return nil, utilstransport.ErrorCodecCasting("RegisterByEmail", utilstransport.Request,utilstransport.GRPC)
	}
	return grpccodec.RegisterByEmailReq2pbReq(req), nil
}


func SendOTPToPhoneRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(requests.SendOTPToPhoneRequest)
	if !ok {
		return nil, utilstransport.ErrorCodecCasting("SendOTPToPhone", utilstransport.Request, utilstransport.GRPC)
	}
	return grpccodec.SendOTPToPhoneReq2pbReq(req), nil
}

func SendOTPToEmailRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(requests.SendOTPToEmailRequest)
	if !ok {
		return nil, utilstransport.ErrorCodecCasting("SendOTPToEmail", utilstransport.Request, utilstransport.GRPC)
	}
	return grpccodec.SendOTPToEmailReq2pbReq(req), nil
}

func AuthenticationResponse(_ context.Context, resp interface{}) (interface{}, error)  {
	res, ok := resp.(responses.AuthenticationResponse)
	if !ok {
		return nil, utilstransport.ErrorCodecCasting("AuthenticationResponse",
			utilstransport.Response,
			utilstransport.GRPC)
	}
	if res.Error != "" {
		return nil, utils.Str2Err(res.Error)
	}

	return grpccodec.AuthResp2pbResp(res), nil
}

func SendOTPResponse(_ context.Context, resp interface{}) (interface{}, error)  {
	pbReply := &pb.SendOTPResponse{}
	res, ok := resp.(responses.SendOTPResponse)
	if !ok {
		return nil, utilstransport.ErrorCodecCasting("SendOTPResponse", utilstransport.Response, utilstransport.GRPC)
	}
	if res.Error != "" {
		return nil, utils.Str2Err(res.Error)
	}

	pbReply.Otp = res.OTP
	pbReply.Error = res.Error
	return pbReply, nil
}