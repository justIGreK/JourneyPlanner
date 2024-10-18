package handler

import (
	"JourneyPlanner/internal/models"
	"context"

	_ "JourneyPlanner/docs"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

type PollService interface {
}

type TaskService interface {
}

type UserService interface {
	LoginUser(ctx context.Context, option, password string) (string, error)
	RegisterUser(ctx context.Context, user models.SignUp) error
}

type Handler struct {
	PollService
	TaskService
	UserService
}

func NewHandler(pollService PollService, taskService TaskService, userService UserService) *Handler {
	return &Handler{
		PollService: pollService,
		TaskService: taskService,
		UserService: userService,
	}
}

func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Route("/auth", func(r chi.Router) {
		r.Post("/singUp", h.SignUp)
		r.Post("/signIn", h.SignIn)
	})
	return r
}
