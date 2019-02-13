package core

import (
	"regexp"
)

var nameRegexp = regexp.MustCompile(`^[0-9A-Za-z-_]{1,16}$`)

func checkName(name string) bool {
	return nameRegexp.Match([]byte(name))
}

var emailRegexp = regexp.MustCompile(`^([a-zA-Z0-9_\-\.]+)@([a-zA-Z0-9_\-\.]+)\.([a-zA-Z]{2,5})$`)

func checkEmail(email string) bool {
	return emailRegexp.Match([]byte(email))
}
