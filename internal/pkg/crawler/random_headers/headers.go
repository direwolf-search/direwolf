package random_headers

import (
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"direwolf/internal/pkg/helpers"
)

const (
	userAgentStringsFilePath = "resources/tb_user_agents.txt"
)

const (
	// DireWolfDefaultUserAgent is a UserAgent's string with self DireWolf name
	DireWolfDefaultUserAgent = "DirewolfSearchEngineBot_v0.1.0"
)

const (
	headersTypesNum = 3
	userAgentsNum   = 5
)

const (
	defaultMain headerType = iota
	defaultTbb
	defaultRandom
)

var headerTypeNames = []string{
	"defaultMain",
	"defaultTbb",
	"defaultRandom",
}

type headerType int

func NewHeaderType(headerTypeName string) headerType {
	var (
		ht headerType
	)

	for num, name := range headerTypeNames {
		if name == headerTypeName {
			ht = headerType(num)
		}
	}

	return ht
}

func (ht headerType) Int() int {
	return int(ht)
}

func (ht headerType) String() string {
	return headerTypeNames[ht.Int()]
}

// mainDefaultHeaders is an unspoofed app's headers
var mainDefaultHeaders = &Header{
	items: map[string][]string{
		"User-Agent": {
			DireWolfDefaultUserAgent,
		},
		"Accept": {
			"*/*",
		},
	},
}

// defaultTbbHeaders is a default headers for Tor Browser emulating
var defaultTbbHeaders = &Header{
	items: map[string][]string{
		"User-Agent": {
			"Mozilla/5.0 (Windows NT 10.0; rv:68.0) Gecko/20100101 Firefox/68.0",
		},
		"Accept": {
			"text/html",
			"application/xhtml+xml",
			"application/xml;q=0.9",
			"*/*;q=0.8",
		},
		"Accept-Language": {
			"en-US",
			"en;q=0.5",
		},
		"Accept-Encoding": {
			"gzip,deflate",
		},
		"Accept-Charset": {
			"utf-8",
		},
		"Connection": {
			"keep-alive",
		},
		"Upgrade-Insecure-Requests": {
			"1",
		},
	},
}

// defaultCurlHeaders is a default headers for Curl emulating
var defaultCurlHeaders = &Header{
	items: map[string][]string{
		"User-Agent": {
			"curl/7.47.0",
		},
		"Accept": {
			"*/*",
		},
	},
}

// defaultWgetHeaders is a default headers for Wget emulating
var defaultWgetHeaders = &Header{
	items: map[string][]string{
		"User-Agent": {
			"Wget/1.20.3 (linux-gnu)",
		},
		"Accept": {
			"*/*",
		},
	},
}

// tbbUserAgents is a User Agent strings for Tor Browser emulating
var tbbUserAgents = &UserAgents{
	items: []string{
		"Mozilla/5.0 (Windows NT 6.1; rv:24.0) Gecko/20100101 Firefox/24.0",
		"Mozilla/5.0 (Windows NT 10.0; rv:68.0) Gecko/20100101 Firefox/68.0",
		"Mozilla/5.0 (Windows NT 6.1; rv:52.0) Gecko/20100101 Firefox/52.0",
		"Mozilla/5.0 (Windows NT 6.1; rv:60.0) Gecko/20100101 Firefox/60.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:60.0) Gecko/20100101 Firefox/60.0",
	},
}

// curlUserAgents is a User Agent strings for Curl emulating
var curlUserAgents = &UserAgents{
	items: []string{
		"curl/7.68.0",
		"curl/7.47.0",
		"curl/7.56.1",
		"curl/7.50.2",
		"curl/7.50.1",
	},
}

// wgetUserAgents is a User Agent strings for Wget emulating
var wgetUserAgents = &UserAgents{
	items: []string{
		"Wget/1.20.3 (linux-gnu)",
		"Wget/1.19.4 (linux-gnu)",
		"Wget/1.17.1 (linux-gnu)",
		"Wget/1.16.3 (linux-gnu)",
		"Wget/1.20 (linux-gnu)",
	},
}

type Header struct {
	items map[string][]string
	sync.RWMutex
}

type HeaderItem struct {
	Key   string
	Value []string
}

func NewHeader() *Header {
	return &Header{
		items: make(map[string][]string),
	}
}

