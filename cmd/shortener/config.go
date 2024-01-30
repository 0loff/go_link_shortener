package main

import (
	"flag"
	"os"
)

var config Config

// Config - структура для хранения параметров инициализации приложения
type Config struct {
	Host         string
	ShortURLHost string
	LogLevel     string
	StorageFile  string
	DatabaseDSN  string
}

// ConfigBuilder - структура, возвращающая подготовленный кофиг в ходе инициализации приложения
type ConfigBuilder struct {
	config Config
}

// SetHost - метод установки значения хоста в конфиг инициализации приложения
func (cb ConfigBuilder) SetHost(host string) ConfigBuilder {
	cb.config.Host = host
	return cb
}

// SetShortLinkHost - метод установки значения хоста для сокращенных urls в конфиг инициализации приложения
func (cb ConfigBuilder) SetShortLinkHost(shortURLHost string) ConfigBuilder {
	cb.config.ShortURLHost = shortURLHost
	return cb
}

// SetLogLevel - метод установки уровня логирования в приложении в зависиомсти от режима запуска при инициализации
func (cb ConfigBuilder) SetLogLevel(logLevel string) ConfigBuilder {
	cb.config.LogLevel = logLevel
	return cb
}

// SetStorageFile - метод установки названия и пути к файлу для хранения сокращенных urls в режиме сохранения в файл
func (cb ConfigBuilder) SetStorageFile(storageFile string) ConfigBuilder {
	cb.config.StorageFile = storageFile
	return cb
}

// SetDatabaseDSN - метод установки заначения строки конфига для инициализации БД
func (cb ConfigBuilder) SetDatabaseDSN(databaseDSN string) ConfigBuilder {
	cb.config.DatabaseDSN = databaseDSN
	return cb
}

// Build - метод для формирования результирующего конфига для инициализации приложения
func (cb ConfigBuilder) Build() Config {
	return cb.config
}

// NewConfigBuilder - метод вызываемый для определения значений конфига инициализации при старте приложения
func NewConfigBuilder() {
	var host string
	flag.StringVar(&host, "a", "localhost:8080", "server host")

	var shortURLHost string
	flag.StringVar(&shortURLHost, "b", "http://localhost:8080", "host for short link")

	var logLevel string
	flag.StringVar(&logLevel, "l", "info", "log level")

	var storageFile string
	flag.StringVar(&storageFile, "f", "", "storage file full name")
	// flag.StringVar(&storageFile, "f", "/tmp/short-url-db.json", "storage file full name")

	var databaseDSN string
	flag.StringVar(&databaseDSN, "d", "", "Database DSN config string")
	// flag.StringVar(&databaseDSN, "d", "host=localhost port=5432 user=postgres password=root dbname=urls sslmode=disable", "Database DSN config string")

	flag.Parse()

	if envHost := os.Getenv("SERVER_ADDRES"); envHost != "" {
		host = envHost
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		shortURLHost = envBaseURL
	}

	if envLoglevel := os.Getenv("LOG_LEVEL"); envLoglevel != "" {
		logLevel = envLoglevel
	}

	if envStorageFile := os.Getenv("FILE_STORAGE_PATH"); envStorageFile != "" {
		storageFile = envStorageFile
	}

	if envStorageFile := os.Getenv("DATABASE_DSN"); envStorageFile != "" {
		databaseDSN = envStorageFile
	}

	config = new(ConfigBuilder).
		SetHost(host).
		SetShortLinkHost(shortURLHost).
		SetLogLevel(logLevel).
		SetStorageFile(storageFile).
		SetDatabaseDSN(databaseDSN).
		Build()
}
