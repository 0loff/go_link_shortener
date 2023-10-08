package main

import (
	"flag"
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

func NewConfigBuilder() error {
	var host string
	flag.StringVar(&host, "a", "localhost:8080", "server host")

	var shortLinkHost string
	flag.StringVar(&shortLinkHost, "b", "localhost:8080", "host for short link")

	flag.Parse()

	config = new(ConfigBuilder).
		SetHost(host).
		SetShortLinkHost(shortLinkHost).
		Build()

	return nil
}
