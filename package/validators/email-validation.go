package validators

import "regexp"

func IsEmailValid(email string) bool {
	var regexEmail = regexp.MustCompile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
	if len(email) < 3 || len(email) > 254 || !regexEmail.MatchString(email) {
		return false
	}
	return true
}