func (h *Header) Set(key string, value []string) {
	h.Lock()
	defer h.Unlock()

	h.items[key] = value
}

func (h *Header) Get(key string) ([]string, bool) {
	h.Lock()
	defer h.Unlock()

	value, ok := h.items[key]

	return value, ok
}

func (h *Header) Iter() <-chan HeaderItem {
	c := make(chan HeaderItem)

	f := func() {
		h.Lock()
		defer h.Unlock()

		for k, v := range h.items {
			c <- HeaderItem{k, v}
		}
		close(c)
	}

	go f()

	return c
}

func (h *Header) GetItems() map[string][]string {
	return h.items
}

type UserAgents struct {
	items []string
	sync.RWMutex
}

type UserAgentsItem struct {
	Index int
	Value string
}

func NewUserAgents() *UserAgents {
	return &UserAgents{
		items: make([]string, 0),
	}
}

func (ua *UserAgents) Append(item string) {
	ua.Lock()
	defer ua.Unlock()

	ua.items = append(ua.items, item)
}

func (ua *UserAgents) Iter() <-chan UserAgentsItem {
	c := make(chan UserAgentsItem)

	f := func() {
		ua.Lock()
		defer ua.Unlock()
		for index, value := range ua.items {
			c <- UserAgentsItem{index, value}
		}
		close(c)
	}
	go f()

	return c
}

func (ua *UserAgents) GetByIndex(ind int) <-chan UserAgentsItem {
	c := make(chan UserAgentsItem)

	f := func(i int) {
		ua.Lock()
		defer ua.Unlock()
		for index, value := range ua.items {
			if index == i {
				c <- UserAgentsItem{index, value}
			}
		}
		close(c)
	}
	go f(ind)

	return c
}

type RandomHTTPHeaderGenerator struct {
	sync.RWMutex
	headerTyp   headerType
	tbUAStrings []string
}

func NewRandomHTTPHeaderGenerator(headerTypeName string) *RandomHTTPHeaderGenerator {
	//	headerTypeName = os.Getenv("DW_DEFAULT_CRAWLER_RANDOM_HEADER_TYPE")
	userAgentRawStrings, err := os.ReadFile(userAgentStringsFilePath)
	if err != nil {
		log.Fatal(err)
	}

	return &RandomHTTPHeaderGenerator{
		headerTyp:   NewHeaderType(headerTypeName),
		tbUAStrings: strings.Split(string(userAgentRawStrings), "\n"),
	}
}

func (r *RandomHTTPHeaderGenerator) generateRandomHTTPHeaderForTor() *Header {
	var (
		hh = NewHeader()
	)

	headersType := helpers.RandomInt(headersTypesNum)
	time.Sleep(25 * time.Millisecond)
	userAgent := helpers.RandomInt(userAgentsNum)

	switch headersType {
	case 0:
		hh = defaultTbbHeaders
		uaCh := tbbUserAgents.GetByIndex(userAgent)
		ua := <-uaCh
		hh.Set("User-Agent", []string{ua.Value})

	case 1:
		hh = defaultCurlHeaders
		uaCh := curlUserAgents.GetByIndex(userAgent)
		ua := <-uaCh
		hh.Set("User-Agent", []string{ua.Value})

	case 2:
		hh = defaultWgetHeaders
		uaCh := wgetUserAgents.GetByIndex(userAgent)
		ua := <-uaCh
		hh.Set("User-Agent", []string{ua.Value})
	}

	return hh
}

func (r *RandomHTTPHeaderGenerator) GenerateRandomHTTPHeader() *http.Header {
	r.Lock()
	defer r.Unlock()

	var hh = NewHeader()

	switch r.headerTyp {
	case defaultMain:
		hh = mainDefaultHeaders

	case defaultTbb:
		hh = defaultTbbHeaders

	case defaultRandom:
		hh = r.generateRandomHTTPHeaderForTor()
	}

	h := http.Header(hh.items)

	return &h
}

func (r *RandomHTTPHeaderGenerator) GenerateRandomHTTPHeaderForClearNet() *Header {
	// not implemented yet
	//TODO

	return &Header{}
}

/*
h := randomHeaders()
hCh := h.Iter()
//hh := <- hCh

for n := range hCh {
fmt.Println(n.Key, n.Value)
}
*/
