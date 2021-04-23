package main

//This is used in the channel for sending results from search to server
type LinkMessage struct {
	ret    map[string][]string
	origin string
	depth  int
}

//This is used in the channel for sending search things to find links
type QueryMessage struct {
	pages  []string
	origin string
	depth  int
}

//This is used in the channel for checking if a path has been found
type MatchMessage struct {
	depth      int
	originPage string
}
