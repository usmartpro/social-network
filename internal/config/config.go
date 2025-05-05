package config

import "os"

type Conf struct {
	Logger  LoggerConf
	Storage StorageConf
	HTTP    HTTPConf
	Cache   CacheConf
}

type HTTPConf struct {
	Host string
	Port string
}

type LoggerConf struct {
	Level string
	File  string
}

type StorageConf struct {
	Type string
	Dsn  string
}

type CacheConf struct {
	Type string
	Dsn  string
}

func LoadConfiguration() (*Conf, error) {
	return &Conf{
		Logger: LoggerConf{
			Level: os.Getenv("LOG_LEVEL"),
			File:  os.Getenv("LOG_FILENAME"),
		},
		Storage: StorageConf{
			Type: os.Getenv("STORAGE_TYPE"),
			Dsn:  os.Getenv("STORAGE_DSN"),
		},
		HTTP: HTTPConf{
			Host: os.Getenv("HTTP_HOST"),
			Port: os.Getenv("HTTP_PORT"),
		},
		Cache: CacheConf{
			Dsn: os.Getenv("CACHE_DSN"),
		},
	}, nil
}
