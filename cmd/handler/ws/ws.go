package ws

import (
	"JourneyPlanner/internal/models"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type GroupService interface {
	GetGroupByID(ctx context.Context, groupID, userLogin string) (*models.Group, error)
}
type ChatService interface {
	SaveMessage(ctx context.Context, msg models.Message) error
	GetChatHistory(ctx context.Context, groupID string) ([]models.Message, error)
}

type WebSocketHandler struct {
	ChatService ChatService
	Group       GroupService
	Clients     map[*websocket.Conn]bool
	Upgrader    websocket.Upgrader
}

func NewWebSocketHandler(chatService ChatService, groupService GroupService) *WebSocketHandler {
	return &WebSocketHandler{
		ChatService: chatService,
		Group:       groupService,
		Clients:     make(map[*websocket.Conn]bool),
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

var logs *zap.SugaredLogger

func SetLogger(l *zap.Logger) {
	logs = l.Sugar()
}

type contextKey string

const (
	UserLoginKey contextKey = "user_id"
)

var validate = validator.New()

func (h *WebSocketHandler) HandleConnections(w http.ResponseWriter, r *http.Request) {
	msg := models.Message{
		User:    r.URL.Query().Get("user_login"),
		GroupID: r.URL.Query().Get("group_id"),
		Time:    time.Now().UTC(),
	}
	if err := validate.Struct(msg); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	_, err := h.Group.GetGroupByID(r.Context(), msg.GroupID, msg.User)
	if err != nil {
		logs.Errorf("invalid user: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	ws, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		logs.Errorf("error during connecting:%v", err)
		return
	}
	defer ws.Close()
	h.Clients[ws] = true

	messages, err := h.ChatService.GetChatHistory(r.Context(), msg.GroupID)
	if err != nil {
		ws.WriteJSON("error to loading history")
	} else {

		for _, msg := range messages {
			message := fmt.Sprintf("%v: %v\t %v", msg.User, msg.Content, msg.Time.Format("2006-01-02 15:04:05"))
			ws.WriteMessage(websocket.TextMessage, []byte(message))
		}
	}

	for {
		_, text, err := ws.ReadMessage()
		if err != nil {
			logs.Errorf("error reading message:%v", err)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		msg.Content = string(text)

		if err := h.ChatService.SaveMessage(r.Context(), msg); err != nil {
			logs.Errorf("error reading message: %v", err)
			continue
		}
		h.broadcastMessage(msg)
	}
}

func (h *WebSocketHandler) broadcastMessage(msg models.Message) {
	for client := range h.Clients {
		message := fmt.Sprintf("%v: %v\t %v", msg.User, msg.Content, msg.Time.Format("2006-01-02 15:04:05"))
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			logs.Errorf("error sending message: %v", err)
			client.Close()
			delete(h.Clients, client)
		}
	}
}
