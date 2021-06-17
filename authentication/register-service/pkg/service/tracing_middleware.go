package service

import (
	"context"
	"github.com/al8n/kit-auth/models"
	stdopentracing "github.com/opentracing/opentracing-go"
)

type tracingMiddleware struct {
	tracer stdopentracing.Tracer
	next  Service
}

func (mw tracingMiddleware) SendOTPToEmail(ctx context.Context, email string) (otp string, err error) {
	var (
		span stdopentracing.Span
		spanCtx context.Context
	)

	span, spanCtx = stdopentracing.StartSpanFromContext(ctx, "Send OTP to Phone")
	defer span.Finish()

	otp, err = mw.next.SendOTPToEmail(spanCtx, email)

	if err == nil {
		span.SetTag("OTP code", otp)
	}

	return
}

func (mw tracingMiddleware) RegisterByEmail(ctx context.Context, email, otp string) (token string, userInfo *models.UserInfo, err error) {
	var (
		span stdopentracing.Span
		spanCtx context.Context
	)

	span, spanCtx = stdopentracing.StartSpanFromContext(ctx, "Register By Email")
	defer span.Finish()

	token, userInfo, err = mw.next.RegisterByEmail(spanCtx, email,  otp)

	if err == nil {
		span.SetTag("User ID", userInfo.ID)
	}

	return
}

func (mw tracingMiddleware) SendOTPToPhone(ctx context.Context, prefix, phone string) (otp string, err error) {
	var (
		span stdopentracing.Span
		spanCtx context.Context
	)

	span, spanCtx = stdopentracing.StartSpanFromContext(ctx, "Send OTP to Phone")
	defer span.Finish()

	otp, err = mw.next.SendOTPToPhone(spanCtx, prefix, phone)

	if err == nil {
		span.SetTag("OTP code", otp)
	}

	return
}

func (mw tracingMiddleware) RegisterByPhone(ctx context.Context, prefix, phone, otp string) (token string, userInfo *models.UserInfo, err error) {
	var (
		span stdopentracing.Span
		spanCtx context.Context
	)

	span, spanCtx = stdopentracing.StartSpanFromContext(ctx, "Register By Phone")
	defer span.Finish()

	token, userInfo, err = mw.next.RegisterByPhone(spanCtx, prefix, phone,  otp)

	if err == nil {
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
