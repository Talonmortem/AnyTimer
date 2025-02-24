package telegram

import (
	"fmt"
	"log"
	"net/http"

	api "github.com/Talonmortem/AnyTimer/internal/api"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type TelegramBot struct {
	bot       *tgbotapi.BotAPI
	userState map[int64]string //хранит состояние пользователя
	tasks     map[int64]api.Task
}

func NewTelegramBot(token string) (*TelegramBot, error) {
	// init bot
	log.Println("Starting bot...", token)
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	return &TelegramBot{bot: bot,
		userState: make(map[int64]string),
		tasks:     make(map[int64]api.Task)}, nil
}

func (t *TelegramBot) Start(client *http.Client) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := t.bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Failed to get updates: %v", err)
	}

	//handle updates
	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		t.handleMessage(update.Message, client)
	}
}

// send message

func (t *TelegramBot) sendMessage(chatid int64, text string) {
	msg := tgbotapi.NewMessage(chatid, text)
	_, err := t.bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

// handle message

func (t *TelegramBot) handleMessage(message *tgbotapi.Message, client *http.Client) {
	switch t.userState[message.Chat.ID] {
	case "awaiting_task_name":
		t.getTaskName(message, client)
		return
	case "awaiting_task_schedule":
		t.getTaskSchedule(message, client)
		return
	case "task_complite":
		t.sendTaskCompliteMessage(message, client)
		return
	default:
		switch message.Text {
		case "/start":
			t.sendWelcomeMessage(message, client)
		case "/help":
			t.sendHelpMessage(message, client)
		case "/tasks":
			t.sendTasksList(message, client)
		case "/addtask":
			t.sendAddTaskMessage(message, client)
		default:
			t.sendUnknownCommandMessage(message, client)
		}
	}
}

func (t *TelegramBot) sendTasksList(message *tgbotapi.Message, client *http.Client) {
	// Send API request to get tasks list
	tasks, err := api.GetTasks(client)
	if err != nil {
		t.sendMessage(message.Chat.ID, "Failed to get tasks list.")
		return
	}

	var taskListStr string

	if len(tasks) == 0 {
		taskListStr = "No tasks found."
	}

	for _, task := range tasks {
		taskListStr += fmt.Sprintf("Task ID: %d, Name: %s, Schedule: %s\n", task.ID, task.Name, task.Schedule)
	}

	t.sendMessage(message.Chat.ID, taskListStr)
}

func (t *TelegramBot) sendUnknownCommandMessage(message *tgbotapi.Message, client *http.Client) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Unknown command.")
	_, err := t.bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send unknown command message: %v", err)
	}

	t.sendHelpMessage(message, client)
}
