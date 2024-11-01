package ws

import (
	"JourneyPlanner/cmd/handler"
	"JourneyPlanner/internal/models"
	"context"
	"fmt"
	"net/http"
	"sync"
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
	Upgrader    websocket.Upgrader
	Clients     map[string]*websocket.Conn
	mu          sync.Mutex
}

func NewWebSocketHandler(chatService ChatService, groupService GroupService) *WebSocketHandler {
	return &WebSocketHandler{
		ChatService: chatService,
		Group:       groupService,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		Clients: make(map[string]*websocket.Conn),
	}
}

var logs *zap.SugaredLogger

func SetLogger(l *zap.Logger) {
	logs = l.Sugar()
}

const (
	UserLoginKey handler.ContextKey = "user_id"
)

var validate = validator.New()

func (h *WebSocketHandler) HandleConnections(w http.ResponseWriter, r *http.Request) {
	userLogin, ok := r.Context().Value(UserLoginKey).(string)
	if !ok{
		logs.Error("failed to get value from context")
		http.Error(w, "Forbidden", http.StatusForbidden)
        return
	}
	msg := models.Message{
		User:    userLogin,
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
	defer func() {
		if err := ws.Close(); err != nil {
			logs.Error("close connection error: %v", err)
		}
	}()

	connKey := fmt.Sprintf("%s_%s", msg.User, msg.GroupID)
	h.mu.Lock()
	if _, ok := h.Clients[connKey]; !ok {
		h.Clients[connKey] = ws
	}
	h.mu.Unlock()

	messages, err := h.ChatService.GetChatHistory(r.Context(), msg.GroupID)
	if err != nil {
		err = ws.WriteMessage(websocket.TextMessage, []byte("error to loading history"))
		if err !=nil{
			logs.Errorf("failed to write message: %v", err)
		}
	} else {
		for _, msg := range messages {
			message := fmt.Sprintf("%v: %v\t %v", msg.User, msg.Content, msg.Time.Format("2006-01-02 15:04:05"))
			err = ws.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil{
				logs.Errorf("failed to write message: %v", err)
			}
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
	message := fmt.Sprintf("%v: %v\t %v", msg.User, msg.Content, msg.Time.Format("2006-01-02 15:04:05"))
	h.mu.Lock()
	defer h.mu.Unlock()
	for connKey, conn := range h.Clients {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			logs.Errorf("error sending message: %v", err)
			err = conn.Close()
			if err != nil{
				logs.Errorf("failed to close connection: %v", err)
			}
			delete(h.Clients, connKey)
		}
	}
}

func (h *WebSocketHandler) NotifyUserDisconnect(userLogin, groupID string) {
	connKey := fmt.Sprintf("%s_%s", userLogin, groupID)

	h.mu.Lock()
	if conn, ok := h.Clients[connKey]; ok {
		err := conn.WriteMessage(websocket.TextMessage, []byte("You are not a member of this group anymore"))
		if err != nil{
			logs.Errorf("failed to write message: %v", err)
		}
		err = conn.Close()
		if err != nil{
			logs.Errorf("failed to close connection: %v", err)
		}
		delete(h.Clients, connKey)
	}
	h.mu.Unlock()
}
