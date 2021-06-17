package main

import (
	"github.com/al8n/kit-auth/gateway/config"
	"github.com/al8n/kit-auth/gateway/internal/server"
	boot "github.com/al8n/micro-boot"
	"log"
)

func main() {
	var (
		cfg = config.Get()
		srv = server.Get()
	)

	boot.SetDefaultConfigFileType(".yml")
	boot.SetDefaultConfigFileName("config")

	bt, err := boot.New("share", srv, boot.Root{
		Start:          	&boot.Config{
			Configurator: cfg,
		},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	if err = bt.Execute(); err != nil {
		log.Fatal(err)
		return
	}
}


