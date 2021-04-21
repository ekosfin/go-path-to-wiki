package main

import (
	"encoding/json"

	"cgt.name/pkg/go-mwclient"
)

type search struct {
	client         *mwclient.Client
	requestChannel chan QueryMessage
	linkChannel    chan LinkMessage
}

func newSearch(reqChan chan QueryMessage, linkChan chan LinkMessage) *search {
	c := StartClient()
	return &search{
		client:         c,
		requestChannel: reqChan,
		linkChannel:    linkChan,
	}
}

func (search *search) run() {
	for query := range search.requestChannel {
		//Process query and then wait for the next one
		search.GetLinks(query.pages, query.depth, query.origin)
	}
	//fmt.Println("Shutting down search tread...")
}

func (search *search) GetLinks(pages []string, depth int, origin string) {

	//Setting parameters for GET
	titles := ""
	for i := range pages {
		if titles == "" {
			titles = pages[i]
			continue
		}
		titles = titles + "|" + pages[i]
	}
	parameters := map[string]string{
		"action":  "query",
		"prop":    "links",
		"titles":  titles,
		"pllimit": "max",
	}

	//Initial GET
	resp, err := search.client.Get(parameters)
	if err != nil {
		panic(err)
	}
	data, err := resp.Value.Marshal()
	if err != nil {
		panic(err)
	}
	var test FullResponce
	err = json.Unmarshal(data, &test)
	if err != nil {
		panic(err)
	}

	//Creating return hashtable
	ret := make(map[string][]string)

	//Add links initial request links
	for i := range test.Query.Pages {
		links := ret[test.Query.Pages[i].Title]
		for j := range test.Query.Pages[i].Links {
			links = append(links, test.Query.Pages[i].Links[j].Title)
		}
		ret[test.Query.Pages[i].Title] = links
	}
	done := true
	//Checking if there is more data to GET from the search
	if !test.Complite {
		done = false
	}
	newparameters := map[string]string{
		"action":     "query",
		"prop":       "links",
		"titles":     titles,
		"pllimit":    "max",
		"plcontinue": test.Continue.Plcontinue,
	}

	//Finding all the links on the page
	for !done {

		resp, err := search.client.Get(newparameters)
		if err != nil {
			panic(err)
		}
		data, err := resp.Value.Marshal()
		if err != nil {
			panic(err)
		}
		var test FullResponce
		err = json.Unmarshal(data, &test)
		if err != nil {
			panic(err)
		}
		//Add links
		for i := range test.Query.Pages {
			links := ret[test.Query.Pages[i].Title]
			for j := range test.Query.Pages[i].Links {
				links = append(links, test.Query.Pages[i].Links[j].Title)
			}
			ret[test.Query.Pages[i].Title] = links
		}
		if !test.Complite {
			done = false
			newparameters["plcontinue"] = test.Continue.Plcontinue
		} else {
			done = true
		}
	}
	//Sending finds to linkChannel
	search.linkChannel <- LinkMessage{ret: ret, origin: origin, depth: depth}
}

func StartClient() *mwclient.Client {
	w, err := mwclient.New("https://en.wikipedia.org/w/api.php", "Find_links_bot")
	if err != nil {
		panic(err)
	}
	err = w.Login(USERNAME, PASSWORD)
	if err != nil {
		panic(err)
	}
	return w
}
