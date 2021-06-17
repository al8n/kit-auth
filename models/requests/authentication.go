package requests


type RegisterByPhoneRequest struct {
	Prefix string `json:"prefix"`
	Phone string `json:"phone"`
	OTP string `json:"otp"`
}

type RegisterByEmailRequest struct {
	Email string `json:"email"`
	OTP string `json:"otp"`
}

type LoginByEmailRequest struct {
	Email string `json:"email"`
	OTP string `json:"otp"`
}

type LoginByPhoneRequest struct {
	Prefix string `json:"prefix"`
	Phone string `json:"phone"`
	OTP string `json:"otp"`
}

type SendOTPToEmailRequest struct {
	Email string `json:"email"`
}

type SendOTPToPhoneRequest struct {
	Prefix string `json:"prefix"`
	Phone string `json:"phone"`
}