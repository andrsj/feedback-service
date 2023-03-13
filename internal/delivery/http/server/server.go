package server

import (
	"net/http"
	"time"

	"github.com/andrsj/feedback-service/internal/delivery/http/router"
)

func New(router *router.Router) *http.Server {
	//nolint
	server := &http.Server{
		Addr:              ":8080",
		Handler:           router.GetChiMux(),
		ReadHeaderTimeout: time.Second,
	}

	return server
}
