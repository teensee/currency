package main

import (
	"Currency/internal/app"
	"Currency/internal/config"
	"log"
)

func main() {
	log.Print("Startup, load config")
	cfg := config.GetConfig()

	kernel, err := app.NewKernel(cfg)

	if err != nil {
		log.Fatal(err)
	}

	log.Print("Run application")
	kernel.Run()
}
