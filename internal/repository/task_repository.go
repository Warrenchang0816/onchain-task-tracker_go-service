package repository

import (
	"database/sql"
	"go-service/internal/model"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) FindAll() ([]model.Task, error) {
	rows, err := r.db.Query(`
		SELECT id, title, description, status, priority, due_date, created_at, updated_at
		FROM tasks
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]model.Task, 0)

	for rows.Next() {
		var task model.Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.DueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) Create(task model.Task) (int64, error) {
	var id int64

	err := r.db.QueryRow(`
		INSERT INTO tasks (title, description, status, priority, due_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id
	`,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.DueDate,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *TaskRepository) FindByID(id int64) (*model.Task, error) {
	var task model.Task

	err := r.db.QueryRow(`
		SELECT id, title, description, status, priority, due_date, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *TaskRepository) Update(task model.Task) error {
	_, err := r.db.Exec(`
		UPDATE tasks
		SET title = $1,
			description = $2,
			status = $3,
			priority = $4,
			due_date = $5,
			updated_at = NOW()
		WHERE id = $6
	`,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.DueDate,
		task.ID,
	)

	return err
}

func (r *TaskRepository) UpdateStatus(id int64, status string) error {
	_, err := r.db.Exec(`
		UPDATE tasks
		SET status = $1,
			updated_at = NOW()
		WHERE id = $2
	`, status, id)

	return err
}
