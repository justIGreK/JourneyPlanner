package handler

import (
	"JourneyPlanner/internal/models"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

var logs *zap.SugaredLogger

func SetLogger(l *zap.Logger) {
	logs = l.Sugar()
}

// @Summary AddGroup
// @Tags groups
// @Description Create new group
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param name query string true "name of group"
// @Param invites query []string false "by adding logins you will automatically invite this users" collectionFormat(multi)
// @Router /groups/add [post]
func (h *Handler) AddGroup(w http.ResponseWriter, r *http.Request) {
	userLogin := r.Context().Value(UserLoginKey).(string)
	groupInfo := models.CreateGroup{
		Name:        r.URL.Query().Get("name"),
		Invitations: r.URL.Query()["invites"],
	}
	if err := validate.Struct(groupInfo); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	err := h.CreateGroup(r.Context(), groupInfo.Name, userLogin, groupInfo.Invitations)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Group is created")
}

// @Summary GetGroups
// @Tags groups
// @Description Get a list of all the groups you are a member of
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Router /groups/getlist [get]
func (h *Handler) GetGroups(w http.ResponseWriter, r *http.Request) {
	userLogin := r.Context().Value(UserLoginKey).(string)
	groups, err := h.GetGroupList(r.Context(), userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := map[string]interface{}{
		"discussion": groups,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary GetGroupInfo
// @Tags groups
// @Description Get full info about group you are a member of
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param group_id query string true "id of group"
// @Router /groups/getgroupinfo [get]
func (h *Handler) GetGroupInfo(w http.ResponseWriter, r *http.Request) {
	userLogin := r.Context().Value(UserLoginKey).(string)
	groupId := r.URL.Query().Get("group_id")
	groupDetails, err := h.GetGroup(r.Context(), groupId, userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := map[string]interface{}{
		"group": groupDetails,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary LeaveFromGroup
// @Tags groups
// @Description Leave from group
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param group_id query string true "id of group"
// @Router /groups/leaveGroup [post]
func (h *Handler) LeaveFromGroup(w http.ResponseWriter, r *http.Request) {
	userLogin := r.Context().Value(UserLoginKey).(string)
	groupId := r.URL.Query().Get("group_id")
	err := h.LeaveGroup(r.Context(), groupId, userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}


// @Summary GiveLeaderRole
// @Tags groups
// @Description Give another member of group leader role
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param user_login query string true "member login"
// @Param group_id query string true "id of group"
// @Router /groups/givelead [put]
func (h *Handler) ChangeLeader(w http.ResponseWriter, r *http.Request) {
	userLogin := r.Context().Value(UserLoginKey).(string)
	groupId := r.URL.Query().Get("group_id")
	memberLogin := r.URL.Query().Get("user_login")
	err := h.GiveLeaderRole(r.Context(), groupId, userLogin, memberLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}


// @Summary DeleteGroup
// @Tags groups
// @Description Delete group by id
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param group_id query string true "id of group"
// @Router /groups/delete [delete]
func (h *Handler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	userLogin := r.Context().Value(UserLoginKey).(string)
	groupId := r.URL.Query().Get("group_id")
	err := h.GroupService.DeleteGroup(r.Context(), groupId, userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
