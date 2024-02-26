package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/zhangjinpeng87/tistream/pkg/dispatcher/server"
	"github.com/zhangjinpeng87/tistream/pkg/utils"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "config.toml", "Config file")
	flag.Parse()

	var globalConfig utils.GlobalConfig
	if _, err := toml.DecodeFile(configFile, &globalConfig); err != nil {
		fmt.Printf("Decode config file failed: %v\n", err)
		os.Exit(1)
	}

	metaServer, err := server.NewDispatchServer(&globalConfig.Dispatcher)
	if err != nil {
		fmt.Printf("Create meta server failed: %v\n", err)
		os.Exit(1)
	}
	if err := metaServer.Prepare(); err != nil {
		fmt.Printf("Prepare meta server failed: %v\n", err)
		os.Exit(1)
	}

	if err := metaServer.Start(); err != nil {
		fmt.Printf("Start meta server failed: %v\n", err)
		os.Exit(1)
	}

	// Wait for the server to stop.
	if err := metaServer.Stop(); err != nil {
		fmt.Printf("Wait for meta server failed: %v\n", err)
		os.Exit(1)
	}
}