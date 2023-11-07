package repository

type URLKeeper interface {
	FindByID(id string) (string, bool)
	FindByLink(link string) string
	SetShortURL(encodedString string, url string) string
	GetNumberOfEntries() int
	PingDBConnect() error
}
