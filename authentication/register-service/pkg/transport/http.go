package transport

import (
	bootapi "github.com/al8n/micro-boot/api"
	"github.com/al8n/kit-auth/authentication/register-service/internal/codec/httpcodec"
	"github.com/al8n/kit-auth/authentication/register-service/internal/codec/httpcodec/httpdecode"
	"github.com/al8n/kit-auth/authentication/register-service/internal/codec/httpcodec/httpencode"
	serviceendpoint "github.com/al8n/kit-auth/authentication/register-service/pkg/endpoint"
	"github.com/al8n/kit-auth/authentication/register-service/pkg/service"
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
		r = mux.NewRouter()
	)
	{
		ste := apis[service.SendOTPToEmail]
		r.Methods(ste.Method).Path(ste.Path).Handler(httptransport.NewServer(
			endpoints.SendOTPToEmailEndpoint,
			httpdecode.SendOTPToEmailRequest,
			httpencode.SendOTPResponse,
			append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "SendOTPToEmail", logger)))...,
		))
	}
	{
		rbe := apis[service.RegisterByEmailServiceName]

		r.Methods(rbe.Method).Path(rbe.Path).Handler(httptransport.NewServer(
			endpoints.RegisterByEmailEndpoint,
			httpdecode.RegisterByEmailRequest,
			httpencode.AuthenticationResponse,
			append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "RegisterByEmail", logger)))...,
		))
	}
	{
		stp := apis[service.SendOTPToPhone]
		r.Methods(stp.Method).Path(stp.Path).Handler(httptransport.NewServer(
			endpoints.SendOTPToPhoneEndpoint,
			httpdecode.SendOTPToPhoneRequest,
			httpencode.SendOTPResponse,
			append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "SendOTPToPhone", logger)))...,
		))
	}
	{
		rbp := apis[service.RegisterByPhoneServiceName]

		r.Methods(rbp.Method).Path(rbp.Path).Handler(httptransport.NewServer(
			endpoints.RegisterByEmailEndpoint,
			httpdecode.RegisterByPhoneRequest,
			httpencode.AuthenticationResponse,
			append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "RegisterByPhone", logger)))...,
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

	// Each individual endpoint is an http/transport.Client (which implements
	// endpoint.Endpoint) that gets wrapped with various middlewares. If you
	// made your own client library, you'd do this work there, so your server
	// could rely on a consistent set of client behavior.
	var registerByEmailEndpoint endpoint.Endpoint
	{
		var (
			name = service.RegisterByEmailServiceName
			rbe = apis[name]
		)

		registerByEmailEndpoint = httptransport.NewClient(
			rbe.Method,
			copyURL(u, rbe.Path),
			httpencode.GenericRequest,
			httpdecode.AuthenticationResponse,
			httptransport.ClientBefore(opentracing.ContextToHTTP(otTracer, logger)),
		).Endpoint()

		registerByEmailEndpoint = opentracing.TraceClient(otTracer, name)(registerByEmailEndpoint)

		// We construct a single ratelimiter middleware, to limit the total outgoing
		// QPS from this client to all methods on the remote instance. We also
		// construct per-endpoint circuitbreaker middlewares to demonstrate how
		// that's done, although they could easily be combined into a single breaker
		// for the entire remote instance, too.
		registerByEmailEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(
			rate.Every(rbe.RateLimit.Duration),
			rbe.RateLimit.Delta),
		)(registerByEmailEndpoint)

		registerByEmailEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(
			rbe.Breaker.Standardize()),
		)(registerByEmailEndpoint)
	}

	var registerByPhoneEndpoint endpoint.Endpoint
	{
		var (
			name = service.RegisterByPhoneServiceName
			rbp = apis[name]
		)

		registerByPhoneEndpoint = httptransport.NewClient(
			rbp.Method,
			copyURL(u, rbp.Path),
			httpencode.GenericRequest,
			httpdecode.AuthenticationResponse,
			httptransport.ClientBefore(opentracing.ContextToHTTP(otTracer, logger)),
		).Endpoint()

		registerByPhoneEndpoint = opentracing.TraceClient(otTracer, name)(registerByPhoneEndpoint)

		registerByPhoneEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(
			rate.Every(rbp.RateLimit.Duration),
			rbp.RateLimit.Delta),
		)(registerByPhoneEndpoint)

		registerByPhoneEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(
			rbp.Breaker.Standardize()),
		)(registerByPhoneEndpoint)
	}

	var sendOTPToPhoneEndpoint endpoint.Endpoint
	{
		var (
			name = service.SendOTPToPhone
			stp = apis[name]
		)

		sendOTPToPhoneEndpoint = httptransport.NewClient(
			stp.Method,
			copyURL(u, stp.Path),
			httpencode.GenericRequest,
			httpdecode.SendOTPResponse,
			httptransport.ClientBefore(opentracing.ContextToHTTP(otTracer, logger)),
		).Endpoint()

		sendOTPToPhoneEndpoint = opentracing.TraceClient(otTracer, name)(sendOTPToPhoneEndpoint)

		sendOTPToPhoneEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(
			rate.Every(stp.RateLimit.Duration),
			stp.RateLimit.Delta),
		)(sendOTPToPhoneEndpoint)

		sendOTPToPhoneEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(
			stp.Breaker.Standardize()),
		)(sendOTPToPhoneEndpoint)
	}

	var sendOTPToEmailEndpoint endpoint.Endpoint
	{
		var (
			name = service.SendOTPToEmail
			stp = apis[name]
		)

		sendOTPToEmailEndpoint = httptransport.NewClient(
			stp.Method,
			copyURL(u, stp.Path),
			httpencode.GenericRequest,
			httpdecode.SendOTPResponse,
			httptransport.ClientBefore(opentracing.ContextToHTTP(otTracer, logger)),
		).Endpoint()

		sendOTPToEmailEndpoint = opentracing.TraceClient(otTracer, name)(sendOTPToEmailEndpoint)

		sendOTPToEmailEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(
			rate.Every(stp.RateLimit.Duration),
			stp.RateLimit.Delta),
		)(sendOTPToEmailEndpoint)

		sendOTPToEmailEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(
			stp.Breaker.Standardize()),
		)(sendOTPToEmailEndpoint)
	}

	// Returning the endpoint.Set as a service.Service relies on the
	// endpoint.Set implementing the Service methods. That's just a simple bit
	// of glue code.
	return serviceendpoint.Set{
		RegisterByEmailEndpoint: registerByEmailEndpoint,
		RegisterByPhoneEndpoint: registerByPhoneEndpoint,
		SendOTPToEmailEndpoint: sendOTPToEmailEndpoint,
		SendOTPToPhoneEndpoint: sendOTPToPhoneEndpoint,
	}, nil
}

func copyURL(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}
