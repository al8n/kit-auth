package models

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Email               string             `bson:"email,omitempty" json:"email,omitempty"`
	Username            string             `bson:"username, omitempty" json:"username"`
	Phone               string             `bson:"phone,omitempty" json:"phone,omitempty"`

	WeChat              string             `bson:"wechat,omitempty" json:"wechat,omitempty"`
	Weibo               string             `bson:"weibo,omitempty" json:"weibo,omitempty"`
	Facebook            string             `bson:"facebook,omitempty" json:"facebook,omitempty"`
	Twitter             string             `bson:"twitter,omitempty" json:"twitter,omitempty"`

	Password            string             `bson:"password" json:"password"`
	Salt                string             `bson:"salt" json:"sort"`
	Avatar              string             `bson:"avatar" json:"avatar"`
	Description         string             `bson:"description" json:"description"`
	Age                 uint8              `bson:"age" json:"age"`
	Gender              uint8              `bson:"gender" json:"gender"`
	University          string             `bson:"university" json:"university"`
	Major               string             `bson:"major" json:"major"`
	EnrollAt            int64              `bson:"enroll_at" json:"enroll_at"`
	City                string             `bson:"city" json:"city"`
	Country             string             `bson:"country" json:"country"`
	Membership          bool               `bson:"membership" json:"membership"`
	MembershipAt        int64              `bson:"membership_at" json:"membership_at"`
	MembershipExpiredAt int64              `bson:"membership_expired_at" json:"membership_expired_at"`
	CreatedAt           int64              `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt           int64              `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	DeactivatedAt       int64              `bson:"deactivated_at,omitempty" json:"deactivated_at,omitempty"`
}

type UserInfo struct {
	ID                  string `bson:"id" json:"id"`
	Email               string `bson:"email,omitempty" json:"email,omitempty"`
	Username            string `bson:"username" json:"username"`
	Phone               string `bson:"phone" json:"phone"`

	WeChat              string `bson:"wechat,omitempty" json:"wechat,omitempty"`
	Weibo               string `bson:"weibo,omitempty" json:"weibo,omitempty"`
	Facebook            string `bson:"facebook,omitempty" json:"facebook,omitempty"`
	Twitter             string `bson:"twitter,omitempty" json:"twitter,omitempty"`

	Avatar              string `bson:"avatar" json:"avatar"`
	Description         string `bson:"description" json:"description"`
	Age                 uint32 `bson:"age" json:"age"`
	Gender              uint32 `bson:"gender" json:"gender"`
	University          string `bson:"university" json:"university"`
	Major               string `bson:"major" json:"marjor"`
	EnrollAt            int64  `bson:"enroll_at" json:"enroll_at"`
	City                string `bson:"city" json:"city"`
	Country             string `bson:"country" json:"country"`
	Membership          bool   `bson:"membership" json:"membership"`
	MembershipAt        int64  `bson:"membership_at" json:"membership_at"`
	MembershipExpiredAt int64  `bson:"membership_expired_at" json:"membership_expired_at"`
	CreatedAt           int64  `bson:"created_at" json:"created_at"`
}

func (user *User) ToInfo() *UserInfo {
	return UserToInfo(user)
}

func UserToInfo(user *User) *UserInfo {
	var (
		err   error
		phone = ""
		email = ""
	)

	_, err = uuid.Parse(user.Phone)
	if err != nil {
		phone = user.Phone
	}

	_, err = uuid.Parse(user.Email)
	if err != nil {
		email = user.Email
	}

	return &UserInfo{
		ID:                  user.ID.Hex(),
		Email:               email,
		Username:            user.Username,
		Phone:               phone,
		WeChat: 			 user.WeChat,
		Weibo: 			 	 user.Weibo,
		Facebook:            user.Facebook,
		Twitter:             user.Twitter,
		Avatar:              user.Avatar,
		Description:         user.Description,
		Age:                 uint32(user.Age),
		Gender:              uint32(user.Gender),
		University:          user.University,
		Major:               user.Major,
		EnrollAt:            user.EnrollAt,
		City:                user.City,
		Country:             user.Country,
		Membership:          user.Membership,
		MembershipAt:        user.MembershipAt,
		MembershipExpiredAt: user.MembershipExpiredAt,
		CreatedAt:           user.CreatedAt,
	}
}
