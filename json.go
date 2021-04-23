package main

//This is for defining the structure of the json responces
type FullResponce struct {
	Continue Continue `json:"continue"`
	Query    Query    `json:"query"`
	Complite bool     `json:"batchcomplete"`
}

type Continue struct {
	Continue   string `json:"continue"`
	Plcontinue string `json:"plcontinue"`
}

type Query struct {
	Pages []Page `json:"pages"`
}

type Page struct {
	Links  []Link `json:"links"`
	Ns     int    `json:"ns"`
	PageId int    `json:"pageid"`
	Title  string `json:"title"`
}

type Link struct {
	Ns    int    `json:"ns"`
	Title string `json:"title"`
}
