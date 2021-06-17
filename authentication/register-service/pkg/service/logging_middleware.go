package service

import (
	"context"
	"github.com/al8n/kit-auth/models"
	"github.com/go-kit/kit/log"
)

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{logger: logger, next: next}
	}
}

func (mw *loggingMiddleware) SendOTPToEmail(ctx context.Context, email string) (otp string, err error)  {
	defer func() {
		mw.logger.Log("method", "SendOTPToEmail", "email", email, "err", err)
	}()
	return mw.next.SendOTPToEmail(ctx, email)
}

func (mw *loggingMiddleware) RegisterByEmail(ctx context.Context, email, otp string) (token string, userInfo *models.UserInfo, err error)  {
	defer func() {
		mw.logger.Log("method", "RegisterByEmail", "email", email, "err", err)
	}()
	return mw.next.RegisterByEmail(ctx, email, otp)
}

func (mw *loggingMiddleware) SendOTPToPhone(ctx context.Context, prefix, phone string) (otp string, err error)  {
	defer func() {
		mw.logger.Log("method", "SendOTPToPhone", "phone", "+" + prefix + phone, "err", err)
	}()
	return mw.next.SendOTPToPhone(ctx, prefix, phone)
}

func (mw *loggingMiddleware) RegisterByPhone(ctx context.Context, prefix, phone, otp string) (token string, userInfo *models.UserInfo, err error)  {
	defer func() {
		mw.logger.Log("method", "RegisterByPhone", "phone", "+" + prefix + phone, "err", err)
	}()
	return mw.next.RegisterByPhone(ctx, prefix, phone, otp)
}

