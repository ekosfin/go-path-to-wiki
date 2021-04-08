package main

import (
	"encoding/json"
	"fmt"

	"cgt.name/pkg/go-mwclient"
)

func GetLinks(w *mwclient.Client, page string) {
	var links []string
	//Initial get
	parameters := map[string]string{
		"action":  "query",
		"prop":    "links",
		"titles":  page,
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
	//Add links
	for i := range test.Query.Pages[0].Links {
		links = append(links, test.Query.Pages[0].Links[i].Title)
	}
	//Finding all the links on the page
	for !test.Complite {
		parameters := map[string]string{
			"action":     "query",
			"prop":       "links",
			"titles":     page,
			"pllimit":    "max",
			"plcontinue": test.Continue.Plcontinue,
		}
		resp, err := w.Get(parameters)
		if err != nil {
			panic(err)
		}
		data, err := resp.Value.Marshal()
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(data, &test)
		if err != nil {
			panic(err)
		}
		//Add links
		for i := range test.Query.Pages[0].Links {
			links = append(links, test.Query.Pages[0].Links[i].Title)
		}

	}
	//Count links
	x := 0
	for range links {
		x++
	}
	fmt.Printf("There are %d, links on page %s\n", x, page)
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
