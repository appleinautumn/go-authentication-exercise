package handler

import (
	"net/http"
)

type UserHandler interface {
	List(w http.ResponseWriter, r *http.Request)
}
