package service

import (
	"context"
	"github.com/al8n/kit-auth/models"
	"github.com/al8n/kit-auth/utils"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	stdopentracing "github.com/opentracing/opentracing-go"
)

type Middleware func(Service) Service


type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{logger: logger, next: next}
	}
}

func (mw *loggingMiddleware) LoginByEmail(ctx context.Context, email, password string) (token string, userInfo *models.UserInfo, err error) {
	defer func() {
		mw.logger.Log("method", "LoginByEmail", "email", email, "err", err)
	}()
	return mw.next.LoginByEmail(ctx, email, password)
}


type instrumentingMiddleware struct {
	ctrs map[string]metrics.Counter
	next  Service
}

func (mw instrumentingMiddleware) LoginByEmail(ctx context.Context, email, password string) (token string, userInfo *models.UserInfo, err error) {
	token, userInfo, err = mw.next.LoginByEmail(ctx, email, password)
	mw.ctrs[LoginByEmailServiceName].Add(1)
	return
}

func InstrumentingMiddleware(ctrs map[string]metrics.Counter) Middleware  {
	return func(next Service) Service {
		return instrumentingMiddleware{
			ctrs: ctrs,
			next:  next,
		}
	}
}

type tracingMiddleware struct {
	tracer stdopentracing.Tracer
	next  Service
}

func (mw tracingMiddleware) LoginByEmail(ctx context.Context, email, password string) (token string, userInfo *models.UserInfo, err error) {
	var (
		span stdopentracing.Span
		spanCtx context.Context
	)

	span, spanCtx = stdopentracing.StartSpanFromContext(ctx, "Login By Email")
	defer span.Finish()

	token, userInfo, err = mw.next.LoginByEmail(spanCtx, email, password)
	if err != nil {
		utils.SetTracerSpanError(span, err)
	} else {
		span.SetTag("User ID", userInfo.ID)
	}

	return
}

func TracingMiddleware(tracer stdopentracing.Tracer) Middleware  {
	return func(next Service) Service {
		return tracingMiddleware{
			tracer: tracer,
			next:  next,
		}
	}
}



