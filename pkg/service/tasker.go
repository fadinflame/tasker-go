package service

import (
	"database/sql"
	"tasker/pkg/consts"
	"tasker/pkg/models"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/log"
	_ "github.com/mattn/go-sqlite3"
)

type TaskService struct {
	sql    *sql.DB
	Logger log.Logger
}

func NewTaskService(db *sql.DB, logger log.Logger) *TaskService {
	return &TaskService{sql: db, Logger: logger}
}

func (s *TaskService) Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS tasks(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        text TEXT,
		is_completed INTEGER DEFAULT 0);
    `

	_, err := s.sql.Exec(query)
	return err
}

func (s *TaskService) CreateTask(task models.Task) (*models.Task, error) {
	level.Debug(s.Logger).Log("service", "CreateTask", "Task Info", task.GetInfo())
	res, err := s.sql.Exec(consts.CreateTaskQuery, task.Title, task.Text)
	if err != nil {
		return nil, err
	}
	task.Id, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (s *TaskService) GetTask(id int64) (*models.Task, error) {
	level.Debug(s.Logger).Log("service", "GetTask", "TaskId", id)
	res := s.sql.QueryRow(consts.GetTaskQuery, id)
	task := models.Task{}
	if err := res.Scan(&task.Id, &task.Title, &task.Text, &task.IsCompleted); err != nil {
		return &task, err
	}
	return &task, nil
}

func (s *TaskService) UpdateTask(id int64, task models.Task) (bool, error) {
	level.Debug(s.Logger).Log("service", "UpdateTask", "TaskId", id, "Data", task.GetInfo())
	_, err := s.sql.Exec(consts.UpdateTaskQuery, task.Title, task.Text, task.IsCompleted, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *TaskService) DeleteTask(id int64) (bool, error) {
	level.Debug(s.Logger).Log("service", "DeleteTask", "TaskId", id)
	_, err := s.sql.Exec(consts.DeleteTaskQuery, id)
	if err != nil {
		return false, err
	}
	return true, nil
}
