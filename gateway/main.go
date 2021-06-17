package main

import (
	"github.com/al8n/kit-auth/gateway/config"
	"github.com/al8n/kit-auth/gateway/internal/server"
	"github.com/al8n/micro-boot"
)

func main() {
	var (
		cfg = config.Get()
		srv = server.Get()
	)

	boot.Init("gateway", srv, cfg)
	//commands.Execute()
}


