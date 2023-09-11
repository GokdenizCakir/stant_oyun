package utils

import "regexp"

func IsPhoneNumber(phone string) bool {
	match, _ := regexp.MatchString(`^(\+\d{1,2}\s?)?\(?\d{3}\)?[\s.-]?\d{3}[\s.-]?\d{4}$`, phone)
	return match
}
