package transport

import (
	"context"
	bootapi "github.com/al8n/micro-boot/api"
	"github.com/al8n/kit-auth/authentication/login-service/internal/codec/grpccodec/grpcdecode"
	"github.com/al8n/kit-auth/authentication/login-service/internal/codec/grpccodec/grpcencode"
	"github.com/al8n/kit-auth/authentication/login-service/pb"
	serviceendpoints "github.com/al8n/kit-auth/authentication/login-service/pkg/endpoint"
	"github.com/al8n/kit-auth/authentication/login-service/pkg/service"
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

const pbName = "pb.Login"

type GRPCServer struct {
	registerByEmail grpctransport.Handler
	loginByEmail grpctransport.Handler
	sendCaptchaToEmail grpctransport.Handler
}

// NewGRPCServer makes a set of endpoints available as a gRPC AuthenticationServer.
func NewGRPCServer(set serviceendpoints.Set, otTracer stdopentracing.Tracer, logger log.Logger, apis bootapi.APIs) pb.LoginServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	return &GRPCServer{
		loginByEmail: grpctransport.NewServer(
			set.LoginByEmailEndpoint,
			grpcdecode.LoginByEmailRequest,
			grpcencode.AuthenticationResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, service.LoginByEmailServiceName, logger)))...,
		),
	}
}

// NewGRPCClient returns an AuthenticationService backed by a gRPC server at the other end
// of the conn. The caller is responsible for constructing the conn, and
// eventually closing the underlying transport. We bake-in certain middlewares,
// implementing the client library pattern.
func NewGRPCClient(conn *grpc.ClientConn, otTracer stdopentracing.Tracer, logger log.Logger, apis bootapi.APIs) service.Service {

	var loginByEmailEndpoint endpoint.Endpoint
	{
		var (
			name = service.LoginByEmailServiceName
			rl = apis[name].RateLimit
			bkr = apis[name].Breaker
		)

		loginByEmailEndpoint = grpctransport.NewClient(
			conn,
			pbName,
			name,
			grpcencode.LoginByEmailRequest,
			grpcdecode.AuthenticationResponse,
			pb.AuthenticationResponse{},
			grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)),
		).Endpoint()

		loginByEmailEndpoint = opentracing.TraceClient(otTracer, name)(loginByEmailEndpoint)

		loginByEmailEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(
			rate.Every(rl.Duration),
			rl.Delta),
		)(loginByEmailEndpoint)

		loginByEmailEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(
			bkr.Standardize()),
		)(loginByEmailEndpoint)
	}

	// Returning the endpoint.Endpoints as a service.Service relies on the
	// endpoint.Set implementing the Service methods. That's just a simple bit
	// of glue code.
	return serviceendpoints.Set{
		LoginByEmailEndpoint: loginByEmailEndpoint,
	}
}

func (g GRPCServer) LoginByEmail(ctx context.Context, request *pb.LoginByEmailRequest) (*pb.AuthenticationResponse, error) {
	_, rep, err := g.loginByEmail.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.AuthenticationResponse), nil
}
