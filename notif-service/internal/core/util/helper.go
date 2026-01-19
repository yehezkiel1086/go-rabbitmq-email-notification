package util

import (
	"encoding/json"
	"fmt"
)

func Serialize(data any) ([]byte, error) {
	return json.Marshal(data)
}

func Deserialize(data []byte, v any) error {
	return json.Unmarshal(data, &v)
}

func GenerateConfirmationURL(token string) string {
	return fmt.Sprintf("http://127.0.0.1:8080/api/v1/confirm-email?token=%s", token)
}
