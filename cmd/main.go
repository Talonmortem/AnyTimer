package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Talonmortem/AnyTimer/internal/config"
	"github.com/Talonmortem/AnyTimer/internal/db"
	"github.com/Talonmortem/AnyTimer/internal/handlers"
	"github.com/Talonmortem/AnyTimer/internal/scheduler"
	"github.com/go-chi/chi"
)

func main() {
	//loading config
	cfg := config.LoadConfig("configs/config.yaml")

	// connect to db
	database, err := db.Connect(db.Config(cfg.Database))
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer database.Close()

	// create router
	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// create task handler
	taskHandler := &handlers.TaskHandler{DB: database}

	// register routes
	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", taskHandler.GetTasks)          // Получить список задач
		r.Get("/{id}", taskHandler.GetTaskByID)   // Получить задачу по ID
		r.Post("/", taskHandler.CreateTask)       // Создать новую задачу
		r.Put("/{id}", taskHandler.UpdateTask)    // Обновить задачу по ID
		r.Delete("/{id}", taskHandler.DeleteTask) // Удалить задачу по ID
	})

	// Создаем планировщик задач
	taskScheduler := scheduler.NewScheduler(database)

	// Запускаем планировщик
	go taskScheduler.Start()

	// start server
	port := cfg.Server.Port
	log.Printf("Starting server on port %d", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
