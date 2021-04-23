package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/gosuri/uilive"
)

type pageInfo struct {
	depth  int
	origin string
}

type calls struct {
	sent    int
	recived int
	canSend bool
}

type server struct {
	syncScannedPages *sync.Map
	syncCalls        *sync.Map
	scannedPages     *int
	linkChannel      chan LinkMessage
	requestChannel   chan QueryMessage
	matchChannel     chan MatchMessage
	startingPage     string
	endingPage       string
	found            bool
	foundPath        foundPath
	waitGroup        *sync.WaitGroup
	sendWaitGroup    *sync.WaitGroup
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
	var swg sync.WaitGroup
	var sM sync.Map
	var cM sync.Map
	sp := 0

	return &server{
		syncScannedPages: &sM,
		syncCalls:        &cM,
		linkChannel:      make(chan LinkMessage, 50),
		requestChannel:   make(chan QueryMessage),
		matchChannel:     make(chan MatchMessage),
		startingPage:     start,
		endingPage:       end,
		found:            false,
		foundPath:        fP,
		waitGroup:        &wg,
		sendWaitGroup:    &swg,
		scannedPages:     &sp,
	}
}

func (s *server) run() {
	for scan := range s.linkChannel {
		s.handleScan(scan)
	}
}

func (s *server) increment() {
	*s.scannedPages++
}

func (s *server) handleScan(scan LinkMessage) {
	temp, ok := s.syncCalls.Load(scan.depth)
	if ok {
		val := temp.(calls)
		val.recived = val.recived + 1
		s.syncCalls.Store(scan.depth, val)
		if val.sent == val.recived {
			val := calls{
				canSend: true,
				sent:    0,
				recived: 0,
			}
			newDepth := scan.depth + 1
			s.syncCalls.Store(newDepth, val)
		}
	} else {
		val := calls{
			canSend: true,
			sent:    1,
			recived: 1,
		}
		s.syncCalls.Store(scan.depth, val)
	}

	pI := pageInfo{
		depth:  scan.depth,
		origin: scan.origin,
	}
	for key, value := range scan.ret {
		temp, ok := s.syncScannedPages.Load(key)
		if !ok {
			s.syncScannedPages.Store(key, pI)
		} else {
			val := temp.(pageInfo)
			if pageInfo(val).depth > pI.depth {
				s.syncScannedPages.Store(key, pI)
			}
		}
		s.waitGroup.Add(1)
		newDepth := scan.depth + 1
		go s.handleLinks(value, newDepth, key)
	}
	scan.ret = nil
}

func (s *server) match() {
	for scan := range s.matchChannel {
		s.waitGroup.Add(1)
		previusPage := scan.originPage
		path := previusPage + " -> " + s.endingPage
		for previusPage != s.startingPage {
			temp, ok := s.syncScannedPages.Load(previusPage)
			if !ok {
				panic("Cannot find value")
			}
			previusPage = temp.(pageInfo).origin
			path = previusPage + " -> " + path
		}
		if scan.depth < s.foundPath.depth {
			s.foundPath.depth = scan.depth
			s.foundPath.path = path
		}
		s.waitGroup.Done()
	}
}

func (s *server) handleLinks(links []string, depth int, origin string) {
	//Readying sending package
	sending := 0
	var pages []string
	query := QueryMessage{
		pages:  pages,
		origin: origin,
		depth:  depth,
	}

	//Check for match
	for link := range links {
		if s.foundPath.depth < depth {
			break
		}
		s.increment()
		if links[link] == s.endingPage {
			//Match found sending message
			mm := MatchMessage{
				depth:      depth,
				originPage: origin,
			}
			s.waitGroup.Add(1)
			s.matchChannel <- mm
			s.waitGroup.Done()
			s.found = true

		}
	}
	//Wait for depth
	canSend := false
	for !canSend && !s.found {
		temp, ok := s.syncCalls.Load(query.depth)
		if ok {
			val := temp.(calls)
			canSend = val.canSend
		}
		time.Sleep(200)
	}
	//Send links
	for link := range links {
		//Check for duplicates
		_, ok := s.syncScannedPages.Load(link)
		if ok {
			continue
		}

		//Add to sending
		pages = append(pages, links[link])
		sending = sending + 1
		//Due to wikipedia limitations max api request for pages is 50 this ensures that it is respected
		if sending <= 49 {
			//Send request
			if !s.found {
				s.sendWaitGroup.Add(1)
				query.pages = pages
				s.sentLinks(query)
			}
			sending = 0
			pages = nil
		}

	}
	if !s.found {
		//Send remaining pages
		s.sendWaitGroup.Add(1)
		query.pages = pages
		s.sentLinks(query)
	}
	s.waitGroup.Done()
}

func (s *server) finish() {
	term := uilive.New()
	term.Start()
	for s.foundPath.depth == 9999 {
		fmt.Fprintf(term, "Searching... %d pages scanned.\n", *s.scannedPages)
		time.Sleep(time.Millisecond * 100)
	}
	term.Stop()
	close(s.requestChannel)
	s.waitGroup.Wait()
	close(s.linkChannel)
	close(s.matchChannel)
	fmt.Println("Shortest path found: ")
	fmt.Printf("Depth: %d\n", s.foundPath.depth)
	fmt.Printf("Path was: %s\n", s.foundPath.path)
}

func (s *server) sentLinks(query QueryMessage) {
	defer func() {
		if r := recover(); r != nil {
			s.sendWaitGroup.Done()
		}
	}()
	canSend := false
	for !canSend {
		temp, ok := s.syncCalls.Load(query.depth)
		if ok {
			val := temp.(calls)
			canSend = val.canSend
		}
		time.Sleep(200)
	}

	if s.foundPath.depth >= query.depth {
		temp, ok := s.syncCalls.Load(query.depth)
		if ok {
			val := temp.(calls)
			val.sent = val.sent + 1
			s.syncCalls.Store(query.depth, val)
			s.requestChannel <- query
			s.sendWaitGroup.Done()
		} else {
			s.sentLinks(query)
		}
	}
}

func (s *server) startUp() {
	var pages []string
	pages = append(pages, s.startingPage)
	query := QueryMessage{
		pages:  pages,
		origin: s.startingPage,
		depth:  0,
	}
	val := calls{
		canSend: true,
		sent:    0,
		recived: 0,
	}
	s.syncCalls.Store(0, val)
	s.sendWaitGroup.Add(1)
	go s.sentLinks(query)
}
