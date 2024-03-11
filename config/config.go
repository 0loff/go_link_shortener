package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
)

const (
	defaultServerAddress = "localhost:8080"
	defaultBaseURL       = "http://localhost:8080"
	defaultLogLevel      = "info"
)

// Config - this is a structure for storing app init params
type Config struct {
	ServerAddress string `json:"server_address"`
	BaseURL       string `json:"base_url"`
	LogLevel      string `json:"log_level"`
	StorageFile   string `json:"file_storage_path"`
	DatabaseDSN   string `json:"database_dsn"`
	TrustedSubnet string `json:"trusted_subnet"`
	EnableHTTPS   bool   `json:"enable_https"`
}

// ConfigBuilder - структура, возвращающая подготовленный кофиг в ходе инициализации приложения
type ConfigBuilder struct {
	c Config
}

// SetServerAddress - метод установки значения хоста в конфиг инициализации приложения
func (cb *ConfigBuilder) SetServerAddress(host string) *ConfigBuilder {
	if cb.c.ServerAddress != "" && host == defaultServerAddress {
		return cb
	}

	cb.c.ServerAddress = host
	return cb
}

// SetBaseURL - метод установки значения хоста для сокращенных urls в конфиг инициализации приложения
func (cb *ConfigBuilder) SetBaseURL(baseURL string) *ConfigBuilder {
	if cb.c.BaseURL != "" && baseURL == defaultBaseURL {
		return cb
	}

	cb.c.BaseURL = baseURL
	return cb
}

// SetLogLevel - метод установки уровня логирования в приложении в зависиомсти от режима запуска при инициализации
func (cb *ConfigBuilder) SetLogLevel(logLevel string) *ConfigBuilder {
	if cb.c.LogLevel != "" && logLevel == defaultLogLevel {
		return cb
	}

	cb.c.LogLevel = logLevel
	return cb
}

// SetStorageFile - метод установки названия и пути к файлу для хранения сокращенных urls в режиме сохранения в файл
func (cb *ConfigBuilder) SetStorageFile(storageFile string) *ConfigBuilder {
	if cb.c.StorageFile != "" && storageFile == "" {
		return cb
	}

	cb.c.StorageFile = storageFile
	return cb
}

// SetDatabaseDSN - метод установки заначения строки конфига для инициализации БД
func (cb *ConfigBuilder) SetDatabaseDSN(databaseDSN string) *ConfigBuilder {
	if cb.c.DatabaseDSN != "" && databaseDSN == "" {
		return cb
	}

	cb.c.DatabaseDSN = databaseDSN
	return cb
}

// SetTrustedSubnet - this is the method for setting the trusted subnet CIDR value
func (cb *ConfigBuilder) SetTrustedSubnet(trustedSubnet string) *ConfigBuilder {
	if cb.c.TrustedSubnet != "" && trustedSubnet == "" {
		return cb
	}

	cb.c.TrustedSubnet = trustedSubnet
	return cb
}

// SetEnableHTTPS - this is setting the https enable flag
func (cb *ConfigBuilder) SetEnableHTTPS(enableHTTPS string) *ConfigBuilder {
	if cb.c.EnableHTTPS && enableHTTPS == "" {
		return cb
	}

	isEnable, err := strconv.ParseBool(enableHTTPS)
	if err != nil {
		cb.c.EnableHTTPS = false
		return cb
	}

	cb.c.EnableHTTPS = isEnable
	return cb
}

// Build - метод для формирования результирующего конфига для инициализации приложения
func (cb *ConfigBuilder) Build() Config {
	return cb.c
}

func (cb *ConfigBuilder) loadFromFile(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&cb.c)
	if err != nil {
		log.Fatal(err)
	}
}

// NewConfigBuilder - метод вызываемый для определения значений конфига инициализации при старте приложения
func NewConfigBuilder() Config {
	var cb ConfigBuilder

	var configFile string
	flag.StringVar(&configFile, "c", "", "path to json config file")

	var serverAddress string
	flag.StringVar(&serverAddress, "a", defaultServerAddress, "server host")

	var baseURL string
	flag.StringVar(&baseURL, "b", defaultBaseURL, "host for short link")

	var logLevel string
	flag.StringVar(&logLevel, "l", defaultLogLevel, "log level")

	var storageFile string
	flag.StringVar(&storageFile, "f", "", "storage file full name")
	// flag.StringVar(&storageFile, "f", "/tmp/short-url-db.json", "storage file full name")

	var databaseDSN string
	flag.StringVar(&databaseDSN, "d", "", "Database DSN config string")
	// flag.StringVar(&databaseDSN, "d", "host=localhost port=5432 user=postgres password=root dbname=urls sslmode=disable", "Database DSN config string")

	var trustedSubnet string
	flag.StringVar(&trustedSubnet, "t", "", "trusted subnet for metrics endpoint")
	// flag.StringVar(&trustedSubnet, "t", "192.168.0.0/24", "trusted subnet for metrics endpoint")

	var enableHTTPS string
	flag.StringVar(&enableHTTPS, "s", "", "Is HTTPS server mode enabled")

	flag.Parse()

	if envConfigFile := os.Getenv("CONFIG"); envConfigFile != "" {
		configFile = envConfigFile
	}

	if envServerSddress := os.Getenv("SERVER_ADDRES"); envServerSddress != "" {
		serverAddress = envServerSddress
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		baseURL = envBaseURL
	}

	if envLoglevel := os.Getenv("LOG_LEVEL"); envLoglevel != "" {
		logLevel = envLoglevel
	}

	if envStorageFile := os.Getenv("FILE_STORAGE_PATH"); envStorageFile != "" {
		storageFile = envStorageFile
	}

	if envDatabaseDSN := os.Getenv("DATABASE_DSN"); envDatabaseDSN != "" {
		databaseDSN = envDatabaseDSN
	}

	if envTrustedSubnet := os.Getenv("TRUSTED_SUBNET"); envTrustedSubnet != "" {
		trustedSubnet = envTrustedSubnet
	}

	if envEnableHTTPS := os.Getenv("ENABLE_HTTPS"); envEnableHTTPS != "" {
		enableHTTPS = envEnableHTTPS
	}

	if configFile != "" {
		cb.loadFromFile(configFile)
	}

	return cb.
		SetServerAddress(serverAddress).
		SetBaseURL(baseURL).
		SetLogLevel(logLevel).
		SetStorageFile(storageFile).
		SetDatabaseDSN(databaseDSN).
		SetTrustedSubnet(trustedSubnet).
		SetEnableHTTPS(enableHTTPS).
		Build()
}
