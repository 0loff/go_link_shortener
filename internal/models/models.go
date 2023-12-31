package models

type CreateURLRequestPayload struct {
	URL string
}

type CreateURLResponsePayload struct {
	Result string `json:"result"`
}

type BatchURLRequestEntry struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchURLResponseEntry struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type BatchInsertURLEntry struct {
	ShortURL    string `db:"short_url"`
	OriginalURL string `db:"origin_url"`
}
