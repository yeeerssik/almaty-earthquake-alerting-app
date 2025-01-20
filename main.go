package main

import (
	"almeqapp/config"
	"almeqapp/internal/client"
	"almeqapp/internal/scheduler"
	"almeqapp/internal/service"
	"almeqapp/internal/usecase"
	"log"
)

func main() {
	c, err := config.New()
	if err != nil {
		log.Panic(err)
	}

	tgService, err := service.New(c)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Successfully authenticated as %s\n", tgService.Bot.Self.UserName)

	eqClient := client.New(c)

	eqService := usecase.New(c, eqClient, tgService)

	scheduler.StartScheduler(c, eqService)
}
