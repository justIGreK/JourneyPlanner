package service

import (
	"JourneyPlanner/internal/models"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskRepository interface {
	AddTask(ctx context.Context, task models.Task) error
	GetTaskList(ctx context.Context, userLogin string, groupID primitive.ObjectID) ([]models.Task, error)
	UpdateTask(ctx context.Context, taskID primitive.ObjectID, updates models.Task) error
	DeleteTask(ctx context.Context, taskID primitive.ObjectID) error
	GetTaskById(ctx context.Context, taskID, groupID primitive.ObjectID) (*models.Task, error)
}

type TaskSrv struct {
	TaskRepository
	GroupRepository
}

func NewTaskSrv(taskRepo TaskRepository, groupRepo GroupRepository) *TaskSrv {
	return &TaskSrv{TaskRepository: taskRepo, GroupRepository: groupRepo}
}

func (s *TaskSrv) CreateTask(ctx context.Context, taskInfo models.CreateTask, userLogin string) error {
	groupOID, err := primitive.ObjectIDFromHex(taskInfo.GroupID)
	if err != nil {
		return err
	}
	_, err = s.GroupRepository.GetGroupById(ctx, groupOID, userLogin)
	if err != nil {
		return errors.New("this group is not exist or you are not member of it")
	}
	dateTimeStr := fmt.Sprintf("%sT%s:00Z", taskInfo.StartTime.StartDate, taskInfo.StartTime.StartTime)

	startTime, err := time.Parse(dateformat, dateTimeStr)
	if err != nil {
		return fmt.Errorf("invalid date or time format: %v", err)
	}
	now := time.Now().UTC()
	if startTime.Before(now) {
		return errors.New("you cant add tasks to past time")
	}
	totalDuration := calculateDuration(taskInfo.Duration) 
	endTime := startTime.Add(time.Duration(totalDuration) * time.Minute)

	newTask := models.Task{
		GroupID:   groupOID,
		Title:     taskInfo.Title,
		StartTime: startTime,
		Duration:  totalDuration,
		EndTime:   endTime,
	}
	existingTasks, err := s.TaskRepository.GetTaskList(ctx, userLogin, groupOID)
	if err != nil {
		return err
	}

	for _, existingTask := range existingTasks {
		if doTasksOverlap(existingTask, newTask) {
			return fmt.Errorf("task overlaps with an existing task: %s", existingTask.Title)
		}
	}
	err = s.TaskRepository.AddTask(ctx, newTask)
	if err != nil {
		return err
	}
	return nil
}

func doTasksOverlap(existingTask, newTask models.Task) bool {
	return existingTask.EndTime.After(newTask.StartTime) && existingTask.StartTime.Before(newTask.EndTime)
}

func (s *TaskSrv) GetTaskList(ctx context.Context, groupID, userLogin string) ([]models.Task, error) {
	groupOID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return nil, err
	}
	_, err = s.GroupRepository.GetGroupById(ctx, groupOID, userLogin)
	if err != nil {
		return nil, errors.New("this group is not exist or you are not member of it")
	}
	tasks, err := s.TaskRepository.GetTaskList(ctx, userLogin, groupOID)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

var dateformat string = "2006-01-02T15:04:05Z"

func (s *TaskSrv) UpdateTask(ctx context.Context, taskID, userLogin string, updateTask models.CreateTask) error {
	groupOID, err := primitive.ObjectIDFromHex(updateTask.GroupID)
	if err != nil {
		return fmt.Errorf("group id: %v", err)
	}
	taskOID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return fmt.Errorf("group id: %v", err)
	}
	group, err := s.GroupRepository.GetGroupById(ctx, groupOID, userLogin)
	if err != nil {
		return errors.New("this group is not exist or you are not member of it")
	}

	if group.LeaderLogin != userLogin {
		return errors.New("you have no permissions to do this")
	}
	task, err := s.TaskRepository.GetTaskById(ctx, taskOID, groupOID)
	if err != nil {
		return errors.New("task was not found")
	}
	var startTime time.Time
	if !updateTask.StartTime.IsFullEmpty() {
		dateTimeStr := fmt.Sprintf("%sT%s:00Z", updateTask.StartTime.StartDate, updateTask.StartTime.StartTime)
		startTime, err = time.Parse(dateformat, dateTimeStr)
		if err != nil {
			return fmt.Errorf("invalid date or time format: %v", err)
		}
		now := time.Now().UTC()
		if startTime.Before(now) {
			return errors.New("you cant add tasks to past time")
		}
	}else{
		startTime = task.StartTime
	}
	var endTime time.Time
	var totalDuration int
	if !updateTask.StartTime.IsFullEmpty() && !updateTask.Duration.IsEmpty() {
		totalDuration = calculateDuration(updateTask.Duration)
		endTime = startTime.Add(time.Duration(totalDuration) * time.Minute)
	} else if updateTask.StartTime.IsFullEmpty() && !updateTask.Duration.IsEmpty() {
		totalDuration = calculateDuration(updateTask.Duration) 
		endTime = task.StartTime.Add(time.Duration(totalDuration) * time.Minute)
	} else if !updateTask.StartTime.IsFullEmpty() && updateTask.Duration.IsEmpty() {
		totalDuration = task.Duration
		endTime = startTime.Add(time.Duration(totalDuration) * time.Minute)
		
	}
	fmt.Println(totalDuration) 
	updates := models.Task{
		Title:     updateTask.Title,
		StartTime: startTime,
		Duration: totalDuration,
		EndTime:   endTime,
	}
	existingTasks, err := s.TaskRepository.GetTaskList(ctx, userLogin, groupOID)
	if err != nil {
		return err
	}

	for _, existingTask := range existingTasks {
		if existingTask.ID != taskOID{
			if doTasksOverlap(existingTask, updates) {
				return fmt.Errorf("task overlaps with an existing task: %s", existingTask.Title)
			}
		}
	}
	err = s.TaskRepository.UpdateTask(ctx, taskOID, updates)
	if err != nil {
		return err
	}

	return nil
}
func calculateDuration(dur models.Duration) int {
    totalMinutes := (dur.DurDays * 24 * 60) + (dur.DurHours * 60) + dur.DurMinutes
    return totalMinutes
}
func (s *TaskSrv) DeleteTask(ctx context.Context, taskID, groupID, userLogin string) error {
	groupOID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return fmt.Errorf("group id: %v", err)
	}
	taskOID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return fmt.Errorf("group id: %v", err)
	}
	group, err := s.GroupRepository.GetGroupById(ctx, groupOID, userLogin)
	if err != nil {
		return errors.New("this group is not exist or you are not member of it")
	}
	if group.LeaderLogin != userLogin {
		return errors.New("you have no permissions to do this")
	}
	err = s.TaskRepository.DeleteTask(ctx, taskOID)
	if err != nil {
		return err
	}
	return nil
}
