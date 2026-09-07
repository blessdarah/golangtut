package lib

import (
	"encoding/json"
	"strings"
)

type HttpValidationError struct {
	fields map[string]string
}

func (e *HttpValidationError) Fields() map[string]string {
	if e == nil {
		return nil
	}

	return e.fields
}

func (e *HttpValidationError) Error() string {
	if e == nil {
		return ""
	}

	res, _ := json.Marshal(map[string]map[string]string{
		"errors": e.fields,
	})
	return string(res)
}

func FormatError(err string) *HttpValidationError {
	if err == "" {
		return nil
	}

	errMap := make(map[string]string)

	err = strings.TrimSuffix(err, ".")

	for e := range strings.SplitSeq(err, "; ") {
		key, value, _ := strings.Cut(e, ": ")
		errMap[key] = key + " " + value
	}

	return &HttpValidationError{fields: errMap}
}
