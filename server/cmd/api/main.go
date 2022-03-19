package main

import (
	"log"
	"math/rand"
	"server/config"
	"server/internal/kernel/usecase"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	log.Println("Starting server")
	cfgFile, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
		//return
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
		//return
	}
	log.Println("Config loaded")

	kernelUS := usecase.NewKernelUS(cfg)
	kernelUS.StartListening()

}
