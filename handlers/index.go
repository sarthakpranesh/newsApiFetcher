package handlers

import (
	"log"
	"net/http"
	"github.com/sarthakpranesh/newsApiFetcher/data"
)

type Running struct {
	l *log.Logger
}

func NewRunning(l *log.Logger) *Running{
	return &Running{l}
}

func (h *Running) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Println("Server alive")
	serverStatus := data.GetStatus()
	err := serverStatus.ToJSON(w)
	if err != nil {
		h.l.Println("Unable to encode data")
		http.Error(w, "Unable to encode data", http.StatusInternalServerError)
	}
}