package transform

import (
	"encoding/base64"
)

func b64enc(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func b64dec(s string) string {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		// return original if decode fails
		return s
	}
	return string(b)
}
