// ginman
// For the full copyright and license information, please view the LICENSE.txt file.

package ginman

import (
	"encoding/json"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

var (
	// Ref: https://github.com/go-playground/validator/pull/1024
	base64RawURLRegexString = "^(?:[A-Za-z0-9-_]{4})*(?:[A-Za-z0-9-_]{2,4})$"
	base64RawURLRegex       = regexp.MustCompile(base64RawURLRegexString)
)

// validateDuration validates a duration field.
var validateDuration validator.Func = func(fl validator.FieldLevel) bool {
	_, err := time.ParseDuration(fl.Field().String())
	return err == nil
}

// validateJSON validates a JSON field.
var validateJSON validator.Func = func(fl validator.FieldLevel) bool {
	switch v := fl.Field().Interface().(type) {
	case string:
		return json.Unmarshal([]byte(v), &json.RawMessage{}) == nil
	default:
		_, err := json.Marshal(v)
		return err == nil
	}
}

// validateBase64Any validates any base64 (base64 or base64url) field.
var validateBase64Any validator.Func = func(fl validator.FieldLevel) bool {
	if validate.Var(fl.Field().Interface(), "base64") == nil || validate.Var(fl.Field().Interface(), "base64url") == nil {
		return true
	} else if base64RawURLRegex.MatchString(fl.Field().String()) {
		// Ref: https://github.com/go-playground/validator/pull/1024
		return true
	}
	return false
}
