package handlers

import (
	"github.com/sarthakpranesh/newsApiFetcher/data"
	"log"
	"net/http"
)

type Articles struct {
	l *log.Logger
}

func NewArticles(l *log.Logger) *Articles {
	return &Articles{l}
}

func (a *Articles) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		a.GetAllArticles(w, r)
		return
	}
}

func (a *Articles) GetAllArticles(w http.ResponseWriter, r *http.Request) {
	la := data.GetArticles()
	err := la.ToJSON(w)
	if err != nil {
		a.l.Fatal(err)
		http.Error(w, "Unable to Encode Articles", http.StatusInternalServerError)
	}
}
