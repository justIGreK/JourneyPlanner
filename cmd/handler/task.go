package handler

import (
	"JourneyPlanner/internal/models"
	"encoding/json"
	"net/http"
	"strconv"
)

// @Summary AddTask
// @Tags Tasks
// @Description Create new task
// @Security BearerAuth
// @Produce  json
// @Param group_id query string true "Id of group"
// @Param title query string true "Task Details"
// @Param start_time query models.StartTime true "Tasks start time"
// @Param duration query models.Duration true "Tasks duration"
// @Router /tasks/add [post]
func (h *Handler) AddTask(w http.ResponseWriter, r *http.Request) {
	userLogin, ok := r.Context().Value(UserLoginKey).(string)
	if !ok{
		logs.Error("failed to get value from context")
		http.Error(w, "Forbidden", http.StatusForbidden)
        return
	}
	var durDays int
	var durHours int
	var durMinutes int
	var err error
	daysStr := r.URL.Query().Get("days")
	if daysStr != "" {
		durDays, err = strconv.Atoi(daysStr)
		if err != nil {
			logs.Infof("Error converting days: %v", err)
			http.Error(w, "Invalid days duration parameter", http.StatusBadRequest)
			return
		}
	}
	hoursStr := r.URL.Query().Get("hours")
	logs.Info(hoursStr)
	if hoursStr != "" {
		durHours, err = strconv.Atoi(hoursStr)
		if err != nil {
			logs.Infof("Error converting hours: %v", err)
			http.Error(w, "Invalid hours duration parameter", http.StatusBadRequest)
			return
		}
	}
	minutesStr := r.URL.Query().Get("minutes")
	if minutesStr != "" {
		durMinutes, err = strconv.Atoi(minutesStr)
		if err != nil {
			logs.Infof("Error converting minutes: %v", err)
			http.Error(w, "Invalid minutes duration parameter", http.StatusBadRequest)
			return
		}
	}
	taskInfo := models.CreateTask{
		GroupID: r.URL.Query().Get("group_id"),
		Title:   r.URL.Query().Get("title"),
		StartTime: models.StartTime{
			StartDate: r.URL.Query().Get("start_date"),
			StartTime: r.URL.Query().Get("start_time"),
		},
		Duration: models.Duration{
			DurDays:    durDays,
			DurHours:   durHours,
			DurMinutes: durMinutes,
		},
	}
	if taskInfo.IsEmpty() {
		http.Error(w, "Task is empty or missing required fields", http.StatusBadRequest)
		return
	}
	err = h.Task.CreateTask(r.Context(), taskInfo, userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode("Task is created")
	if err != nil {
		logs.Error("failed to encode JSON: %v", err)
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
}

// @Summary GetTasks
// @Tags Tasks
// @Description Create new task
// @Security BearerAuth
// @Produce  json
// @Param group_id query string true "Id of group"
// @Router /tasks/getlist [get]
func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	userLogin, ok := r.Context().Value(UserLoginKey).(string)
	if !ok{
		logs.Error("failed to get value from context")
		http.Error(w, "Forbidden", http.StatusForbidden)
        return
	}
	groupID := r.URL.Query().Get("group_id")
	tasks, err := h.Task.GetTaskList(r.Context(), groupID, userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(tasks) == 0 {
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode("Tasklist is empty")
		if err != nil {
			logs.Error("failed to encode JSON: %v", err)
			http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
			return
		}
		return
	}
	response := map[string]interface{}{
		"tasks": tasks,
	}
	w.Header().Set("Content-Type", "application/json")
	err =json.NewEncoder(w).Encode(response)
	if err != nil {
		logs.Error("failed to encode JSON: %v", err)
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
}

// @Summary UpdateTask
// @Tags Tasks
// @Description update existing task
// @Security BearerAuth
// @Produce  json
// @Param group_id query string true "Id of group"
// @Param task_id query string true "task id"
// @Param title query string false "Task Details"
// @Param start_time query models.StartTime false "Tasks start time"
// @Param duration query models.Duration false "Tasks duration"
// @Router /tasks/update [put]
func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	userLogin, ok := r.Context().Value(UserLoginKey).(string)
	if !ok{
		logs.Error("failed to get value from context")
		http.Error(w, "Forbidden", http.StatusForbidden)
        return
	}
	taskID := r.URL.Query().Get("task_id")
	var durDays int
	var durHours int
	var durMinutes int
	var err error
	daysStr := r.URL.Query().Get("days")
	if daysStr != "" {
		durDays, err = strconv.Atoi(daysStr)
		if err != nil {
			logs.Infof("Error converting days: %v", err)
			http.Error(w, "Invalid days duration parameter", http.StatusBadRequest)
			return
		}
	}
	hoursStr := r.URL.Query().Get("hours")
	if hoursStr != "" {
		durHours, err = strconv.Atoi(hoursStr)
		if err != nil {
			logs.Infof("Error converting hours: %v", err)
			http.Error(w, "Invalid hours duration parameter", http.StatusBadRequest)
			return
		}
	}
	minutesStr := r.URL.Query().Get("minutes")
	if minutesStr != "" {
		durMinutes, err = strconv.Atoi(minutesStr)
		if err != nil {
			logs.Infof("Error converting minutes: %v", err)
			http.Error(w, "Invalid minutes duration parameter", http.StatusBadRequest)
			return
		}
	}
	updateTask := models.CreateTask{
		GroupID: r.URL.Query().Get("group_id"),
		Title:   r.URL.Query().Get("title"),
		StartTime: models.StartTime{
			StartDate: r.URL.Query().Get("start_date"),
			StartTime: r.URL.Query().Get("start_time"),
		},
		Duration: models.Duration{
			DurDays:    durDays,
			DurHours:   durHours,
			DurMinutes: durMinutes,
		},
	}
	if updateTask.IsEmptyUpdate() {
		http.Error(w, "No new details", http.StatusBadRequest)
		return
	}
	if updateTask.StartTime.IsMissingPart() {
		http.Error(w, "if you change start time, you need to fill and another part, and vice versa", http.StatusBadRequest)
		return
	}
	err = h.Task.UpdateTask(r.Context(), taskID, userLogin, updateTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode("Done")
	if err != nil {
		logs.Error("failed to encode JSON: %v", err)
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
}

// @Summary DeleteTask
// @Tags Tasks
// @Description Delete existing task
// @Security BearerAuth
// @Produce  json
// @Param group_id query string true "Id of group"
// @Param task_id query string true "Id of group"
// @Router /tasks/delete [delete]
func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	userLogin, ok := r.Context().Value(UserLoginKey).(string)
	if !ok{
		logs.Error("failed to get value from context")
		http.Error(w, "Forbidden", http.StatusForbidden)
        return
	}
	groupID := r.URL.Query().Get("group_id")
	taskID := r.URL.Query().Get("task_id")
	err := h.Task.DeleteTask(r.Context(), taskID, groupID, userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode("Done")
	if err != nil {
		logs.Error("failed to encode JSON: %v", err)
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
}
