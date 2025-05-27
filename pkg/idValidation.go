package pkg

import (
	goaway "github.com/TwiN/go-away"
	"regexp"
)

func IsValidID(id string) (bool, string) {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]{3,20}$`, id)
	if !matched {
		return false, "ID must be 3–20 characters and only contain letters, numbers, - or _"
	}

	if !goaway.IsProfane(id) {
		return false, "Please use appropriate language."
	}
	return true, ""
}
