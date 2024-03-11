package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	schemaregistry "github.com/zhangjinpeng87/tistream/pkg/schemaregistry/server"
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

	schemaServer, err := schemaregistry.NewSchemaServer(&globalConfig)
	if err != nil {
		fmt.Printf("Create schema server failed: %v\n", err)
		os.Exit(1)
	}

	if err := schemaServer.Start(); err != nil {
		fmt.Printf("Start schema server failed: %v\n", err)
		os.Exit(1)
	}

	// Wait for the server to stop.
	if err := schemaServer.Wait(); err != nil {
		fmt.Printf("Wait for schema server failed: %v\n", err)
		os.Exit(1)
	}
}
