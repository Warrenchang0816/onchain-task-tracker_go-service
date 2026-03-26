package service

import (
	"go-service/internal/model"
	"go-service/internal/repository"
)

type TaskService struct {
	taskRepo *repository.TaskRepository
}

func NewTaskService(taskRepo *repository.TaskRepository) *TaskService {
	return &TaskService{taskRepo: taskRepo}
}

func (s *TaskService) GetTasks() ([]model.Task, error) {
	return s.taskRepo.FindAll()
}

func (s *TaskService) CreateTask(task model.Task) (int64, error) {
	return s.taskRepo.Create(task)
}

func (s *TaskService) UpdateTask(task model.Task) error {
	return s.taskRepo.Update(task)
}

func (s *TaskService) UpdateTaskStatus(id int64, status string) error {
	return s.taskRepo.UpdateStatus(id, status)
}
