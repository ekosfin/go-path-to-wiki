package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"
)

/*
//Create a bot linked to your account at: https://www.mediawiki.org/wiki/Special:BotPasswords
const (
	USERNAME = "insert username here"
	PASSWORD = "insert bot password here"
)
*/

func main() {
	go http.ListenAndServe(":8080", nil)
	startTime := time.Now()
	start := "Water"
	end := "Apple"
	s := newServer(start, end)
	go s.run()
	search := newSearch(s.requestChannel, s.linkChannel)
	fmt.Printf("Finding path from %s to %s \n", start, end)
	go search.run()
	go s.match()
	go s.startUp()
	s.finish()
	duration := time.Since(startTime)
	fmt.Printf("Execution time was: %f seconds\n", duration.Seconds())
}
