package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Port              int
	Host              string
	Username          string
	SecretStoragePath string
	WsEndpoing        string
}

func Load() *Config {
	conf := &Config{}
	port, _ := os.LookupEnv("PORT")
	host, _ := os.LookupEnv("HOST")
	username, _ := os.LookupEnv("USERNAME")
	secretStoragePath, _ := os.LookupEnv("SECRET_STORE_PATH")
	wsEndpoint, _ := os.LookupEnv("WS_ENDPOINT")

	portNumber, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal("Error converting port number.")
		conf.Port = 8080
	} else {
		conf.Port = portNumber
	}
	conf.Host = host
	conf.Username = username
	conf.SecretStoragePath = secretStoragePath
	conf.WsEndpoing = wsEndpoint

	return conf
}
