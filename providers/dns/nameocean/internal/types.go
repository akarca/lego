package internal

import "fmt"

// APIError is the error body of the NameOcean API.
type APIError struct {
	Message string `json:"error"`
	Code    string `json:"code"`
	Field   string `json:"field,omitempty"`
}

func (a *APIError) Error() string {
	if a.Field != "" {
		return fmt.Sprintf("%s: %s (%s)", a.Code, a.Message, a.Field)
	}

	return fmt.Sprintf("%s: %s", a.Code, a.Message)
}

// Record is a DNS record.
type Record struct {
	ID       int    `json:"id,omitempty"`
	Type     string `json:"type"`
	Hostname string `json:"hostname"`
	Text     string `json:"text,omitempty"`
	TTL      string `json:"ttl,omitempty"`
}
