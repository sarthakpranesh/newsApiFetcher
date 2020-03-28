package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"github.com/sarthakpranesh/newsApiFetcher/handlers"
)

func main() {
	fmt.Println("Trying to start the server")
	l := log.New(os.Stdout, "Go-Server -> ", log.LstdFlags)

	sr := handlers.NewRunning(l)

	sm := http.NewServeMux()
	sm.Handle("/", sr)
	
	
	s := &http.Server{
		Addr:              ":8080",
		Handler:           sm,
		TLSConfig:         nil,
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}

	go func () {
		fmt.Println("Starting Server")
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <- sigChan
	l.Println("Received a terminate request, graceful shutdown initiated", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}