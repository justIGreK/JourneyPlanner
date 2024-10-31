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
	err := h.Group.CreateGroup(r.Context(), groupInfo.Name, userLogin, groupInfo.Invitations)
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
// @Produce  json
// @Router /groups/getlist [get]
func (h *Handler) GetGroups(w http.ResponseWriter, r *http.Request) {
	userLogin := r.Context().Value(UserLoginKey).(string)
	groups, err := h.Group.GetGroupList(r.Context(), userLogin)
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
// @Produce  json
// @Param group_id query string true "id of group"
// @Router /groups/getgroupinfo [get]
func (h *Handler) GetGroupInfo(w http.ResponseWriter, r *http.Request) {
	userLogin := r.Context().Value(UserLoginKey).(string)
	groupId := r.URL.Query().Get("group_id")
	groupDetails, err := h.Group.GetGroupByID(r.Context(), groupId, userLogin)
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

// @Summary Get group blacklist
// @Tags blacklist
// @Description Get blacklist of group
// @Security BearerAuth
// @Produce  json
// @Param group_id query string true "id of group"
// @Router /groups/blacklist [get]
func (h *Handler) GetBlacklist(w http.ResponseWriter, r *http.Request) {
	userLogin := r.Context().Value(UserLoginKey).(string)
	groupId := r.URL.Query().Get("group_id")
	blacklist, err := h.Group.GetBlacklist(r.Context(), groupId, userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := map[string]interface{}{
		"blacklist": blacklist,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary BanMember
// @Tags blacklist
// @Description Kick and ban member from group
// @Security BearerAuth
// @Produce  json
// @Param group_id query string true "id of group"
// @Param memberLogin query string true "login of member"
// @Router /groups/ban [put]
func (h *Handler) BanMember(w http.ResponseWriter, r *http.Request) {
	userLogin := r.Context().Value(UserLoginKey).(string)
	groupId := r.URL.Query().Get("group_id")
	member := r.URL.Query().Get("memberLogin")
	if member == userLogin {
		http.Error(w, "you cant kick yourself", http.StatusBadRequest)
		return
	}
	err := h.Group.BanMember(r.Context(), groupId, member, userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode("Done")
}

// @Summary UnbanMember
// @Tags blacklist
// @Description Unban member in group
// @Security BearerAuth
// @Produce  json
// @Param group_id query string true "id of group"
// @Param memberLogin query string true "login of member"
// @Router /groups/unban [put]
func (h *Handler) UnbanMember(w http.ResponseWriter, r *http.Request) {
	userLogin := r.Context().Value(UserLoginKey).(string)
	groupId := r.URL.Query().Get("group_id")
	member := r.URL.Query().Get("memberLogin")
	err := h.Group.UnbanMember(r.Context(), groupId, member, userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode("Done")
}

// @Summary LeaveFromGroup
// @Tags groups
// @Description Leave from group
// @Security BearerAuth
// @Produce  json
// @Param group_id query string true "id of group"
// @Router /groups/leaveGroup [post]
func (h *Handler) LeaveFromGroup(w http.ResponseWriter, r *http.Request) {
	userLogin := r.Context().Value(UserLoginKey).(string)
	groupId := r.URL.Query().Get("group_id")
	err := h.Group.LeaveGroup(r.Context(), groupId, userLogin)
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
// @Produce  json
// @Param user_login query string true "member login"
// @Param group_id query string true "id of group"
// @Router /groups/givelead [put]
func (h *Handler) ChangeLeader(w http.ResponseWriter, r *http.Request) {
	userLogin := r.Context().Value(UserLoginKey).(string)
	groupId := r.URL.Query().Get("group_id")
	memberLogin := r.URL.Query().Get("user_login")
	err := h.Group.GiveLeaderRole(r.Context(), groupId, userLogin, memberLogin)
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
// @Produce  json
// @Param group_id query string true "id of group"
// @Router /groups/delete [delete]
func (h *Handler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	userLogin := r.Context().Value(UserLoginKey).(string)
	groupId := r.URL.Query().Get("group_id")
	err := h.Group.DeleteGroup(r.Context(), groupId, userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

// @Summary Invite user to group
// @Tags invites
// @Description Invite user to group
// @Security BearerAuth
// @Produce  json
// @Param group_id query string true "id of group"
// @Param user_login query string true "invited user"
// @Router /groups/invite [post]
func (h *Handler) Invite(w http.ResponseWriter, r *http.Request) {
	userLogin := r.Context().Value(UserLoginKey).(string)
	inviteDetails := models.CreateInvite{
		GroupID: r.URL.Query().Get("group_id"),
		User:    r.URL.Query().Get("user_login"),
	}
	if err := validate.Struct(inviteDetails); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	err := h.Group.InviteUser(r.Context(), inviteDetails.GroupID, userLogin, inviteDetails.User)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Done")
}

// @Summary Get invite list
// @Tags invites
// @Description Get your list of invites
// @Security BearerAuth
// @Produce  json
// @Router /groups/invitelist [get]
func (h *Handler) GetInviteList(w http.ResponseWriter, r *http.Request) {
	userLogin := r.Context().Value(UserLoginKey).(string)
	invites, err := h.Group.GetInviteList(r.Context(), userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(invites) == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("Your current invitelist is empty")
		return
	}
	response := map[string]interface{}{
		"invites": invites,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

// @Summary Decline invite
// @Tags invites
// @Description Decline invite
// @Security BearerAuth
// @Produce  json
// @Param invite_id query string true "Id of invite"
// @Router /groups/declineinvite [post]
func (h *Handler) DeclineInvite(w http.ResponseWriter, r *http.Request) {
	userLogin := r.Context().Value(UserLoginKey).(string)
	inviteID := r.URL.Query().Get("invite_id")
	err := h.Group.DeclineInvite(r.Context(), userLogin, inviteID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("invite declined")

}

func (h *Handler) JoinGroup(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	err := h.Group.JoinGroup(r.Context(), token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Done")
}
