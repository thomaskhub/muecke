package main

import (
	"flag"
	"fmt"

	"github.com/thomaskhub/muecke/bridge"
	"github.com/thomaskhub/muecke/utils"
)

func main() {
	configFile := flag.String("config", "config.yaml", "Path to the config file")
	flag.Parse()

	cfg := utils.ParseConfig(*configFile)
	fmt.Printf("cfg: %v\n", cfg)

	bridge.StartBridge(cfg)

	select {}
}
