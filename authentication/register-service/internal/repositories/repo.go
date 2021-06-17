package repositories

import (
	"context"
	"errors"
	"github.com/al8n/kit-auth/authentication/common"
	"github.com/al8n/kit-auth/authentication/register-service/config"
	"github.com/al8n/kit-auth/authentication/signature"
	"github.com/al8n/kit-auth/models"
	"github.com/al8n/kit-auth/utils"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	component = "MongoDB"
	duplicatedKeyCode = 11000
)

var (
	ErrorPwdOrEmail = errors.New("email or password is incorrect")
)

type Repo struct {
	MongoDB *mongo.Client
	Redis *redis.Client
}

func NewRepo() (repo *Repo, err error ) {

	var (
		client *mongo.Client
		opt *options.ClientOptions
	)

	opt, err = config.GetConfig().Mongo.Standardize()
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

func (repo Repo) SendOTPToEmail(ctx context.Context, email string) (otp string, err error) {
	return
}

func (repo Repo) RegisterByEmail (ctx context.Context, email, _ string) (token string, userInfo *models.UserInfo, err error) {
	var (
		rst *mongo.InsertOneResult
		collection *mongo.Collection
		user *models.User
		//hashPwd []byte
		span stdopentracing.Span
		spanCtx context.Context
		cfg = config.GetConfig()
	)

	span, spanCtx = stdopentracing.StartSpanFromContext(ctx, "mongo.insertOne", stdopentracing.Tag{
		Key: string(ext.Component),
		Value: component,
	})

	defer span.Finish()

	collection = repo.MongoDB.Database(cfg.Mongo.DB).Collection(cfg.Mongo.Collection)

	//hashPwd, err = bcrypt.GenerateFromPassword([]byte(password), salt)
	//if err != nil {
	//	utils.SetAndLogTracerSpanError(span, err)
	//	return "", nil, err
	//}

	now := time.Now().Unix()
	user = &models.User{
		//Username: username,
		Email: email,
		//Password: string(hashPwd),
		//Phone: uuid.New().String(),
		CreatedAt: now,
		UpdatedAt: now,
	}

	rst, err = collection.InsertOne(spanCtx, user)
	if err != nil {
		utils.SetTracerSpanError(span)
		if wE, ok := err.(mongo.WriteException); ok {
			var errs []error

			for _, wErr := range wE.WriteErrors {
				if wErr.Error() != "" {
					span.SetTag("error code", wErr.Code)
				}

				if wErr.Code == duplicatedKeyCode {
					errs = append(errs, common.ErrorDuplicateEmail)
				} else {
					errs = append(errs, wErr)
				}
			}
			utils.LogTracerError(span, errs...)
		} else {
			utils.LogTracerError(span, err)
		}

		return "", nil, err
	}

	user.ID = rst.InsertedID.(primitive.ObjectID)
	userInfo = user.ToInfo()

	token, err = signature.Sign(spanCtx, userInfo.ID, cfg.Service.Secret, config.MachineID, cfg.Service.ExpAt, cfg.Service.Method)

	if err != nil {
		utils.SetAndLogTracerSpanError(span, err)
		return "", nil, err
	}
	return
}


func (repo Repo) SendOTPToPhone(ctx context.Context, prefix, phone string) (otp string, err error) {
	return
}

func (repo Repo) RegisterByPhone (ctx context.Context, prefix, phone, _ string) (token string, userInfo *models.UserInfo, err error) {
	var (
		rst *mongo.InsertOneResult
		collection *mongo.Collection
		user *models.User
		//hashPwd []byte
		span stdopentracing.Span
		spanCtx context.Context
		cfg = config.GetConfig()
	)

	span, spanCtx = stdopentracing.StartSpanFromContext(ctx, "mongo.insertOne", stdopentracing.Tag{
		Key: string(ext.Component),
		Value: component,
	})

	defer span.Finish()

	collection = repo.MongoDB.Database(cfg.Mongo.DB).Collection(cfg.Mongo.Collection)

	now := time.Now().Unix()
	user = &models.User{
		Phone: phone,
		CreatedAt: now,
		UpdatedAt: now,
	}

	rst, err = collection.InsertOne(spanCtx, user)
	if err != nil {
		utils.SetTracerSpanError(span)
		if wE, ok := err.(mongo.WriteException); ok {
			var errs []error

			for _, wErr := range wE.WriteErrors {
				if wErr.Error() != "" {
					span.SetTag("error code", wErr.Code)
				}

				if wErr.Code == duplicatedKeyCode {
					errs = append(errs, common.ErrorDuplicateEmail)
				} else {
					errs = append(errs, wErr)
				}
			}
			utils.LogTracerError(span, errs...)
		} else {
			utils.LogTracerError(span, err)
		}

		return "", nil, err
	}

	user.ID = rst.InsertedID.(primitive.ObjectID)
	userInfo = user.ToInfo()

	token, err = signature.Sign(spanCtx, userInfo.ID, cfg.Service.Secret, config.MachineID, cfg.Service.ExpAt, cfg.Service.Method)

	if err != nil {
		utils.SetAndLogTracerSpanError(span, err)
		return "", nil, err
	}
	return
}

func (repo Repo) SendOTP(email string) (otp string, err error) {
	return
}
