package service

import (
	"context"
	"github.com/al8n/kit-auth/authentication/common"
	"github.com/al8n/kit-auth/authentication/register-service/internal/repositories"
	"github.com/al8n/kit-auth/models"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	stdopentracing "github.com/opentracing/opentracing-go"
)

const (
	RegisterByEmailServiceName = "RegisterByEmail"
	RegisterByPhoneServiceName = "RegisterByPhone"
	SendOTPToEmail = "SendOTPToEmail"
	SendOTPToPhone = "SendOTPToPhone"
)

type Service interface {
	RegisterByEmail(ctx context.Context, email, otp string) (token string, userInfo *models.UserInfo, err error)
	RegisterByPhone(ctx context.Context, prefix, phone, otp string) (token string, userInfo *models.UserInfo, err error)

	SendOTPToPhone(ctx context.Context, prefix, phone string) (otp string, err error)
	SendOTPToEmail(ctx context.Context, email string) (otp string, err error)
}

// New returns a basic Service with all of the expected middlewares wired in.
func New(logger log.Logger, counters map[string]metrics.Counter, tracer stdopentracing.Tracer) (svc Service,  err error) {

	svc, err = NewBasicService()
	if err != nil {
		return nil, err
	}

	svc = LoggingMiddleware(logger)(svc)
	svc = InstrumentingMiddleware(counters)(svc)
	svc = TracingMiddleware(tracer)(svc)

	return svc, nil
}

type basicService struct {
	repo *repositories.Repo
}

func NewBasicService() (svc Service, err error ) {
	var repo *repositories.Repo

	repo, err = repositories.NewRepo()
	if err != nil {
		return
	}

	return &basicService{
		repo: repo,
	}, nil
}

func (svc basicService) SendOTPToEmail(ctx context.Context, email string) (otp string, err error) {
	if !common.EmailValidator(email) {
		return "", common.ErrorInvalidEmail
	}

	otp, err = svc.repo.SendOTPToEmail(ctx, email)
	if err != nil {
		return "", err
	}

	return
}

func (svc basicService) RegisterByEmail(ctx context.Context, email, otp string) (token string, userInfo *models.UserInfo, err error) {
	if !common.EmailValidator(email) {
		return "", nil, common.ErrorInvalidEmail
	}

	//if !common.PasswordValidator(password) {
	//	return "", nil, common.ErrorEmailIdentity
	//}

	token, userInfo, err = svc.repo.RegisterByEmail(ctx, email, otp)
	if err != nil {
		return "", nil, err
	}

	return
}

func (svc basicService) SendOTPToPhone(ctx context.Context, prefix, phone string) (otp string, err error) {


	otp, err = svc.repo.SendOTPToPhone(ctx, prefix, phone)
	if err != nil {
		return "", err
	}

	return
}

func (svc basicService) RegisterByPhone(ctx context.Context, prefix, phone, otp string) (token string, userInfo *models.UserInfo, err error) {
	//if !common.EmailValidator(email) {
	//	return "", nil, common.ErrorInvalidEmail
	//}

	//if !common.PasswordValidator(password) {
	//	return "", nil, common.ErrorEmailIdentity
	//}

	token, userInfo, err = svc.repo.RegisterByPhone(ctx, prefix, phone, otp)
	if err != nil {
		return "", nil, err
	}

	return
}