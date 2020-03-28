package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"github.com/sarthakpranesh/newsApiFetcher/handlers"
	"github.com/sarthakpranesh/newsApiFetcher/data"
	"github.com/joho/godotenv"
)

func fetchData() {
	var url = getEnvVariable("URL")
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	client := http.Client{
		Timeout: 2*time.Second,
	}

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(err)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	ARF := &data.ApiRequestFmt{}
	jsonErr := json.Unmarshal(body, &ARF)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	data.SetArticles(ARF.AllArticles)
}

func setInterval(f func(), milliseconds int, async bool) chan bool {
	// How often to fire the passed in function
	// in seconds
	interval := time.Duration(milliseconds) * time.Second

	// Setup the ticket and the channel to signal
	// the ending of the interval
	ticker := time.NewTicker(interval)
	clear := make(chan bool)

	// Put the selection in a go routine
	// so that the for loop is none blocking
	go func() {
		for {

			select {
			case <-ticker.C:
				if async {
					// This won't block
					go f()
				} else {
					// This will block
					f()
				}
			case <-clear:
				ticker.Stop()
				return
			}

		}
	}()

	// We return the channel so we can pass in
	// a value to it to clear the interval
	return clear
}

func getEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		return os.Getenv(key)
	}
	return os.Getenv(key)
}

func main() {
	fmt.Println("Trying to start the server")

	l := log.New(os.Stdout, "Go-Server -> ", log.LstdFlags)

	sr := handlers.NewRunning(l)
	ah := handlers.NewArticles(l)

	sm := http.NewServeMux()
	sm.Handle("/", sr)
	sm.Handle("/articles", ah)
	
	
	s := &http.Server{
		Addr:              "0.0.0.0:"+getEnvVariable("PORT"),
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

	go func () {
		fetchData()
	}()

	go func () {
		_ = setInterval(fetchData , 1800, true)
	}()

	// makes a signal channel for terminating the server peacefully
	// without sudden death ;>
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <- sigChan
	l.Println("Received a terminate request, graceful shutdown initiated", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}