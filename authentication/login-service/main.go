package main

import (
	"github.com/al8n/kit-auth/authentication/login-service/config"
	"github.com/al8n/kit-auth/authentication/login-service/pkg/server"
	boot "github.com/al8n/micro-boot"
	"log"
)

func main() {
	var (
		cfg = config.GetConfig()
	)

	boot.SetDefaultConfigFileType(".yml")
	boot.SetDefaultConfigFileName("config")


	bt, err := boot.New("Login", &server.Server{}, boot.Root{
		Start:          	&boot.Config{
			Configurator: cfg,
		},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	if err := bt.Execute(); err != nil {
		log.Fatal(err)
		return
	}
}
