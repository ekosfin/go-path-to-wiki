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
	GetLinks(w, "Finnkino")
	GetLinks(w, "Finland")
	GetLinks(w, "COVID-19")
}
