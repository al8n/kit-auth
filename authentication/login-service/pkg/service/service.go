package service

import (
	"context"
	"github.com/al8n/kit-auth/authentication/common"
	"github.com/al8n/kit-auth/authentication/login-service/internal/repositories"
	"github.com/al8n/kit-auth/models"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	stdopentracing "github.com/opentracing/opentracing-go"
)

const (
	LoginByEmailServiceName = "LoginByEmail"
)

type Service interface {
	LoginByEmail(ctx context.Context, email, password string) (token string, userInfo *models.UserInfo, err error)
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

func (svc basicService) LoginByEmail(ctx context.Context, email, password string) (token string, userInfo *models.UserInfo, err error)  {
	if common.EmailValidator(email) {
		return "", nil, common.ErrorInvalidEmail
	}

	if common.PasswordValidator(password) {
		return "",nil, common.ErrorEmailIdentity
	}

	token, userInfo, err = svc.repo.LoginByEmail(ctx,email, password)
	if err != nil {
		return "", nil, err
	}

	return
}

