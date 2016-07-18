package config

import (
	"encoding/json"
	"errors"
	"os"
)

type Configuration struct {
	ListenAddress  int            `json:"listen_address"`
	DatabaseConfig DatabaseConfig `json:"database_config"`
	FacebookConfig FacebookConfig `json:"facebook_config"`
	JWTKey         string         `json:"jwt_key"`
	SigningKey     []byte
}
type DatabaseConfig struct {
	DatabaseUri string `json:"database_address"`
	Username    string `json:"user"`
	Password    string `json:"password"`
	Port        int    `json:"port"`
	DBName      string `json:"database_name"`
}

type FacebookConfig struct {
	ID              string `json:"fb_id"`
	Key             string `json:"fb_key"`
	CallbackAddress string `json:"server_callback_address"`
}

var (
	ErrCantFindConfig = errors.New("Config file is missing")
)

func NewConfig(configfile string) (*Configuration, error) {
	if _, err := os.Stat(configfile); os.IsNotExist(err) {
		return nil, ErrCantFindConfig
	}

	configuration := Configuration{}

	file, _ := os.Open(configfile)
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&configuration)
	if err != nil {
		return nil, err
	}

	configuration.SigningKey = []byte(configuration.JWTKey)

	return &configuration, nil
}

func MustNewConfig(configfile string) *Configuration {
	config, err := NewConfig(configfile)

	if err != nil {
		panic(err)
	}

	return config
}
