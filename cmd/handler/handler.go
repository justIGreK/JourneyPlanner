package handler

import (
	"JourneyPlanner/internal/models"
	"context"

	_ "JourneyPlanner/docs"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

type GroupService interface {
	CreateGroup(ctx context.Context, groupName, userLogin string, invites []string) error
	GetGroupList(ctx context.Context, userLogin string) ([]models.GroupList, error)
	GetGroup(ctx context.Context, groupID, userLogin string) (*models.Group, error)
	LeaveGroup(ctx context.Context, groupID, userLogin string) error
	DeleteGroup(ctx context.Context, groupID, userLogin string) error
	GiveLeaderRole(ctx context.Context, groupID, userLogin, memberLogin string) error
}
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
	GroupService
}

func NewHandler(pollService PollService, taskService TaskService,
	userService UserService, groupService GroupService) *Handler {
	return &Handler{
		PollService:  pollService,
		TaskService:  taskService,
		UserService:  userService,
		GroupService: groupService,
	}
}

func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Route("/auth", func(r chi.Router) {
		r.Post("/singUp", h.SignUp)
		r.Post("/signIn", h.SignIn)
	})
	r.Route("/groups", func(r chi.Router) {
		r.Use(h.AuthMiddleware)
		r.Post("/add", h.AddGroup)
		r.Get("/getlist", h.GetGroups)
		r.Get("/getgroupinfo", h.GetGroupInfo)
		r.Post("/leaveGroup", h.LeaveFromGroup)
		r.Put("/givelead", h.ChangeLeader)
		r.Delete("/delete", h.DeleteGroup)
	})
	return r
}
