package main

import (
	"log"

	"github.com/mmarci96/codebox-wss-addon/pkg/config"
)

func main() {
	conf := config.Load()
	log.Println("Starting server on", conf.Host, conf.Port)

}
