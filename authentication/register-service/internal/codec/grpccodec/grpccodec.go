package grpccodec

import (
	"github.com/al8n/kit-auth/models/requests"
	"github.com/al8n/kit-auth/models/responses"
	"github.com/al8n/kit-auth/authentication/register-service/pb"
)

func SendOTPToEmailReq2pbReq(req requests.SendOTPToEmailRequest) (pbReq *pb.SendOTPToEmailRequest) {
	pbReq = &pb.SendOTPToEmailRequest{
		Email: req.Email,
	}
	return
}

func SendOTPToPhoneReq2pbReq(req requests.SendOTPToPhoneRequest) (pbReq *pb.SendOTPToPhoneRequest) {
	pbReq = &pb.SendOTPToPhoneRequest{
		Prefix: req.Prefix,
		Phone: req.Phone,
	}
	return
}

func RegisterByPhoneReq2pbReq(req requests.RegisterByPhoneRequest) (pbReq *pb.RegisterByPhoneRequest) {
	pbReq = &pb.RegisterByPhoneRequest{}
	pbReq.Prefix = req.Prefix
	pbReq.Phone = req.Phone
	pbReq.Otp = req.OTP
	return
}

func RegisterByEmailReq2pbReq(req requests.RegisterByEmailRequest) (pbReq *pb.RegisterByEmailRequest) {
	pbReq = &pb.RegisterByEmailRequest{}
	pbReq.Email = req.Email
	pbReq.Otp = req.OTP
	return
}

func SendOTPpbResp2Resp(pbResp pb.SendOTPResponse) (res *responses.SendOTPResponse) {
	res = &responses.SendOTPResponse{
		OTP:   pbResp.Otp,
		Error: pbResp.Error,
	}
	return
}

func SendOTPResp2pbResp(res *responses.SendOTPResponse) (pbResp *pb.SendOTPResponse) {
	pbResp = &pb.SendOTPResponse{
		Otp:   res.OTP,
		Error: res.Error,
	}
	return
}

func AuthpbResp2Resp(pbReply pb.AuthenticationResponse) (res *responses.AuthenticationResponse) {
	res = &responses.AuthenticationResponse{}

	// Ugly value assignment, but faster than using reflect
	res.Token = pbReply.Token
	res.Error = pbReply.Error
	res.User.ID = pbReply.User.Id
	res.User.Email = pbReply.User.Email
	res.User.Username = pbReply.User.Username
	res.User.Phone = pbReply.User.Phone
	res.User.Avatar = pbReply.User.Avatar
	res.User.Description = pbReply.User.Description
	res.User.Age = pbReply.User.Age
	res.User.Gender = pbReply.User.Gender
	res.User.University = pbReply.User.University
	res.User.Major = pbReply.User.Major
	res.User.City = pbReply.User.City
	res.User.Country = pbReply.User.Country
	res.User.Membership = pbReply.User.Membership
	res.User.EnrollAt = pbReply.User.EnrollAt
	res.User.MembershipAt = pbReply.User.MembershipAt
	res.User.MembershipExpiredAt = pbReply.User.MembershipExpiredAt
	res.User.CreatedAt = pbReply.User.CreatedAt
	return
}

func AuthResp2pbResp(res responses.AuthenticationResponse) (pbReply *pb.AuthenticationResponse)  {
	pbReply = &pb.AuthenticationResponse{}

	// Ugly value assignment, but faster than using reflect
	pbReply.Token = res.Token
	pbReply.Error = res.Error
	pbReply.User.Id = res.User.ID
	pbReply.User.Email = res.User.Email
	pbReply.User.Username = res.User.Username
	pbReply.User.Phone = res.User.Phone
	pbReply.User.Avatar = res.User.Avatar
	pbReply.User.Description = res.User.Description
	pbReply.User.Age = res.User.Age
	pbReply.User.Gender = res.User.Gender
	pbReply.User.University = res.User.University
	pbReply.User.Major = res.User.Major
	pbReply.User.City = res.User.City
	pbReply.User.Country = res.User.Country
	pbReply.User.Membership = res.User.Membership
	pbReply.User.EnrollAt = res.User.EnrollAt
	pbReply.User.MembershipAt = res.User.MembershipAt
	pbReply.User.MembershipExpiredAt = res.User.MembershipExpiredAt
	pbReply.User.CreatedAt = res.User.CreatedAt
	return
}