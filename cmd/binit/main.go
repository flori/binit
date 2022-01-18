package main

import (
	"log"
	"os"

	"github.com/flori/binit"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	var config binit.Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err.Error())
	}
	binit.Run(config, os.Args[1:])
}

func init() {
	binit.ConfigureLogging()
}
