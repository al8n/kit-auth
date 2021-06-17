module github.com/al8n/kit-auth/gateway

go 1.14

replace github.com/al8n/kit-auth/authentication/service => ../authentication/service

require (
	github.com/al8n/micro-boot v0.0.0-20210617075526-1fbbdc53c9b2
	github.com/go-kit/kit v0.10.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/openzipkin/zipkin-go v0.2.5
	google.golang.org/grpc v1.38.0
)
