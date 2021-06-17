package server

import (
	"context"
	authenticationService "github.com/al8n/kit-auth/authentication/service"
	authenticationtransport "github.com/al8n/kit-auth/authentication/"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"google.golang.org/grpc"
	"io"
	"net/url"
	"strings"
)

func AuthenticationFactory(makeEndpoint func(authenticationService.AuthenticationService) endpoint.Endpoint, tracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {

		conn, err := grpc.Dial(instance, grpc.WithInsecure())
		if err != nil {
			return nil, nil, err
		}
		service :=


	}
}
