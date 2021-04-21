package main

import (
	"fmt"
	"sync"
	"time"
)

type pageInfo struct {
	depth  int
	origin string
}

type server struct {
	scannedPages   map[string]*pageInfo
	linkChannel    chan LinkMessage
	requestChannel chan QueryMessage
	matchChannel   chan MatchMessage
	startingPage   string
	endingPage     string
	found          bool
	foundPath      foundPath
	waitGroup      *sync.WaitGroup
}

type foundPath struct {
	depth int
	path  string
}

func newServer(start string, end string) *server {
	fP := foundPath{
		depth: 9999,
		path:  "None",
	}
	var wg sync.WaitGroup
	return &server{
		scannedPages:   make(map[string]*pageInfo),
		linkChannel:    make(chan LinkMessage),
		requestChannel: make(chan QueryMessage),
		matchChannel:   make(chan MatchMessage),
		startingPage:   start,
		endingPage:     end,
		found:          false,
		foundPath:      fP,
		waitGroup:      &wg,
	}
}

func (s *server) run() {
	for scan := range s.linkChannel {
		pI := &pageInfo{
			depth:  scan.depth,
			origin: scan.origin,
		}
		scannedPages := s.duplicateMap()
		for key, value := range scan.ret {
			val, ok := s.scannedPages[key]
			if !ok {
				s.scannedPages[key] = pI
				fmt.Printf("PAGE: %s added\n", key)
			} else {
				if val.depth > pI.depth {
					s.scannedPages[key] = pI
					fmt.Printf("PAGE: %s faster path found\n", key)
				}
			}
			s.waitGroup.Add(1)
			newDepth := scan.depth + 1
			go s.handleLinks(value, newDepth, key, scannedPages)
		}
	}
	fmt.Println("Shutting down server tread...")
}

func (s *server) match() {
	for scan := range s.matchChannel {
		previusPage := scan.originPage
		path := previusPage + " -> " + s.endingPage
		for previusPage != s.startingPage {
			previusPage = s.scannedPages[previusPage].origin
			path = previusPage + " -> " + path
		}
		fmt.Printf("One path found at depth: %d path: %s\n", scan.depth, path)
		if scan.depth < s.foundPath.depth {
			s.foundPath.depth = scan.depth
			s.foundPath.path = path
		}
	}
	fmt.Println("Shutting down server match tread...")
}

func (s *server) handleLinks(links []string, depth int, origin string, scannedPages map[string]*pageInfo) {
	//Now one step deeper
	depth += depth
	//Readying sending package
	sending := 0
	var pages []string
	//Looping links
	for link := range links {
		if links[link] == s.endingPage {
			//Match found sending message
			mm := MatchMessage{
				depth:      depth,
				originPage: origin,
			}
			s.matchChannel <- mm
			s.found = true
			break
		}
		//Check if duplicate site
		_, ok := scannedPages[links[link]]
		if ok {
			//Page exists do not send request
		} else {
			//Add to sending
			pages = append(pages, links[link])
			sending += sending
			//Due to wikipedia limitations max api request for pages is 50 this ensures that it is respected
			if sending <= 49 {
				query := QueryMessage{
					pages:  pages,
					origin: origin,
					depth:  depth,
				}
				//Send request
				if !s.found {
					s.requestChannel <- query
				}
				sending = 0
				pages = nil
			}
		}
	}
	if !s.found {
		//Send remaining pages
		query := QueryMessage{
			pages:  pages,
			origin: origin,
			depth:  depth,
		}
		s.requestChannel <- query
	}
	s.waitGroup.Done()
}

func (s *server) duplicateMap() map[string]*pageInfo {
	copy := make(map[string]*pageInfo)
	for key, value := range s.scannedPages {
		copy[key] = value
	}
	return copy
}

func (s *server) finish() {
	for s.foundPath.depth == 9999 {
		time.Sleep(500)
	}
	close(s.requestChannel)
	s.waitGroup.Wait()
	close(s.linkChannel)
	fmt.Println("Shortest path found:")
	fmt.Printf("Path was: %s\n", s.foundPath.path)
}
