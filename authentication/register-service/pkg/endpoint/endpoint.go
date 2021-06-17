package endpoint

import (
	"context"
	"errors"
	"fmt"
	bootapi "github.com/al8n/micro-boot/api"
	authservice "github.com/al8n/kit-auth/authentication/register-service/pkg/service"
	"github.com/al8n/kit-auth/models"
	"github.com/al8n/kit-auth/models/requests"
	"github.com/al8n/kit-auth/models/responses"
	"github.com/al8n/kit-auth/utils"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
)


type MakeEndpointFunc = func(authservice.Service) endpoint.Endpoint

type Set struct {
	RegisterByEmailEndpoint endpoint.Endpoint
	RegisterByPhoneEndpoint endpoint.Endpoint
	SendOTPToEmailEndpoint endpoint.Endpoint
	SendOTPToPhoneEndpoint endpoint.Endpoint
}

func (s Set) RegisterByPhone(ctx context.Context, prefix, phone, otp string) (token string, userInfo *models.UserInfo, err error) {
	var (
		resp interface{}
		response responses.AuthenticationResponse
	)

	resp, err = s.RegisterByPhoneEndpoint(ctx, requests.RegisterByPhoneRequest{
		Prefix:    prefix,
		Phone: phone,
		OTP:  otp,
	})

	response = resp.(responses.AuthenticationResponse)
	return response.Token, &response.User, utils.Str2Err(response.Error)
}

func (s Set) SendOTPToPhone(ctx context.Context, prefix, phone string) (otp string, err error) {
	var (
		resp interface{}
		response responses.SendOTPResponse
	)

	resp, err = s.SendOTPToPhoneEndpoint(ctx, requests.SendOTPToPhoneRequest{
		Prefix:    prefix,
		Phone: phone,
	})

	response = resp.(responses.SendOTPResponse)
	return response.OTP, utils.Str2Err(response.Error)
}

func (s Set) SendOTPToEmail(ctx context.Context, email string) (otp string, err error) {
	var (
		resp interface{}
		response responses.SendOTPResponse
	)

	resp, err = s.SendOTPToEmailEndpoint(ctx, requests.SendOTPToEmailRequest{
		Email:    email,
	})

	response = resp.(responses.SendOTPResponse)
	return response.OTP, utils.Str2Err(response.Error)
}

func (s Set) RegisterByEmail(ctx context.Context, email, captcha string) (token string, userInfo *models.UserInfo, err error) {
	var (
		resp interface{}
		response responses.AuthenticationResponse
	)

	resp, err = s.RegisterByEmailEndpoint(ctx, requests.RegisterByEmailRequest{
		Email:    email,
		OTP:  captcha,
	})

	response = resp.(responses.AuthenticationResponse)
	return response.Token, &response.User, utils.Str2Err(response.Error)
}

func New(svc authservice.Service, logger log.Logger, duration map[string]metrics.Histogram, otTracer stdopentracing.Tracer, apis bootapi.APIs) (set *Set, err error) {

	set = &Set{
		SendOTPToEmailEndpoint: MakeEndpoint(
			svc,
			apis[authservice.SendOTPToEmail],
			logger,
			duration[authservice.SendOTPToEmail],
			otTracer,
			MakeSendOTPToEmailEndpoint),
		RegisterByEmailEndpoint: MakeEndpoint(
			svc,
			apis[authservice.RegisterByEmailServiceName],
			logger,
			duration[authservice.RegisterByEmailServiceName],
			otTracer,
			MakeRegisterByEmailEndpoint),
		SendOTPToPhoneEndpoint: MakeEndpoint(
			svc,
			apis[authservice.SendOTPToPhone],
			logger,
			duration[authservice.SendOTPToPhone],
			otTracer,
			MakeSendOTPToPhoneEndpoint),
		RegisterByPhoneEndpoint: MakeEndpoint(
			svc,
			apis[authservice.RegisterByPhoneServiceName],
			logger,
			duration[authservice.RegisterByPhoneServiceName],
			otTracer,
			MakeRegisterByPhoneEndpoint),

	}

	return
}

func MakeEndpoint(svc authservice.Service, api bootapi.API, logger log.Logger, duration metrics.Histogram, otTracer stdopentracing.Tracer,  makeFN MakeEndpointFunc) endpoint.Endpoint {
	var ep endpoint.Endpoint
	{
		ep = makeFN(svc)
		// RegisterByEmail is limited to 1000 requests per second with burst of 1 request.
		// Note, rate is defined as a time interval between requests.
		ep = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(api.RateLimit.Duration), api.RateLimit.Delta))(ep)
		ep = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(api.Breaker.Standardize()))(ep)
		ep = opentracing.TraceServer(otTracer, api.Name)(ep)
		ep = LoggingMiddleware(log.With(logger, api.GetGoKitLoggerKVs()))(ep)
		ep = InstrumentingMiddleware(duration.With(api.Instrument...))(ep)
	}
	return ep
}

func MakeSendOTPToEmailEndpoint(svc authservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var (
			req requests.SendOTPToEmailRequest
			captcha string
		)

		req = request.(requests.SendOTPToEmailRequest)
		captcha, err = svc.SendOTPToEmail(ctx, req.Email)
		if err != nil {
			return responses.SendOTPResponse{
				Error: err.Error(),
			}, nil
		}

		return responses.SendOTPResponse{
			OTP: captcha,
		}, nil
	}
}

func MakeRegisterByEmailEndpoint(svc authservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var (
			req requests.RegisterByEmailRequest
			token string
			user *models.UserInfo
		)

		req = request.(requests.RegisterByEmailRequest)
		token, user, err = svc.RegisterByEmail(ctx, req.Email, req.OTP)
		if err != nil {
			return responses.AuthenticationResponse{
				Error: err.Error(),
			}, nil
		}

		return responses.AuthenticationResponse{
			Token:    token,
			User:      *user,
			Error:    "",
		}, nil
	}
}


func MakeSendOTPToPhoneEndpoint(svc authservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var (
			req requests.SendOTPToPhoneRequest
			captcha string
		)

		req = request.(requests.SendOTPToPhoneRequest)
		captcha, err = svc.SendOTPToPhone(ctx, req.Prefix, req.Phone)
		if err != nil {
			return responses.SendOTPResponse{
				Error: err.Error(),
			}, nil
		}

		return responses.SendOTPResponse{
			OTP: captcha,
		}, nil
	}
}

func MakeRegisterByPhoneEndpoint(svc authservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var (
			req requests.RegisterByPhoneRequest
			token string
			user *models.UserInfo
		)

		req = request.(requests.RegisterByPhoneRequest)
		token, user, err = svc.RegisterByPhone(ctx, req.Prefix, req.Phone, req.OTP)
		if err != nil {
			return responses.AuthenticationResponse{
				Error: err.Error(),
			}, nil
		}

		return responses.AuthenticationResponse{
			Token:    token,
			User:      *user,
			Error:    "",
		}, nil
	}
}


func ErrorNoMakeFunc(name string) error {
	return errors.New(fmt.Sprintf("no endpoint make function found for %s", name))
}