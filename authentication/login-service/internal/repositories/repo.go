package repositories

import (
	"context"
	"github.com/al8n/kit-auth/authentication/common"
	"github.com/al8n/kit-auth/authentication/login-service/config"
	"github.com/al8n/kit-auth/authentication/signature"
	"github.com/al8n/kit-auth/models"
	"github.com/al8n/kit-auth/utils"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

const (
	component = "MongoDB"
)

type Repo struct {
	MongoDB *mongo.Client
}

func NewRepo() (repo *Repo, err error ) {

	var (
		client *mongo.Client
		opt *options.ClientOptions
		cfg = config.GetConfig()
	)

	opt, err = cfg.Mongo.Standardize()
	if err != nil {
		return nil, err
	}

	if client, err = mongo.Connect(context.TODO(), opt); err != nil {
		return
	}

	return &Repo{
		MongoDB: client,
	}, nil
}

func (repo Repo) LoginByEmail (ctx context.Context, email, password string) (token string, userInfo *models.UserInfo, err error)  {
	var (
		collection *mongo.Collection
		user models.User
		span stdopentracing.Span
		spanCtx context.Context
		cfg = config.GetConfig()
	)

	span, spanCtx = stdopentracing.StartSpanFromContext(ctx, "mongo.findOne", stdopentracing.Tag{
		Key: string(ext.Component),
		Value: component,
	})

	defer span.Finish()

	collection = repo.MongoDB.Database(cfg.Mongo.DB).Collection(cfg.Mongo.DB)

	err = collection.FindOne(spanCtx, bson.D{{Key: "email", Value: email}}).Decode(&user)
	if err != nil {
		utils.SetAndLogTracerSpanError(span, err)
		return "", nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(password))

	if err != nil {
		utils.SetAndLogTracerSpanError(span, common.ErrorEmailIdentity)
		return "", nil, common.ErrorEmailIdentity
	}

	userInfo = user.ToInfo()
	token, err = signature.Sign(spanCtx, userInfo.ID, cfg.Service.Secret, config.MachineID, cfg.Service.ExpAt, cfg.Service.Method)
	if err != nil {
		utils.SetAndLogTracerSpanError(span, err)
	}
	return
}
