package tasks

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Task struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Schedule string `json:"schedule"`
	//TODO: add more fields
}

func GetAllTasks(ctx context.Context, db *pgxpool.Pool) ([]Task, error) {
	rows, err := db.Query(ctx, "SELECT id, name, schedule FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Name, &t.Schedule); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func GetTasksByID(ctx context.Context, db *pgxpool.Pool, id int) (Task, error) {
	row, err := db.Query(ctx, "SELECT id, name, schedule FROM tasks WHERE id = $1", id)
	if err != nil {
		return Task{}, err
	}
	defer row.Close()

	var t Task

	if err := row.Scan(&t.ID, &t.Name, &t.Schedule); err != nil {
		return Task{}, err
	}

	return t, nil
}

func CreateTask(ctx context.Context, db *pgxpool.Pool, t Task) (Task, error) {
	var id int
	if err := db.QueryRow(ctx, "INSERT INTO tasks (name, schedule) VALUES ($1, $2) RETURNING id", t.Name, t.Schedule).Scan(&id); err != nil {
		return Task{}, err
	}

	return Task{ID: id, Name: t.Name, Schedule: t.Schedule}, nil
}

func UpdateTask(ctx context.Context, db *pgxpool.Pool, id int, t Task) (Task, error) {
	_, err := db.Exec(ctx, "UPDATE tasks SET name = $1, schedule = $2 WHERE id = $3", t.Name, t.Schedule, id)
	if err != nil {
		return Task{}, err
	}
	t.ID = id
	return t, nil
}

func DeleteTask(ctx context.Context, db *pgxpool.Pool, id int) error {
	_, err := db.Exec(ctx, "DELETE FROM tasks WHERE id = $1", id)
	return err
}
