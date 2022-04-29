package configs

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Server   ServerConfig
	SeverCSV SeverCSVConfig
	MongoDB  MongoDBConfig
	Log      LogConfig
}

type ServerConfig struct {
	Host    string `envconfig:"host"`
	Port    string `envconfig:"port"`
	Network string `envconfig:"network"`
}

type SeverCSVConfig struct {
	Host   string `envconfig:"host"`
	Port   string `envconfig:"port"`
	Folder string `envconfig:"folder"`
}

type MongoDBConfig struct {
	User     string `envconfig:"user"`
	Password string `envconfig:"password"`
}

type LogConfig struct {
	Prefix string `envconfig:"prefix"`
}

const (
	serverGroup    = "server"
	csvServerGroup = "csv_server"
	mongodbGroup   = "mongodb"
	logGroup       = "log"
)

func NewConfig() (*Config, error) {
	config := new(Config)

	if err := envconfig.Process(serverGroup, &config.Server); err != nil {
		return &Config{}, err
	}
	if err := envconfig.Process(csvServerGroup, &config.SeverCSV); err != nil {
		return &Config{}, err
	}
	if err := envconfig.Process(mongodbGroup, &config.MongoDB); err != nil {
		return &Config{}, err
	}
	if err := envconfig.Process(logGroup, &config.Log); err != nil {
		return &Config{}, err
	}
	return config, nil
}
