package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
)

// Receive the message from the channel and send it to Telegram
func telegramSender(bot *tgbotapi.BotAPI, updates chan string) {
	for upd := range updates {
		msg := tgbotapi.NewMessage(ACCOUNT, upd)
		_, err := bot.Send(msg)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// Get Chat ID in order to limit the audience
func getUpdates(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, _ := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "CHAT ID: "+strconv.FormatInt(update.Message.Chat.ID, 10))
		_, err := bot.Send(msg)
		if err != nil {
			fmt.Println(err)
		}
	}
}
