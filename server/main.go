package main

/*
//Create a bot linked to your account at: https://www.mediawiki.org/wiki/Special:BotPasswords
const (
	USERNAME = "insert username here"
	PASSWORD = "insert bot password here"
)
*/

func main() {

	start := "Donald Trump"
	end := "Nightwish"
	s := newServer(start, end)
	go s.run()
	search := newSearch(s.requestChannel, s.linkChannel)
	go search.run(start)
	go s.match()
	s.finish()
}
