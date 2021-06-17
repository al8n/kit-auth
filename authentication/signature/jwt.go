package signature

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	stdjwt "github.com/dgrijalva/jwt-go"
	stdopentracing "github.com/opentracing/opentracing-go"
	"strings"
	"time"
)

const signatureOP = "Signature"

type AuthenticationClaims struct {
	ID string `json:"id"`
	stdjwt.StandardClaims
}

func keyFunc(token *stdjwt.Token, secret string) (interface{}, error)  {
	return []byte(secret), nil
}

func Sign(ctx context.Context, id, secret string, machineID uint64, expiration time.Duration, method stdjwt.SigningMethod) (token string, err error)  {
	var (
		now      time.Time
		jtiRaw   []byte
		jti      string
		claims   AuthenticationClaims
		rawToken *stdjwt.Token
		span     stdopentracing.Span
	)

	span, _ = stdopentracing.StartSpanFromContext(ctx, signatureOP)
	defer span.Finish()

	now = time.Now()

	jtiRaw = make([]byte, 128)
	_, err = rand.Read(jtiRaw)
	jti = replaceSlashesAndPlus(base64.StdEncoding.EncodeToString(jtiRaw))
	if err != nil {
		return "", err
	}

	claims = AuthenticationClaims{
		ID:             id,
		StandardClaims: stdjwt.StandardClaims{
			ExpiresAt: now.Add(expiration).Unix(),
			Issuer: fmt.Sprintf("%d", machineID),
			Id: jti,
			IssuedAt: now.Unix(),
		},
	}


	rawToken = stdjwt.NewWithClaims(method, claims)

	token, err = rawToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	span.SetTag("exp", claims.ExpiresAt)
	span.SetTag("isa", claims.IssuedAt)
	return
}


func replaceSlashesAndPlus(str string) string {
	str = strings.Replace(str, "/", "0", -1)
	str = strings.Replace(str, "+", "0", -1)
	return str
}

