package transport

import (
	bootapi "github.com/al8n/micro-boot/api"
	"github.com/al8n/kit-auth/authentication/login-service/internal/codec/httpcodec"
	"github.com/al8n/kit-auth/authentication/login-service/internal/codec/httpcodec/httpdecode"
	"github.com/al8n/kit-auth/authentication/login-service/internal/codec/httpcodec/httpencode"
	serviceendpoint "github.com/al8n/kit-auth/authentication/login-service/pkg/endpoint"
	"github.com/al8n/kit-auth/authentication/login-service/pkg/service"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"net/url"
	"strings"
)

// NewHTTPHandler returns an HTTP handler that makes a set of endpoints
// available on predefined paths.
func NewHTTPHandler(endpoints serviceendpoint.Set, otTracer stdopentracing.Tracer, logger log.Logger, apis bootapi.APIs) *mux.Router {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(httpcodec.ErrorEncoder),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	var (
		r *mux.Router
		lbe bootapi.API
	)
	{
		r = mux.NewRouter()

		lbe = apis[service.LoginByEmailServiceName]
		r.Methods(lbe.Method).Path(lbe.Path).Handler(httptransport.NewServer(
			endpoints.LoginByEmailEndpoint,
			httpdecode.LoginByEmailRequest,
			httpencode.AuthenticationResponse,
			append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "LoginByEmail", logger)))...,
		))
	}

	return r
}

// NewHTTPClient returns an Authentication backed by an HTTP server living at the
// remote instance. We expect instance to come from a service discovery system,
// so likely of the form "host:port". We bake-in certain middlewares,
// implementing the client library pattern.
func NewHTTPClient(instance string, otTracer stdopentracing.Tracer, logger log.Logger, apis bootapi.APIs) (service.Service, error) {
	// Quickly sanitize the instance string.
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}

	// The LoginByEmail endpoint is the same thing, with slightly different
	// middlewares to demonstrate how to specialize per-endpoint.
	var loginByEmail endpoint.Endpoint
	{
		var (
			name = service.LoginByEmailServiceName
			lbe = apis[name]
		)

		loginByEmail = httptransport.NewClient(
			lbe.Method,
			copyURL(u, lbe.Path),
			httpencode.GenericRequest,
			httpdecode.AuthenticationResponse,
			httptransport.ClientBefore(opentracing.ContextToHTTP(otTracer, logger)),
		).Endpoint()

		loginByEmail = opentracing.TraceClient(otTracer, lbe.Name)(loginByEmail)


		// We construct a single ratelimiter middleware, to limit the total outgoing
		// QPS from this client to all methods on the remote instance. We also
		// construct per-endpoint circuitbreaker middlewares to demonstrate how
		// that's done, although they could easily be combined into a single breaker
		// for the entire remote instance, too.
		loginByEmail = ratelimit.NewErroringLimiter(rate.NewLimiter(
			rate.Every(lbe.RateLimit.Duration),
			lbe.RateLimit.Delta),
		)(loginByEmail)

		loginByEmail = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(
			lbe.Breaker.Standardize()),
		)(loginByEmail)
	}

	// Returning the endpoint.Set as a service.Service relies on the
	// endpoint.Set implementing the Service methods. That's just a simple bit
	// of glue code.
	return serviceendpoint.Set{
		LoginByEmailEndpoint: loginByEmail,
	}, nil
}

func copyURL(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}
