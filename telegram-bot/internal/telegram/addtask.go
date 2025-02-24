package telegram

import (
	"fmt"
	"log"
	"net/http"

	api "github.com/Talonmortem/AnyTimer/internal/api"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// add task functions

// send add task message
func (t *TelegramBot) sendAddTaskMessage(message *tgbotapi.Message, client *http.Client) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Please enter the task name.")
	_, err := t.bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send add task message: %v", err)
	}

	t.userState[message.Chat.ID] = "awaiting_task_name"
}

func (t *TelegramBot) getTaskName(message *tgbotapi.Message, client *http.Client) {
	t.tasks[message.Chat.ID] = api.Task{Name: message.Text}
	t.userState[message.Chat.ID] = "awaiting_task_schedule"
	t.processTaskSchedule(message, client)
}

// process task schedule

func (t *TelegramBot) processTaskSchedule(message *tgbotapi.Message, client *http.Client) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Please enter the task schedule.")
	_, err := t.bot.Send(msg)
	if err != nil {
		log.Printf("Failed to process task schedule: %v", err)
	}

	t.userState[message.Chat.ID] = "awaiting_task_schedule"
}

func (t *TelegramBot) getTaskSchedule(message *tgbotapi.Message, client *http.Client) {
	t.userState[message.Chat.ID] = "task_complite"
	t.tasks[message.Chat.ID] = api.Task{Name: t.tasks[message.Chat.ID].Name, Schedule: message.Text}
	t.sendTaskCompliteMessage(message, client)
}

// complete task

func (t *TelegramBot) sendTaskCompliteMessage(message *tgbotapi.Message, client *http.Client) {
	task, err := api.CreateTask(client, t.tasks[message.Chat.ID])
	if err != nil {
		t.sendMessage(message.Chat.ID, fmt.Sprintf("Failed to create task: %v", err))
		return
	}
	t.sendMessage(message.Chat.ID, (fmt.Sprintf("Task ID: %d, Name: %s, Schedule: %s\n", task.ID, task.Name, task.Schedule)))
	delete(t.userState, message.Chat.ID)
	delete(t.tasks, message.Chat.ID)
}

//
