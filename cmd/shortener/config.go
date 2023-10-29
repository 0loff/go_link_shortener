package main

import (
	"flag"
	"os"
)

var config Config

type Config struct {
	Host          string
	ShortLinkHost string
	LogLevel      string
	StorageFile   string
}

type ConfigBuilder struct {
	config Config
}

func (cb ConfigBuilder) SetHost(host string) ConfigBuilder {
	cb.config.Host = host
	return cb
}

func (cb ConfigBuilder) SetShortLinkHost(shortLinkHost string) ConfigBuilder {
	cb.config.ShortLinkHost = shortLinkHost
	return cb
}

func (cb ConfigBuilder) SetLogLevel(logLevel string) ConfigBuilder {
	cb.config.LogLevel = logLevel
	return cb
}

func (cb ConfigBuilder) SetStorageFile(storageFile string) ConfigBuilder {
	cb.config.StorageFile = storageFile
	return cb
}

func (cb ConfigBuilder) Build() Config {
	return cb.config
}

func NewConfigBuilder() {
	var host string
	flag.StringVar(&host, "a", "localhost:8080", "server host")

	var shortLinkHost string
	flag.StringVar(&shortLinkHost, "b", "http://localhost:8080", "host for short link")

	var logLevel string
	flag.StringVar(&logLevel, "l", "info", "log level")

	var storageFile string
	flag.StringVar(&storageFile, "f", "/tmp/short-url-db.json", "storage file full name")

	flag.Parse()

	if envHost := os.Getenv("SERVER_ADDRES"); envHost != "" {
		host = envHost
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		shortLinkHost = envBaseURL
	}

	if envLoglevel := os.Getenv("LOG_LEVEL"); envLoglevel != "" {
		logLevel = envLoglevel
	}

	if envStorageFile := os.Getenv("FILE_STORAGE_PATH"); envStorageFile != "" {
		storageFile = envStorageFile
	}

	config = new(ConfigBuilder).
		SetHost(host).
		SetShortLinkHost(shortLinkHost).
		SetLogLevel(logLevel).
		SetStorageFile(storageFile).
		Build()
}
