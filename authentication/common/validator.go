package common

import "regexp"

const (
	MaxEmailLength = 320)

var (
	emailPattern = regexp.MustCompile(`^\w+(?:\.\w+){0,3}@(?:[a-zA-Z0-9]+\.){1,3}[a-z]{2,5}$`)
	upperPattern = regexp.MustCompile(`[A-Z]+`)
	lowerPattern = regexp.MustCompile(`[a-z]+`)
	digitPattern = regexp.MustCompile(`[0-9]+`)
	symbolPattern = regexp.MustCompile(`[!@#$%^&*()_+=\-:;,.|><'"?}{\[\]\\/]+`)
)

func EmailValidator(email string) (match bool) {
	if len(email) > MaxEmailLength {
		return false
	}
	return emailPattern.Match([]byte(email))
}

func PasswordValidator(password string) (match bool) {
	pwdLength := len(password)
	if pwdLength < 8 || pwdLength > 16  {
		return false
	}

	var (
		ruleMatch = 0
		bp = []byte(password)
	)

	if upperPattern.Match(bp) {
		ruleMatch++
	}

	if lowerPattern.Match(bp) {
		ruleMatch++
	}

	if digitPattern.Match(bp) {
		ruleMatch++
	}

	if symbolPattern.Match(bp) {
		ruleMatch++
	}

	if ruleMatch < 3 {
		return false
	}

	return true
}
