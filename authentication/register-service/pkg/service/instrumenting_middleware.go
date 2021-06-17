package service

import (
	"context"
	"github.com/al8n/kit-auth/models"
	"github.com/go-kit/kit/metrics"

)

type instrumentingMiddleware struct {
	ctrs map[string]metrics.Counter
	next  Service
}

func (mw instrumentingMiddleware) SendOTPToEmail(ctx context.Context, email string) (otp string, err error) {
	otp, err = mw.next.SendOTPToEmail(ctx, email)
	mw.ctrs[SendOTPToEmail].Add(1)
	return
}

func (mw instrumentingMiddleware) RegisterByEmail(ctx context.Context, email, otp string) (token string, userInfo *models.UserInfo, err error) {
	token, userInfo, err = mw.next.RegisterByEmail(ctx, email, otp)
	mw.ctrs[RegisterByEmailServiceName].Add(1)
	return
}

func (mw instrumentingMiddleware) SendOTPToPhone(ctx context.Context, prefix, phone string) (otp string, err error) {
	otp, err = mw.next.SendOTPToPhone(ctx, prefix, phone)
	mw.ctrs[SendOTPToPhone].Add(1)
	return
}

func (mw instrumentingMiddleware) RegisterByPhone(ctx context.Context, prefix, phone, otp string) (token string, userInfo *models.UserInfo, err error) {
	token, userInfo, err = mw.next.RegisterByPhone(ctx, prefix, phone, otp)
	mw.ctrs[RegisterByPhoneServiceName].Add(1)
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





