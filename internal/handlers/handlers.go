package handlers

import (
	"encoding/json" // Для работы с JSON
	"log"
	"net/http" // Для обработки HTTP-запросов
	"strconv"

	// Для преобразования строк в числа
	// Для маршрутизации
	"github.com/Talonmortem/AnyTimer/internal/tasks"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TaskHandler - обработчик задач
type TaskHandler struct {
	DB *pgxpool.Pool
}

// GetTasks обрабатывает запросы на получение списка задач (GET /tasks)

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	taskList, err := tasks.GetAllTasks(r.Context(), h.DB)
	if err != nil {
		http.Error(w, "Failed to retrieve tasks", http.StatusInternalServerError)
		return
	}

	// Устанавливаем Content-Type в JSON
	w.Header().Set("Content-Type", "application/json")

	// Преобразуем задачи в JSON и отправляем клиенту

	if err := json.NewEncoder(w).Encode(taskList); err != nil {
		http.Error(w, "Failed to encode tasks", http.StatusInternalServerError)
		return
	}
}

func (h *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Получаем зачадку по ID
	task, err := tasks.GetTasksByID(r.Context(), h.DB, id)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// Устанавливаем Content-Type в JSON
	w.Header().Set("Content-Type", "application/json")

	// Преобразуем задачи в JSON и отправляем клиенту
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "Failed to encode tasks", http.StatusInternalServerError)
	}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task tasks.Task

	// Разбираем тело запроса

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Println("Failed to decode request body:", err)
		http.Error(w, "Failed to decode request", http.StatusBadRequest)
		return
	}

	newTask, err := tasks.CreateTask(r.Context(), h.DB, task)
	if err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	//Возвращаем созданую задачу клиенту
	if err := json.NewEncoder(w).Encode(newTask); err != nil {
		http.Error(w, "Failed to encode task", http.StatusInternalServerError)
	}
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Разбираем тело запроса

	var task tasks.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Println("Failed to decode request body:", err)
		http.Error(w, "Failed to decode request", http.StatusBadRequest)
		return
	}

	// Обновляем задачу в БД
	updatedTask, err := tasks.UpdateTask(r.Context(), h.DB, id, task)
	if err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	// Возвращаем обновленную задачу клиенту
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedTask); err != nil {
		http.Error(w, "Failed to encode task", http.StatusInternalServerError)
	}
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	if err := tasks.DeleteTask(r.Context(), h.DB, id); err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
