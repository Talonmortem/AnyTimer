package telegram

import (
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (t *TelegramBot) sendWelcomeMessage(message *tgbotapi.Message, client *http.Client) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Welcome to the AnyTimer bot!")
	_, err := t.bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send welcome message: %v", err)
	}
}

func (t *TelegramBot) sendHelpMessage(message *tgbotapi.Message, client *http.Client) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "This is the help message.")
	_, err := t.bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send help message: %v", err)
	}
}
