package ce

import (
	"context"
	"errors"
	"log"
	"net/http"
	neturl "net/url"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	"github.com/gocolly/colly/v2/proxy"

	"direwolf/internal/domain/model/host"
	"direwolf/internal/domain/model/link"
	parser "direwolf/internal/pkg/crawler/html_parser"
	rd "direwolf/internal/pkg/crawler/random_delay"
	rh "direwolf/internal/pkg/crawler/random_headers"
	"direwolf/internal/pkg/generic"
	"direwolf/internal/pkg/helpers"
	"direwolf/internal/pkg/links"
)

var (
	ErrWorkersNum      = errors.New("error of workersNum setup")
	ErrVisitUrl        = errors.New("error of url visiting")
	ErrCheckUrlVisited = errors.New("error of check is url being visited")
)

// CollyConfig implements domain/config.Config interface
type CollyConfig interface {
	RandomHTTPHeaderTypeName() string
	RandomDelayRangeName() string
	WorkersNum() int
	TorGate() string
}

type CollyEngine struct {
	queue                     *Queue
	repo                      generic.CommonRepository
	engine                    *colly.Collector
	htmlParser                parser.HTMLParser
	workersNum                int
	crawlerRules              map[string]*colly.LimitRule
	isTor                     bool
	torGate                   string
	randomHTTPHeaderGenerator *rh.RandomHTTPHeaderGenerator
	randomDelayGenerator      *rd.RandomDelayGenerator
	sync.RWMutex
}

func NewCollyEngine(isTor bool, parser parser.HTMLParser, config CollyConfig) *CollyEngine {
	var (
		torLimitRule = &colly.LimitRule{
			DomainRegexp: links.GetOnionV3URLPatternString(),
		}
		clearNetLimitRule = &colly.LimitRule{}
	)

	collyCollector := colly.NewCollector(
		colly.Debugger(&debug.LogDebugger{}),
		colly.Async(true),
		colly.URLFilters(links.GetOnionV3URLPattern()),
	)

	return &CollyEngine{
		engine:     collyCollector,
		htmlParser: parser,
		workersNum: config.WorkersNum(),
		crawlerRules: map[string]*colly.LimitRule{
			"torLimits":      torLimitRule,
			"clearNetLimits": clearNetLimitRule,
		},
		isTor:                     isTor,
		torGate:                   config.TorGate(),
		randomHTTPHeaderGenerator: rh.NewRandomHTTPHeaderGenerator(config.RandomHTTPHeaderTypeName()),
		randomDelayGenerator:      rd.NewRandomDelayGenerator(config.RandomDelayRangeName()),
	}
}

func (cen *CollyEngine) GenerateRandomHeader() *http.Header {
	cen.Lock()
	defer cen.Unlock()

	return cen.randomHTTPHeaderGenerator.GenerateRandomHTTPHeader()
}

func (cen *CollyEngine) SetRandomDelay() {
	cen.crawlerRules["torLimits"].RandomDelay = time.Duration(
		cen.randomDelayGenerator.GenerateRandomDelay(),
	) * time.Millisecond
}

func (cen *CollyEngine) SetParallelism(workersNum int) {
	cen.crawlerRules["torLimits"].Parallelism = workersNum
}

func (cen *CollyEngine) SetTorGate(gate string) {
	ps, err := proxy.RoundRobinProxySwitcher(gate)
	if err != nil {
		log.Fatal(err)
	}

	cen.engine.SetProxyFunc(ps)
}

func (cen *CollyEngine) SetHTMLParser(p interface{}) {
	if v, ok := p.(parser.HTMLParser); ok {
		cen.htmlParser = v
	}
}

func (cen *CollyEngine) Visit(
	ctx context.Context,
	url string,
) {
	var (
		addRequestToQueue = func(ctx context.Context, someUrl string, arg ...interface{}) error {
			u, err := neturl.Parse(someUrl)
			if err != nil {
				return err
			}
			req := &colly.Request{
				URL:     u,
				Headers: cen.GenerateRandomHeader(),
				Ctx:     colly.NewContext(),
				Depth:   0,
				Method:  "GET",
			}
			if len(arg) == 1 {
				err = cen.queue.AddRequest(ctx, arg[0], req, cen.repo.Insert)
				if err != nil {
					return err
				}
			}
			if len(arg) == 0 {
				err = cen.queue.AddRequest(ctx, nil, req, cen.repo.Insert)
				if err != nil {
					return err
				}
			}
			return nil
		}
	)
	// Set error handler
	cen.engine.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	cen.engine.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL)
	})

	// get full HTML page and parse it to host.Host{}
	cen.engine.OnHTML("html", func(e *colly.HTMLElement) {
		cen.Lock()
		h, err := cen.htmlParser.ParseHTML(url, e.Response.Body)
		cen.Unlock()

		if err != nil {
			log.Println(err)
		}

		//err = cen.repo.Insert(ctx, h)
		err = cen.SaveHost(ctx, h)
		if err != nil {
			log.Println(err)
		}

	})

	// On every <a> element which has href attribute call callback
	cen.engine.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		l := &link.Link{
			From:    url,
			Body:    href,
			Snippet: e.Text,
			IsV3:    true,
		}
		err := cen.SaveLink(ctx, l)
		if err != nil {
			log.Println(err)
		}
		ok, err := cen.engine.HasVisited(e.Request.AbsoluteURL(href))
		if err != nil {
			log.Println(err)
		}

		if !ok {
			err = addRequestToQueue(ctx, e.Request.AbsoluteURL(href), l)
			if err != nil {
				log.Println(err)
			}
		}
	})

	ok, err := cen.engine.HasVisited(url)
	if err != nil {
		log.Println(err)
	}

	if !ok {
		err = addRequestToQueue(ctx, url)
		if err != nil {
			log.Println(err)
		}
	}

	err = cen.queue.Run(cen.engine)
	if err != nil {
		log.Println(err)
		return
	}

	cen.engine.Wait()
}

func (cen *CollyEngine) VisitAll(ctx context.Context, urls ...string) {
	for _, u := range urls {
		cen.Visit(ctx, u)
	}
}

func (cen *CollyEngine) Init() {
	if cen.isTor {
		cen.SetTorGate(cen.torGate)
	}

	cen.SetQueue()
	cen.SetRandomDelay()
	cen.SetParallelism(cen.workersNum)
	cen.SetHTMLParser(cen.htmlParser)
}

func (cen *CollyEngine) SetQueue() {
	cen.queue = NewQueue(cen.workersNum)
}

func (cen *CollyEngine) SaveLink(ctx context.Context, l *link.Link) error {
	err := cen.repo.Insert(ctx, l)
	if err != nil {
		return err
	}

	return nil
}

func (cen *CollyEngine) SaveHost(ctx context.Context, h *host.Host) error {
	exits, err := cen.repo.Exists(ctx, h.URL)
	if err != nil {
		return err
	}

	if !exits {
		err := cen.repo.Insert(ctx, h)
		if err != nil {
			return err
		}
	} else {
		updated, err := cen.repo.Updated(ctx, h.URL, helpers.GetMd5(h.Body))
		if err != nil {
			return err
		}

		if updated {
			err := cen.repo.Update(ctx, h.URL)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (cen *CollyEngine) GetName() string {
	path := reflect.TypeOf(cen).Elem().PkgPath()
	sp := strings.Split(path, "/")

	return sp[len(sp)-1]
}

// facebookwkhpilnemxj7asaniu7vnjjbiltxjqhye3mhbshg7kx5tfyd.onion
