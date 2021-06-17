package responses

import (
	"github.com/al8n/kit-auth/models"
)

type SendOTPResponse struct {
	OTP string 	`json:"otp,omitempty"`
	Error  string	`json:"err,omitempty"`
}

type AuthenticationResponse struct {
	Token string 	`json:"token"`
	User models.UserInfo `json:"user"`
	Error  string `json:"err,omitempty"`
}
