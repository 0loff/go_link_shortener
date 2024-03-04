package models

// Структура body запроса при создании сокращенного URL
type CreateURLRequestPayload struct {
	URL string
}

// Структура body ответа при запросе на создание сокращенного URL в формате json
type CreateURLResponsePayload struct {
	Result string `json:"result"`
}

// Структура единичной записи в объекте json при запросе на создание нескольких заисей в объекте body запроса
type BatchURLRequestEntry struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// Структура единичной записи в объекте ответа при запросе на создание множества URL
type BatchURLResponseEntry struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// Структура для сохранения записи о сокращенном URL в файле или слайсе
type URLEntry struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	IsDeleted   bool   `json:"-"`
}

// Структура единичной удаленной записи в body ответа на запрос о множественном удалении сокращенных урлов
type DelURLEntry struct {
	UserID   string `db:"user_id"`
	ShortURL string `db:"short_url"`
}

// Stats - Using servise statistic structure
type Metrics struct {
	Urls  int `json:"urls"`
	Users int `json:"users"`
}
