package handler

import (
	"JourneyPlanner/internal/models"
	"JourneyPlanner/internal/service"
	"context"
	"net/http"

	_ "JourneyPlanner/docs"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

type GroupService interface {
	CreateGroup(ctx context.Context, groupName, userLogin string) error
	GetGroupList(ctx context.Context, userLogin string) ([]models.GroupList, error)
	GetGroupByID(ctx context.Context, groupID, userLogin string) (*models.Group, error)
	LeaveGroup(ctx context.Context, groupID, userLogin string) error
	DeleteGroup(ctx context.Context, groupID, userLogin string) error
	GiveLeaderRole(ctx context.Context, groupID, userLogin, memberLogin string) error
	InviteUser(ctx context.Context, groupID, userLogin, invitedUser string) error
	GetInviteList(ctx context.Context, userLogin string) ([]models.InvitationList, error)
	DeclineInvite(ctx context.Context, userLogin, inviteID string) error
	JoinGroup(ctx context.Context, token string) error
	GetBlacklist(ctx context.Context, groupID, userLogin string) (*models.BlackList, error)
	BanMember(ctx context.Context, groupID, memberLogin, userLogin string) error
	UnbanMember(ctx context.Context, groupID, memberLogin, userLogin string) error
}
type PollService interface {
	CreatePoll(ctx context.Context, pollInfo models.CreatePoll, userLogin string) error
	GetPollList(ctx context.Context, groupID, userLogin string) (*models.PollList, error)
	DeletePollByID(ctx context.Context, pollID, groupID, userLogin string) error
	ClosePoll(ctx context.Context, pollID, groupID, userLogin string) error
	VotePoll(ctx context.Context, userLogin string, vote models.AddVote) error
}

type TaskService interface {
	CreateTask(ctx context.Context, taskInfo models.CreateTask, userLogin string) error
	GetTaskList(ctx context.Context, groupID, userLogin string) ([]models.Task, error)
	UpdateTask(ctx context.Context, taskID, userLogin string, task models.CreateTask) error
	DeleteTask(ctx context.Context, taskID, groupID, userLogin string) error
}

type UserService interface {
	LoginUser(ctx context.Context, option, password string) (string, error)
	RegisterUser(ctx context.Context, user models.SignUp) error
	ValidatePasetoToken(tokenString string) (*service.TokenPayload, error)
}

type WsHandler interface {
	HandleConnections(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	Poll      PollService
	Task      TaskService
	User      UserService
	Group     GroupService
}

func NewHandler(pollService PollService, taskService TaskService,
	userService UserService, groupService GroupService) *Handler {
	return &Handler{
		Poll:      pollService,
		Task:      taskService,
		User:      userService,
		Group:     groupService,
	}
}

func (h *Handler) InitRoutes(wsHandler WsHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Get("/join-group", h.JoinGroup)
	r.Route("/auth", func(r chi.Router) {
		r.Post("/singUp", h.SignUp)
		r.Post("/signIn", h.SignIn)
	})
	r.Route("/groups", func(r chi.Router) {
		r.Use(h.AuthMiddleware)
		r.Get("/ws", wsHandler.HandleConnections)	
		r.Post("/add", h.AddGroup)
		r.Get("/getlist", h.GetGroups)
		r.Get("/getgroupinfo", h.GetGroupInfo)
		r.Post("/leaveGroup", h.LeaveFromGroup)
		r.Put("/givelead", h.ChangeLeader)
		r.Delete("/delete", h.DeleteGroup)
		r.Post("/invite", h.Invite)
		r.Get("/invitelist", h.GetInviteList)
		r.Post("/declineinvite", h.DeclineInvite)
		r.Put("/ban", h.BanMember)
		r.Put("/unban", h.UnbanMember)
		r.Get("/blacklist", h.GetBlacklist)
	})
	r.Route("/tasks", func(r chi.Router) {
		r.Use(h.AuthMiddleware)
		r.Post("/add", h.AddTask)
		r.Get("/getlist", h.GetTasks)
		r.Delete("/delete", h.DeleteTask)
		r.Put("/update", h.UpdateTask)
	})
	r.Route("/polls", func(r chi.Router) {
		r.Use(h.AuthMiddleware)
		r.Post("/add", h.CreatePoll)
		r.Get("/getlist", h.GetPolls)
		r.Delete("/delete", h.DeletePoll)
		r.Put("/close", h.ClosePoll)
		r.Put("/vote", h.VotePoll)
	})
	return r
}
