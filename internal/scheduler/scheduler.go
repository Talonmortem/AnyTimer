package scheduler

/*
Cron выражения имеют следующий формат:

* * * * * <command>
| | | | |
| | | | +--- День недели (0 - 7) (воскресенье = 0 или 7)
| | | +----- Месяц (1 - 12)
| | +------- День месяца (1 - 31)
| +--------- Час (0 - 23)
+----------- Минуты (0 - 59)
Примеры:

* * * * * — каждый момент времени (каждую минуту).
0 9 * * * — каждый день в 9:00.
* /5 * * * * — каждые 5 минут.
*/

import (
	"context"
	"log"

	"github.com/Talonmortem/AnyTimer/internal/tasks"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3" // Библиотека для cron
)

type TaskScheduler struct {
	DB   *pgxpool.Pool
	Cron *cron.Cron
}

// Новый планировщик задач

func NewScheduler(db *pgxpool.Pool) *TaskScheduler {

	// Создаем новый cron
	c := cron.New()

	// Инициализируем планировщик
	return &TaskScheduler{
		DB:   db,
		Cron: c,
	}
}

// Запускаем планировщик
func (ts *TaskScheduler) Start() {
	// Запускаем планировщик
	_, err := ts.Cron.AddFunc("*/5 * * * *", func() {
		log.Println("Cron job is running")
	})

	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}

	// Запускаем планировщик
	ts.Cron.Start()
}

// Останавливаем планировщик
func (ts *TaskScheduler) Stop() {
	// Останавливаем планировщик
	ts.Cron.Stop()
}

// Запуск задач из бд по расписанию

func (ts *TaskScheduler) RunTaskBySchedule() {
	tasks, err := tasks.GetAllTasks(context.Background(), ts.DB)
	if err != nil {
		log.Println("Failed to get tasks from database:", err)
		return
	}

	for _, task := range tasks {
		_, err := ts.Cron.AddFunc(task.Schedule, func() {
			//Логика выполнения задачи
			log.Printf("Running task: %s", task.Name)
			//TODO: Вызов функции из tasks
		})
		if err != nil {
			log.Printf("Failed to add task %s to cron job: %v", task.Name, err)
		}
	}
}
