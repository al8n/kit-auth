package config

import (
	"bytes"
	bootapi "github.com/al8n/micro-boot/api"
	bootflag "github.com/al8n/micro-boot/flag"
	"io/ioutil"
	"strings"
	"time"

	"github.com/al8n/kit-auth/authentication/common"
	"github.com/dgrijalva/jwt-go"
)

type Register struct {
	Name 		string 				  `json:"name" yaml:"name"`
	Secret     string                 `json:"secret" yaml:"secret"`
	SecretFile string                 `json:"secret-file" yaml:"secret-file"`
	RawMethod  string                 `json:"method" yaml:"method"`
	Method     jwt.SigningMethod      `json:"-" yaml:"-"`
	ExpAt      time.Duration         `json:"exp-at" yaml:"exp-at"`
	APIs       bootapi.APIs `json:"apis" yaml:"apis"`
}

func (a *Register) BindFlags(fs *bootflag.FlagSet)  {
	fs.StringVar(&a.Name, "auth-name", "Register Service", "the service name")
	fs.StringVar(&a.Secret, "auth-secret", "", "secret for authentication token")
	fs.StringVar(&a.SecretFile, "auth-secret-file", "", "specify a secret file for authentication token. (priority lower than auth-secret)")
	fs.StringVar(&a.RawMethod, "auth-secret-method", "", "signature method for generate token")
	fs.DurationVar(&a.ExpAt, "auth-expiration", 24 * 30 * time.Hour, "the period of validity of token")
}

func (a *Register) Parse() (err error) {
	err = a.parseSecret()
	if err != nil {
		return err
	}
	err = a.APIs.Parse()
	if err != nil {
		return err
	}
	return nil
}

func (a *Register) parseSecret() (err error) {
	secretLength := len(a.Secret)
	if secretLength == 0 {
		if len(a.SecretFile) == 0 {
			return common.ErrorNoAuthenticationSecret
		}
		var file []byte
		file, err = ioutil.ReadFile(a.SecretFile)
		if err != nil {
			return err
		}
		file = bytes.TrimSpace(file)
		a.Secret = string(file)
	}

	a.Secret = strings.TrimSpace(a.Secret)
	secretLength = len(a.Secret)
	switch a.RawMethod {
	case "512", "HS512", "hs512", "HMAC-SSA-512", "hmac-ssa-512":
		if secretLength < HS512 {
			return common.ErrorMismatchAuthenticationSecretAndMethod
		}
		a.Method = jwt.SigningMethodHS512
		return

	case "384", "HS384", "hs384", "HMAC-SSA-384", "hmac-ssa-384":
		if secretLength < HS384 {
			return common.ErrorMismatchAuthenticationSecretAndMethod
		}
		a.Method = jwt.SigningMethodHS384
		return

	case "256", "HS256", "hs256", "HMAC-SSA-256", "hmac-ssa-256":
		if secretLength < HS256 {
			return common.ErrorMismatchAuthenticationSecretAndMethod
		}
		a.Method = jwt.SigningMethodHS256
		return

	default:
		if secretLength < HS256 {
			return common.ErrorMismatchAuthenticationSecretAndMethod
		} else if secretLength >= HS256 && secretLength < HS384 {
			a.Method = jwt.SigningMethodHS256
		} else if secretLength >= HS384 && secretLength < HS512 {
			a.Method = jwt.SigningMethodHS384
		} else {
			a.Method = jwt.SigningMethodHS512
		}
		return nil
	}
}
