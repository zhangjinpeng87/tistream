package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	remotesorter "github.com/zhangjinpeng87/tistream/pkg/remotesorter/server"
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

	sorterServer, err := remotesorter.NewSorterServer(&globalConfig)
	if err != nil {
		fmt.Printf("Create sorter server failed: %v\n", err)
		os.Exit(1)
	}
	if err := sorterServer.Prepare(); err != nil {
		fmt.Printf("Prepare sorter server failed: %v\n", err)
		os.Exit(1)
	}

	if err := sorterServer.Start(); err != nil {
		fmt.Printf("Start sorter server failed: %v\n", err)
		os.Exit(1)
	}

	// Wait for the server to stop.
	if err := sorterServer.Wait(); err != nil {
		fmt.Printf("Wait for sorter server failed: %v\n", err)
		os.Exit(1)
	}
}
