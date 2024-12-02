package main

import (
	"flag"
	"log"
	"worker/internal/app"
	"worker/internal/worker"
)

var (
	baseUrl     string
	rabbitMQUrl string
)

func init() {
	flag.StringVar(&baseUrl, "baseurl", "", "base url")
	flag.StringVar(&rabbitMQUrl, "rabbiturl", "", "rabbitMQ connect url")
}

func main() {
	flag.Parse()
	config := worker.NewConfig(baseUrl, rabbitMQUrl)
	if err := app.Start(config); err != nil {
		log.Fatal()
	}
}
