package main

import (
	"ispjournalctl/cmd"
	"ispjournalctl/config"
	"ispjournalctl/service"
)

var (
	version = "1.0.0"
)

func main() {
	if cfg, err := config.Load(""); err != nil {
		panic(err)
	} else {
		service.JournalServiceClient.ReceiveConfiguration(cfg.GateHost)
	}
	cmd.Execute(version)
}
