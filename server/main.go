package main

import "fmt"

/*
//Create a bot linked to your account at: https://www.mediawiki.org/wiki/Special:BotPasswords
const (
	USERNAME = "insert username here"
	PASSWORD = "insert bot password here"
)
*/

type LinkMessage struct {
	ret map[string][]string
}

func main() {
	w := StartClient()
	var pages []string
	linkChan := make(chan LinkMessage)

	pages = append(pages, "Finnkino")
	pages = append(pages, "COVID-19")
	go GetLinks(w, pages, linkChan)

	links := <-linkChan
	for key, value := range links.ret {
		fmt.Printf("Calculating links for %s page\n", key)
		x := 0
		for range value {
			x++
		}
		fmt.Printf("There are %d links on %s page\n", x, key)
	}
}
