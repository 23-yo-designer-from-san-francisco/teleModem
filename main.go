package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(API_KEY)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = false
	updates := make(chan string)
	go telegramSender(bot, updates)
	updates <- "***Up and running***"
	go modemHandler(updates)
	select {}
}
