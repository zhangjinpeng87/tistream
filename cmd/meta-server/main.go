package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/zhangjinpeng87/tistream/pkg/metaserver/server"
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

	metaServer := server.NewMetaServer(&globalConfig)
	if err := metaServer.Prepare(); err != nil {
		fmt.Printf("Prepare meta server failed: %v\n", err)
		os.Exit(1)
	}

	if err := metaServer.Start(); err != nil {
		fmt.Printf("Start meta server failed: %v\n", err)
		os.Exit(1)
	}
}
