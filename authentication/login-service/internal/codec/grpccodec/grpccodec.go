package grpccodec

import (
	"github.com/al8n/kit-auth/authentication/login-service/pb"
	"github.com/al8n/kit-auth/models/requests"
	"github.com/al8n/kit-auth/models/responses"
)


func LoginByEmailReq2pbReq(req requests.LoginByEmailRequest) (pbReq *pb.LoginByEmailRequest)  {
	pbReq = &pb.LoginByEmailRequest{}
	pbReq.Email = req.Email
	pbReq.Otp = req.OTP
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