package service

import (
	"JourneyPlanner/internal/models"
	"context"
	"errors"
	"fmt"
	"time"
)

type TaskRepository interface {
	AddTask(ctx context.Context, task models.Task, groupID string) error
	GetTaskList(ctx context.Context, userLogin, groupID string) ([]models.Task, error)
	GetTaskById(ctx context.Context, taskID, groupID string) (*models.Task, error)
	UpdateTask(ctx context.Context, taskID string, newTask models.Task) error
	DeleteTask(ctx context.Context, taskID string) error
}

type TaskSrv struct {
	Task  TaskRepository
	Group GroupRepository
}

func NewTaskSrv(taskRepo TaskRepository, groupRepo GroupRepository) *TaskSrv {
	return &TaskSrv{Task: taskRepo, Group: groupRepo}
}

func (s *TaskSrv) CreateTask(ctx context.Context, taskInfo models.CreateTask, userLogin string) error {
	group, err := s.Group.GetGroup(ctx, taskInfo.GroupID, userLogin)
	if err != nil {
		logs.Error(err)
		return errors.New("failed to find group")
	}
	if group == nil{
		return errors.New("group is not found, or you are not a member of it")
	}
	if group.LeaderLogin != userLogin {
		logs.Error(err)
		return errors.New("you have no permissions to do this")
	}
	dateTimeStr := fmt.Sprintf("%sT%s:00Z", taskInfo.StartTime.StartDate, taskInfo.StartTime.StartTime)

	startTime, err := time.Parse(dateformat, dateTimeStr)
	if err != nil {
		logs.Error(err)
		return fmt.Errorf("invalid date or time format: %v", err)
	}
	now := time.Now().UTC()
	if startTime.Before(now) {
		return errors.New("you cant add tasks to past time")
	}
	totalDuration := calculateDuration(taskInfo.Duration)
	endTime := startTime.Add(time.Duration(totalDuration) * time.Minute)

	newTask := models.Task{
		Title:     taskInfo.Title,
		StartTime: startTime,
		Duration:  totalDuration,
		EndTime:   endTime,
	}
	existingTasks, err := s.Task.GetTaskList(ctx, userLogin, taskInfo.GroupID)
	if err != nil {
		logs.Error(err)
		return errors.New("System error")
	}

	for _, existingTask := range existingTasks {
		if doTasksOverlap(existingTask, newTask) {
			return fmt.Errorf("task overlaps with an existing task: %s", existingTask.Title)
		}
	}
	err = s.Task.AddTask(ctx, newTask, taskInfo.GroupID)
	if err != nil {
		logs.Error(err)
		return errors.New("System error")
	}
	return nil
}

func doTasksOverlap(existingTask, newTask models.Task) bool {
	return existingTask.EndTime.After(newTask.StartTime) && existingTask.StartTime.Before(newTask.EndTime)
}

func (s *TaskSrv) GetTaskList(ctx context.Context, groupID, userLogin string) ([]models.Task, error) {
	group, err := s.Group.GetGroup(ctx, groupID, userLogin)
	if err != nil {
		logs.Error(err)
		return nil, errors.New("failed to find group")
	}
	if group == nil{
		return nil, errors.New("group is not found, or you are not a member of it")
	}
	tasks, err := s.Task.GetTaskList(ctx, userLogin, groupID)
	if err != nil {
		logs.Error(err)
		return nil, errors.New("System error")
	}
	return tasks, nil
}

var dateformat string = "2006-01-02T15:04:05Z"

func (s *TaskSrv) UpdateTask(ctx context.Context, taskID, userLogin string, updateTask models.CreateTask) error {
	group, err := s.Group.GetGroup(ctx, updateTask.GroupID, userLogin)
	if err != nil {
		logs.Error(err)
		return errors.New("failed to find group")
	}
	if group == nil{
		return errors.New("group is not found, or you are not a member of it")
	}
	if group.LeaderLogin != userLogin {
		logs.Error(err)
		return errors.New("you have no permissions to do this")
	}
	task, err := s.Task.GetTaskById(ctx, taskID, updateTask.GroupID)
	if err != nil {
		logs.Error(err)
		return errors.New("task was not found")
	}
	var startTime time.Time
	if !updateTask.StartTime.IsFullEmpty() {
		dateTimeStr := fmt.Sprintf("%sT%s:00Z", updateTask.StartTime.StartDate, updateTask.StartTime.StartTime)
		startTime, err = time.Parse(dateformat, dateTimeStr)
		if err != nil {
			logs.Error(err)
			return fmt.Errorf("invalid date or time format: %v", err)
		}
		now := time.Now().UTC()
		if startTime.Before(now) {
			logs.Error(err)
			return errors.New("you cant add tasks to past time")
		}
	} else {
		startTime = task.StartTime
	}
	var endTime time.Time
	var totalDuration int
	startTimeProvided := !updateTask.StartTime.IsFullEmpty()
	durationProvided := !updateTask.Duration.IsEmpty()
	switch {
	case startTimeProvided && durationProvided:
		totalDuration = calculateDuration(updateTask.Duration)
		endTime = startTime.Add(time.Duration(totalDuration) * time.Minute)
	case !startTimeProvided && durationProvided:
		totalDuration = calculateDuration(updateTask.Duration)
		endTime = task.StartTime.Add(time.Duration(totalDuration) * time.Minute)
	case startTimeProvided && !durationProvided:
		totalDuration = task.Duration
		endTime = startTime.Add(time.Duration(totalDuration) * time.Minute)
	}
	updates := models.Task{
		Title:     updateTask.Title,
		StartTime: startTime,
		Duration:  totalDuration,
		EndTime:   endTime,
	}
	existingTasks, err := s.Task.GetTaskList(ctx, userLogin, updateTask.GroupID)
	if err != nil {
		logs.Error(err)
		return errors.New("System error")
	}

	for _, existingTask := range existingTasks {
		if existingTask.ID.Hex() != taskID {
			if doTasksOverlap(existingTask, updates) {
				return fmt.Errorf("task overlaps with an existing task: %s", existingTask.Title)
			}
		}
	}
	err = s.Task.UpdateTask(ctx, taskID, updates)
	if err != nil {
		logs.Error(err)
		return errors.New("System error")
	}

	return nil
}

const (
	HoursInDay    = 24
	MinutesInHour = 60
)

func calculateDuration(dur models.Duration) int {
	totalMinutes := (dur.DurDays * HoursInDay * MinutesInHour) + (dur.DurHours * MinutesInHour) + dur.DurMinutes
	return totalMinutes
}
func (s *TaskSrv) DeleteTask(ctx context.Context, taskID, groupID, userLogin string) error {
	group, err := s.Group.GetGroup(ctx, groupID, userLogin)
	if err != nil {
		logs.Error(err)
		return errors.New("failed to find group")
	}
	if group == nil{
		return errors.New("group is not found, or you are not a member of it")
	}
	if group.LeaderLogin != userLogin {
		return errors.New("you have no permissions to do this")
	}
	err = s.Task.DeleteTask(ctx, taskID)
	if err != nil {
		logs.Error(err)
		return errors.New("System error")
	}
	return nil
}
