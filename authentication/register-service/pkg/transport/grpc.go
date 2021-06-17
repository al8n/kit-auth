package transport

import (
	"context"
	bootapi "github.com/al8n/micro-boot/api"
	"github.com/al8n/kit-auth/authentication/register-service/internal/codec/grpccodec/grpcdecode"
	"github.com/al8n/kit-auth/authentication/register-service/internal/codec/grpccodec/grpcencode"
	"github.com/al8n/kit-auth/authentication/register-service/pb"
	serviceendpoints "github.com/al8n/kit-auth/authentication/register-service/pkg/endpoint"
	"github.com/al8n/kit-auth/authentication/register-service/pkg/service"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/transport"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

const pbName = "pb.Register"

type GRPCServer struct {
	registerByEmail grpctransport.Handler
	registerByPhone grpctransport.Handler
	sendOTPToEmail grpctransport.Handler
	sendOTPToPhone grpctransport.Handler
}

// NewGRPCServer makes a set of endpoints available as a gRPC AuthenticationServer.
func NewGRPCServer(set serviceendpoints.Set, otTracer stdopentracing.Tracer, logger log.Logger, apis bootapi.APIs) pb.RegisterServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	return &GRPCServer{
		registerByEmail: grpctransport.NewServer(
			set.RegisterByEmailEndpoint,
			grpcdecode.RegisterByEmailRequest,
			grpcencode.AuthenticationResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, service.RegisterByEmailServiceName, logger)))...,
		),
		registerByPhone: grpctransport.NewServer(
			set.RegisterByPhoneEndpoint,
			grpcdecode.RegisterByPhoneRequest,
			grpcencode.AuthenticationResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, service.RegisterByPhoneServiceName, logger)))...,
		),
		sendOTPToPhone: grpctransport.NewServer(
			set.SendOTPToPhoneEndpoint,
			grpcdecode.SendOTPToPhoneRequest,
			grpcencode.SendOTPResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, service.SendOTPToPhone, logger)))...,
		),
		sendOTPToEmail: grpctransport.NewServer(
			set.SendOTPToEmailEndpoint,
			grpcdecode.SendOTPToEmailRequest,
			grpcencode.SendOTPResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, service.SendOTPToEmail, logger)))...,
		),
	}
}

