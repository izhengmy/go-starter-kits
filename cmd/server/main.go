package main

import (
	"app/internal/config"
	"app/pkg/viperx"
	"flag"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "conf", "config.yaml", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()

	var conf config.Config
	if err := viperx.Scan(configPath, &conf); err != nil {
		panic(err)
	}

	server, cleanup, err := wireServer(conf)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err := server.Start(); err != nil {
		panic(err)
	}
}
