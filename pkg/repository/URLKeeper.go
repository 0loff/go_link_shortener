package repository

type URLKeeper interface {
	FindByID(id string) string
	FindByLink(link string) string
	SetShortURL(encodedString string, url string)
	GetNumberOfEntries() int
	PingConnect() error
}
