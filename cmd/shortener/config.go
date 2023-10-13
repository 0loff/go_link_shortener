package main

import (
	"flag"
	"os"
)

var config Config

type Config struct {
	Host          string
	ShortLinkHost string
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

func (cb ConfigBuilder) Build() Config {
	return cb.config
}

func NewConfigBuilder() {
	var host string
	flag.StringVar(&host, "a", "localhost:8080", "server host")

	var shortLinkHost string
	flag.StringVar(&shortLinkHost, "b", "http://localhost:8080", "host for short link")

	flag.Parse()

	if envHost := os.Getenv("SERVER_ADDRES"); envHost != "" {
		host = envHost
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		shortLinkHost = envBaseURL
	}

	config = new(ConfigBuilder).
		SetHost(host).
		SetShortLinkHost(shortLinkHost).
		Build()
}
