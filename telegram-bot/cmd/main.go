package main

import (
	"log"
	"net/http"

	"github.com/Talonmortem/AnyTimer/internal/config"
	"github.com/Talonmortem/AnyTimer/telegram-bot/internal/telegram"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.LoadConfig("configs/config.yaml")

	// Создаём Telegram бота
	bot, err := telegram.NewTelegramBot(cfg.Telegram.Token)
	if err != nil {
		log.Fatalf("Failed to create Telegram bot: %v", err)
	}

	// Создаём HTTP клиент для взаимодействия с основным сервисом
	client := &http.Client{}

	// Запускаем бота
	go bot.Start(client)

	// Ожидаем завершения работы
	select {}
}
