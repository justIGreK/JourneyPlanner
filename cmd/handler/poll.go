package handler

import (
	"JourneyPlanner/internal/models"
	"encoding/json"
	"net/http"
	"strconv"
)

const noDuration = 15770000 // if selected no duration, poll will be active for 30 years

// @Summary CreatePoll
// @Tags polls
// @Description Create new poll
// @Security BearerAuth
// @Produce  json
// @Param groupID query string true "id of group"
// @Param title query string true "title of poll"
// @Param firstOption query string true "first option "
// @Param sercondOption query string true "second option"
// @Param duration query uint false "duration of poll in minutes" minimum(0)
// @Router /polls/add [post]
func (h *Handler) CreatePoll(w http.ResponseWriter, r *http.Request) {
	userLogin := r.Context().Value(UserLoginKey).(string)
	var duration uint64
	var err error
	durationStr := r.URL.Query().Get("duration")
	if durationStr != "" {
		duration, err = strconv.ParseUint(durationStr, 10, 32)
		if err != nil {
			http.Error(w, "InvalidTime: "+err.Error(), http.StatusBadRequest)
			return
		}
	}else {
		duration = noDuration
	}
	pollInfo := models.CreatePoll{
		GroupID:      r.URL.Query().Get("groupID"),
		Title:        r.URL.Query().Get("title"),
		FirstOption:  r.URL.Query().Get("firstOption"),
		SecondOption: r.URL.Query().Get("sercondOption"),
		Duration:     duration,
	}
	if err := validate.Struct(pollInfo); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Poll.CreatePoll(r.Context(), pollInfo, userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode("Poll is created")
}

// @Summary GetPolls
// @Tags polls
// @Description Get list of polls
// @Security BearerAuth
// @Produce  json
// @Param groupID query string true "id of group"
// @Router /polls/getlist [get]
func (h *Handler) GetPolls(w http.ResponseWriter, r *http.Request) {
	userLogin := r.Context().Value(UserLoginKey).(string)
	groupID := r.URL.Query().Get("groupID")
	polls, err := h.Poll.GetPollList(r.Context(), groupID, userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := map[string]interface{}{
		"polls": polls,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary Delete Poll
// @Tags polls
// @Description Delete poll by id
// @Security BearerAuth
// @Produce  json
// @Param groupID query string true "id of group"
// @Param pollID query string true "id of poll"
// @Router /polls/delete [delete]
func (h *Handler) DeletePoll(w http.ResponseWriter, r *http.Request) {
	userLogin := r.Context().Value(UserLoginKey).(string)
	groupID := r.URL.Query().Get("groupID")
	pollID := r.URL.Query().Get("pollID")
	err := h.Poll.DeletePollByID(r.Context(), pollID, groupID, userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Done")
}

// @Summary Close Poll
// @Tags polls
// @Description Close poll for voting
// @Security BearerAuth
// @Produce  json
// @Param groupID query string true "id of group"
// @Param pollID query string true "id of poll"
// @Router /polls/close [put]
func (h *Handler) ClosePoll(w http.ResponseWriter, r *http.Request) {
	userLogin := r.Context().Value(UserLoginKey).(string)
	groupID := r.URL.Query().Get("groupID")
	pollID := r.URL.Query().Get("pollID")
	err := h.Poll.ClosePoll(r.Context(), pollID, groupID, userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Done")
}



// @Summary Vote Poll
// @Tags polls
// @Description Vote for poll option
// @Security BearerAuth
// @Produce  json
// @Param groupID query string true "id of group"
// @Param pollID query string true "id of poll"
// @Param option query string true "vote option" Enums(firstOption, secondOption)
// @Router /polls/vote [put]
func (h *Handler) VotePoll(w http.ResponseWriter, r *http.Request) {
	userLogin := r.Context().Value(UserLoginKey).(string)
	vote := models.AddVote{
		GroupID: r.URL.Query().Get("groupID"),
		PollID: r.URL.Query().Get("pollID"),
		Option: r.URL.Query().Get("option"),
	}
	err := h.Poll.VotePoll(r.Context(), userLogin, vote)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Done")
}





