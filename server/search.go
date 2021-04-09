package main

import (
	"encoding/json"

	"cgt.name/pkg/go-mwclient"
)

func GetLinks(w *mwclient.Client, pages []string, linkChan chan LinkMessage) {
	ret := make(map[string][]string)
	titles := ""
	for i := range pages {
		if titles == "" {
			titles = pages[i]
			continue
		}
		titles = titles + "|" + pages[i]
	}
	//Initial get
	parameters := map[string]string{
		"action":  "query",
		"prop":    "links",
		"titles":  titles,
		"pllimit": "max",
	}

	resp, err := w.Get(parameters)
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
	//Add links initial request links
	for i := range test.Query.Pages {
		links := ret[test.Query.Pages[i].Title]
		for j := range test.Query.Pages[i].Links {
			links = append(links, test.Query.Pages[i].Links[j].Title)
		}
		ret[test.Query.Pages[i].Title] = links
	}
	done := true
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

		resp, err := w.Get(newparameters)
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
	linkChan <- LinkMessage{ret: ret}
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
