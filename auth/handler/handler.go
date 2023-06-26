package handler

import (
	"encoding/json"
	"net/http"

	"imp/assessment/auth/request"
	"imp/assessment/auth/service"
	"imp/assessment/util"
)

type authHandler struct {
	service service.AuthService
}

func NewAuthHandler(sv service.AuthService) AuthHandler {
	return &authHandler{
		service: sv,
	}
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get payload
	payload := &request.LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		util.Error(w, http.StatusBadRequest, nil, "Invalid request")
		return
	}

	// validate
	if errors := util.ValidateRequest(payload); len(errors) > 0 {
		util.Error(w, http.StatusBadRequest, errors, "Validation error")
		return
	}

	// login
	accessToken, err := h.service.Login(ctx, payload.Username, payload.Password)
	if err != nil {
		util.Error(w, http.StatusBadRequest, err, err.Error())
		return
	}

	util.Success(w, http.StatusOK, accessToken, "")
}

func (h *authHandler) Signup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get payload
	payload := &request.SignupRequest{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		util.Error(w, http.StatusBadRequest, nil, "Invalid request")
		return
	}

	// validate
	if errors := util.ValidateRequest(payload); len(errors) > 0 {
		util.Error(w, http.StatusBadRequest, errors, "Validation error")
		return
	}

	// register
	user, err := h.service.Signup(ctx, payload.Username, payload.Fullname, payload.Password)
	if err != nil {
		util.Error(w, http.StatusBadRequest, err, err.Error())
		return
	}

	util.Success(w, http.StatusOK, user, "")
}
