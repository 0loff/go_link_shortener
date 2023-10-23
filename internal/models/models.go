package models

type CreateURLRequestPayload struct {
	URL string
}

type CreateURLResponsePayload struct {
	Result string `json:"result"`
}
