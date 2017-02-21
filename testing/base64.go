package testing

import (
	"encoding/base64"
)

// ToBase64 encodes value into base 64
func ToBase64(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}
