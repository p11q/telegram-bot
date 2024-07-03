package main

import (
	"flag"
	"log"

	tgClient "telegram_bot/clients/telegram"
	event_consumer "telegram_bot/consumer/event-consumer"
	"telegram_bot/events/telegram"
	"telegram_bot/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
)

// 7260402433:AAHP_bEdwQ8NRgwwX64gxR1AjRr1hwqpRnc

func main() {

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal()
	}
}

func mustToken() string {

	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not spacified")
	}

	return *token
}
