package main

import (
	"fmt"
	"log"

	"github.com/mmarci96/codebox-wss-addon/pkg/config"
	"github.com/mmarci96/codebox-wss-addon/pkg/server"
)

func main() {
	conf := config.Load()
	log.Println("Starting server on", conf.Host, conf.Port)
	server := server.Setup(conf)

	addr := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	log.Printf("Server running at http://%s", addr)
	if err := server.Run(addr); err != nil {
		log.Fatal("Server error:", err)
	}

}
