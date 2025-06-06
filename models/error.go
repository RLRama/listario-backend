package models

import "fmt"

type APIError struct {
	StatusCode    int    `json:"status_code"`
	PublicMessage string `json:"public_message"`
	InternalError error  `json:"-"`
}

func (e *APIError) Error() string {
	if e.InternalError != nil {
		return fmt.Sprintf("status %d: %s - internal %v", e.StatusCode, e.PublicMessage, e.InternalError)
	}
	return fmt.Sprintf("status %d: %s", e.StatusCode, e.PublicMessage)
}
