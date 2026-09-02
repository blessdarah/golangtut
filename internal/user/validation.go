package user

import (
	"errors"
	"regexp"
	"strings"

	"github.com/twharmon/govalid"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func init() {
	govalid.Rule("email", func(v any) error {
		s, ok := v.(string)
		if !ok {
			return errors.New("email rule expects a string value")
		}
		if !emailRegex.MatchString(s) {
			return govalid.NewValidationError("must be a valid email address")
		}
		return nil
	})

	govalid.Rule("alphanum", func(v any) error {
		s, ok := v.(string)
		if !ok {
			return errors.New("alphanum rule expects a string value")
		}
		for _, r := range s {
			if r < '0' || (r > '9' && r < 'A') || (r > 'Z' && r < 'a') || r > 'z' {
				return govalid.NewValidationError("must only contain letters and digits")
			}
		}
		if strings.TrimSpace(s) == "" {
			return govalid.NewValidationError("must only contain letters and digits")
		}
		return nil
	})
}