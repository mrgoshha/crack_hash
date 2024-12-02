package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"log"
	"manager/internal/app"
	"manager/internal/manager"
)

var (
	configPath  string
	mongoUrl    string
	rabbitMQUrl string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/manager.toml", "path to config file")
	flag.StringVar(&mongoUrl, "mongourl", "", "mongo connect url")
	flag.StringVar(&rabbitMQUrl, "rabbiturl", "", "rabbitMQ connect url")

}

func main() {
	flag.Parse()
	config := manager.NewConfig(mongoUrl, rabbitMQUrl)
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal()
	}

	if err := app.Start(config); err != nil {
		log.Fatal()
	}
}
