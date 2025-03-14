package handler

import (
	"context"
	"net/http"

	"go-authentication-exercise/internal/user/entity"
	"go-authentication-exercise/internal/util"

	// "go-authentication-exercise/internal/user/request"
	"go-authentication-exercise/internal/user/service"
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
	query := r.URL.Query()

	// create paging object
	page := query.Get("page")
	limit := query.Get("limit")
	paging := util.NewPaging(page, limit)

	ctx := r.Context()
	ctx = context.WithValue(ctx, "paging", paging)

	// get users
	users, err := h.service.List(ctx)

	// get total users
	count, err := h.service.Count(ctx)
	if err != nil {
		util.Error(w, http.StatusInternalServerError, nil, err.Error())
		return
	}

	if err != nil {
		util.Error(w, http.StatusInternalServerError, nil, err.Error())
		return
	}

	paginatedData := util.Paginate(paging, listToResponse(users), count)

	util.Success(w, http.StatusOK, paginatedData, "")
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
