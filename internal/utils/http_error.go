package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func ReadErrorResponse(resp *http.Response) error {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
	if len(body) == 0 {
		return fmt.Errorf("unexpected status %s", resp.Status)
	}

	var payload struct {
		Reason  string `json:"reason"`
		Msg     string `json:"msg"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(body, &payload); err == nil {
		for _, msg := range []string{payload.Reason, payload.Msg, payload.Message} {
			msg = strings.TrimSpace(msg)
			if msg != "" {
				return fmt.Errorf("unexpected status %s: %s", resp.Status, msg)
			}
		}
	}

	return fmt.Errorf("unexpected status %s: %s", resp.Status, strings.TrimSpace(string(body)))
}
