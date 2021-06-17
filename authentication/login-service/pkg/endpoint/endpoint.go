package endpoint

import (
	"context"
	"errors"
	"fmt"
	loginservice "github.com/al8n/kit-auth/authentication/login-service/pkg/service"
	"github.com/al8n/kit-auth/models"
	"github.com/al8n/kit-auth/models/requests"
	"github.com/al8n/kit-auth/models/responses"
	"github.com/al8n/kit-auth/utils"
	bootapi "github.com/al8n/micro-boot/api"
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


type MakeEndpointFunc = func(loginservice.Service) endpoint.Endpoint

type Set struct {
	LoginByEmailEndpoint endpoint.Endpoint
}

func (s Set) LoginByEmail(ctx context.Context, email, otp string) (token string, userInfo *models.UserInfo, err error) {
	var (
		resp interface{}
		response responses.AuthenticationResponse
	)

	resp, err = s.LoginByEmailEndpoint(ctx, requests.LoginByEmailRequest{
		Email:    email,
		OTP: otp,
	})

	response = resp.(responses.AuthenticationResponse)
	return response.Token, &response.User, utils.Str2Err(response.Error)
}

func New(svc loginservice.Service, logger log.Logger, duration map[string]metrics.Histogram, otTracer stdopentracing.Tracer, apis bootapi.APIs) (set *Set, err error) {

	set = &Set{
		LoginByEmailEndpoint:    MakeEndpoint(
			svc,
			apis[loginservice.LoginByEmailServiceName],
			logger,
			duration[loginservice.LoginByEmailServiceName],
			otTracer,
			MakeLoginByEmailEndpoint),
	}

	return
}

func MakeEndpoint(svc loginservice.Service, api bootapi.API, logger log.Logger, duration metrics.Histogram, otTracer stdopentracing.Tracer,  makeFN MakeEndpointFunc) endpoint.Endpoint {
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

func MakeLoginByEmailEndpoint(svc loginservice.Service) endpoint.Endpoint  {
	return func(ctx context.Context, request interface{}) ( response interface{}, err error) {
		var (
			req requests.LoginByEmailRequest
			token string
			user *models.UserInfo
		)

		req = request.(requests.LoginByEmailRequest)
		token, user, err = svc.LoginByEmail(ctx, req.Email, req.OTP)
		if err != nil {
			return responses.AuthenticationResponse{
				Error: err.Error(),
			}, nil
		}
		return responses.AuthenticationResponse{
			Token:    token,
			User: *user,
			Error:    "",
		}, nil
	}
}

func ErrorNoMakeFunc(name string) error {
	return errors.New(fmt.Sprintf("no endpoint make function found for %s", name))
}