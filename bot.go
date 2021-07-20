package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
	"time"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(API_KEY)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	updates := make(chan string)
	go telegramSender(bot, updates)
	go modemHandler(updates)
	time.Sleep(100 * time.Second)
}

// Receive the message from the channel and send it to Telegram
func telegramSender(bot *tgbotapi.BotAPI, updates chan string) {
	for upd := range updates {
		msg := tgbotapi.NewMessage(ACCOUNT, upd)
		bot.Send(msg)
	}
}

// Get new modem messages and send them to channel
func modemHandler(updates chan string) {
	for ;; {
		time.Sleep(5 * time.Second)
		updates <- "You have a new message"
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
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "CHAT ID: " +strconv.FormatInt(update.Message.Chat.ID, 10))
		bot.Send(msg)
	}
}