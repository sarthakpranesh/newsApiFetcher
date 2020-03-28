# News Api Fetcher
This simple news API that fetches the latest articles related to Covid-19 outbreak in India.
Powered by [NewsAPI.org](newsapi.org) with

## Pre-requisites
`GO` should be installed on your devices

## API End-Point
1.GET `/articles`   -   returns a list of recent articles related to Covid-19 outbreak in India

## How to start
To start the server you first need a [NewsAPI.org](newsapi.org) apiKey. 
After you have the api key follow the below steps
1. `git clone https://github.com/sarthakpranesh/newsApiFetcher.git`
2. `cd newsApiFetcher`
3. `touch .env` and add `url = "<your-api-key>"` content to the file
3. `go run Main.go` or build a container `docker build -t news-api`
