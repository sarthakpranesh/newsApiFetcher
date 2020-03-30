package handlers

import (
	"log"
	"net/http"
	"github.com/sarthakpranesh/newsApiFetcher/data"
)

// creating a struct with a logger object
// this helps in writing good tests and debugging in future
type Running struct {
	l *log.Logger
}

func NewRunning(l *log.Logger) *Running{
	return &Running{l}
}

// ServeHTTP to serve to requests when called
func (h *Running) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Println("Server alive")
	serverStatus := data.GetStatus()
	err := serverStatus.ToJSON(w)
	if err != nil {
		h.l.Println("Unable to encode data")
		http.Error(w, "Unable to encode data", http.StatusInternalServerError)
	}
}