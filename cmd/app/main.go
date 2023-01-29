package main

import (
	"Currency/internal/app"
	"Currency/internal/config"
	"log"
)

func main() {
	log.Print("Startup, load config")
	cfg := config.GetConfig()

	kernel := app.NewKernel(cfg)
	kernel.
		ConfigureDatabase().
		ConfigureServiceLocator().
		ConfigureRoutes().
		AfterInitializationEvents().
		End()
}
