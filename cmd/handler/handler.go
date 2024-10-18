package handler

type Handler struct {
	PollService
	TaskService
	UserService
}

type PollService interface {
}

type TaskService interface {
}

type UserService interface {
}
