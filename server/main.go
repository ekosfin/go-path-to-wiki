package main

/*
//Create a bot linked to your account at: https://www.mediawiki.org/wiki/Special:BotPasswords
const (
	USERNAME = "insert username here"
	PASSWORD = "insert bot password here"
)
*/

func main() {
	w := StartClient()
	var pages []string
	pages = append(pages, "Finland")
	pages = append(pages, "Finnkino")
	pages = append(pages, "COVID-19")
	GetLinks(w, pages)
}