// NewGRPCClient returns an AuthenticationService backed by a gRPC server at the other end
// of the conn. The caller is responsible for constructing the conn, and
// eventually closing the underlying transport. We bake-in certain middlewares,
// implementing the client library pattern.
func NewGRPCClient(conn *grpc.ClientConn, otTracer stdopentracing.Tracer, logger log.Logger, apis bootapi.APIs) service.Service {

	// Each individual endpoint is an grpc/transport.Client (which implements
	// endpoint.Endpoint) that gets wrapped with various middlewares. If you
	// made your own client library, you'd do this work there, so your server
	// could rely on a consistent set of client behavior.
	var registerByEmailEndpoint endpoint.Endpoint
	{
		var (
			name = service.RegisterByEmailServiceName
			rl = apis[name].RateLimit
			bkr = apis[name].Breaker
		)
		registerByEmailEndpoint = grpctransport.NewClient(
			conn,
			pbName,
			name,
			grpcencode.RegisterByEmailRequest,
			grpcdecode.AuthenticationResponse,
			pb.AuthenticationResponse{},
			grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)),
		).Endpoint()

		registerByEmailEndpoint = opentracing.TraceClient(otTracer, name)(registerByEmailEndpoint)

		// We construct a single ratelimiter middleware, to limit the total outgoing
		// QPS from this client to all methods on the remote instance. We also
		// construct per-endpoint circuitbreaker middlewares to demonstrate how
		// that's done, although they could easily be combined into a single breaker
		// for the entire remote instance, too.
		registerByEmailEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(
			rate.Every(rl.Duration),
			rl.Delta))(registerByEmailEndpoint)


		registerByEmailEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(
			bkr.Standardize()),
		)(registerByEmailEndpoint)
	}

	var registerByPhoneEndpoint endpoint.Endpoint
	{
		var (
			name = service.RegisterByPhoneServiceName
			rl = apis[name].RateLimit
			bkr = apis[name].Breaker
		)
		registerByPhoneEndpoint = grpctransport.NewClient(
			conn,
			pbName,
			name,
			grpcencode.RegisterByPhoneRequest,
			grpcdecode.AuthenticationResponse,
			pb.AuthenticationResponse{},
			grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)),
		).Endpoint()

		registerByPhoneEndpoint = opentracing.TraceClient(otTracer, name)(registerByPhoneEndpoint)

		registerByPhoneEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(
			rate.Every(rl.Duration),
			rl.Delta))(registerByPhoneEndpoint)

		registerByPhoneEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(
			bkr.Standardize()),
		)(registerByPhoneEndpoint)
	}

	var sendOTPToEmailEndpoint endpoint.Endpoint
	{
		var (
			name = service.SendOTPToEmail
			rl = apis[name].RateLimit
			bkr = apis[name].Breaker
		)
		sendOTPToEmailEndpoint = grpctransport.NewClient(
			conn,
			pbName,
			name,
			grpcencode.SendOTPToEmailRequest,
			grpcdecode.SendOTPResponse,
			pb.SendOTPResponse{},
			grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)),
		).Endpoint()

		sendOTPToEmailEndpoint = opentracing.TraceClient(otTracer, name)(sendOTPToEmailEndpoint)

		sendOTPToEmailEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(
			rate.Every(rl.Duration),
			rl.Delta))(sendOTPToEmailEndpoint)

		sendOTPToEmailEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(
			bkr.Standardize()),
		)(sendOTPToEmailEndpoint)
	}

	var sendOTPToPhoneEndpoint endpoint.Endpoint
	{
		var (
			name = service.SendOTPToPhone
			rl = apis[name].RateLimit
			bkr = apis[name].Breaker
		)
		sendOTPToPhoneEndpoint = grpctransport.NewClient(
			conn,
			pbName,
			name,
			grpcencode.SendOTPToPhoneRequest,
			grpcdecode.SendOTPResponse,
			pb.SendOTPResponse{},
			grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)),
		).Endpoint()

		sendOTPToPhoneEndpoint = opentracing.TraceClient(otTracer, name)(sendOTPToPhoneEndpoint)

		sendOTPToPhoneEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(
			rate.Every(rl.Duration),
			rl.Delta))(sendOTPToPhoneEndpoint)

		sendOTPToPhoneEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(
			bkr.Standardize()),
		)(sendOTPToPhoneEndpoint)
	}

	// Returning the endpoint.Endpoints as a service.Service relies on the
	// endpoint.Set implementing the Service methods. That's just a simple bit
	// of glue code.
	return serviceendpoints.Set{
		RegisterByEmailEndpoint: registerByEmailEndpoint,
		RegisterByPhoneEndpoint: registerByPhoneEndpoint,
		SendOTPToEmailEndpoint: sendOTPToEmailEndpoint,
		SendOTPToPhoneEndpoint: sendOTPToPhoneEndpoint,
	}
}

func (g GRPCServer) SendOTPToEmail(ctx context.Context, request *pb.SendOTPToEmailRequest) (*pb.SendOTPResponse, error) {
	_, rep, err := g.sendOTPToEmail.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.SendOTPResponse), nil
}

func (g GRPCServer) RegisterByEmail(ctx context.Context, request *pb.RegisterByEmailRequest) (*pb.AuthenticationResponse, error) {
	_, rep, err := g.registerByEmail.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.AuthenticationResponse), nil
}

func (g GRPCServer) SendOTPToPhone(ctx context.Context, request *pb.SendOTPToPhoneRequest) (*pb.SendOTPResponse, error) {
	_, rep, err := g.sendOTPToPhone.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.SendOTPResponse), nil
}

func (g GRPCServer) RegisterByPhone(ctx context.Context, request *pb.RegisterByPhoneRequest) (*pb.AuthenticationResponse, error) {
	_, rep, err := g.registerByPhone.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.AuthenticationResponse), nil
}

