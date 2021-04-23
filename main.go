package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
)

/*
//Create a bot linked to your account at: https://www.mediawiki.org/wiki/Special:BotPasswords
const (
	USERNAME = "insert username here"
	PASSWORD = "insert bot password here"
)
*/

func main() {
	args := os.Args[1:]
	go http.ListenAndServe(":8080", nil)
	start := args[0]
	end := args[1]
	s := newServer(start, end)
	go s.run()
	search := newSearch(s.requestChannel, s.linkChannel)
	fmt.Printf("Finding path from %s to %s \n", start, end)
	go search.run()
	go search.run()
	go s.match()
	go s.startUp()
	s.finish()

}
