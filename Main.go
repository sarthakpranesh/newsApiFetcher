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

/*
fetchData() - used to retrieve news data from newsApi.org each time
the server is started or once each 30 minutes.
*/
func fetchData() {
	var url = getEnvVariable("URL")

	// generates a request object with request method, url and body
	req, err := http.NewRequest(http.MethodGet, url, nil)

	// defensive programing to make sure no error occurred
	if err != nil {
		log.Fatal(err)
	}

	// making a client object that will be used to request
	// with timeout set to 2s
	client := http.Client{
		Timeout: 2*time.Second,
	}

	// client performs the request using the client.Do method
	res, getErr := client.Do(req)

	// defensive programing to make sure no error occurred
	if getErr != nil {
		log.Fatal(err)
	}

	// using IO utils to ready body out of the response received from the api
	body, readErr := ioutil.ReadAll(res.Body)

	// defensive programing to make sure there was no error
	if readErr != nil {
		log.Fatal(readErr)
	}

	// making the API request format object to decode the JSON into a object
	ARF := &data.ApiRequestFmt{}

	// using unmarshal to decode the body into the request object
	jsonErr := json.Unmarshal(body, &ARF)

	// defensive programing to make sure no error occurred
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	// finally setting the list of articles
	data.SetArticles(ARF.AllArticles)
}

func setInterval(f func(), seconds int, async bool) chan bool {
	// How often to fire the passed in function
	// in seconds
	interval := time.Duration(seconds) * time.Second

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
	// if dot .env is not present then err will load directly from the os
	// for production it will load from the environment pre set
	err := godotenv.Load(".env")
	if err != nil {
		return os.Getenv(key)
	}
	return os.Getenv(key)
}

func main() {
	fmt.Println("Trying to start the server")

	l := log.New(os.Stdout, "Go-Server -> ", log.LstdFlags)

	// creating the handler objects
	sr := handlers.NewRunning(l)
	ah := handlers.NewArticles(l)

	// route handlers, ServeMux
	sm := http.NewServeMux()
	sm.Handle("/", sr)
	sm.Handle("/articles", ah)
	
	// setting up the server object
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

	/*
		go func - is a go routine that are functions or methods that run concurrently with other functions and methods
		they can be thought of lite weight threads, where as the cost of creating a go routine is
		lighter then creating a thread
	*/

	// a go routine to start the server to ListenAndServe
	go func () {
		fmt.Println("Starting Server")
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	// setting the value of the articles to set list of news articles when the server starts
	go func () {
		fetchData()
	}()

	// refreshes the list of articles in each 30 minutes
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