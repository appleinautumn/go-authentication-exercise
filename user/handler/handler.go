package handler

import (
	"encoding/json"
	"net/http"

	"imp/assessment/user/entity"
	"imp/assessment/util"

	// "imp/assessment/user/request"
	"imp/assessment/user/service"
)

type userHandler struct {
	service service.UserService
}

func NewUserHandler(sv service.UserService) UserHandler {
	return &userHandler{
		service: sv,
	}
}

func (h *userHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := h.service.List(ctx)

	if err != nil {
		util.Error(w, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response, err := json.Marshal(map[string]interface{}{
		"data": listToResponse(users),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusAccepted)
	w.Write(response)
}

func listToResponse(list []*entity.User) []map[string]interface{} {
	var res []map[string]interface{}

	for _, m := range list {
		res = append(res, map[string]interface{}{
			"id":         m.Id,
			"username":   m.Username,
			"fullname":   m.Fullname,
			"created_at": m.CreatedAt,
			"updated_at": m.UpdatedAt,
		})
	}

	return res
}
