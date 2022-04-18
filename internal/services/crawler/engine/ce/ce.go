package ce

import (
	"context"
	"errors"
	"net/http"
	neturl "net/url"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	"github.com/gocolly/colly/v2/proxy"

	"direwolf/internal/pkg/helpers"
	"direwolf/internal/pkg/links"
	parser "direwolf/internal/services/crawler/html_parser"
	rd "direwolf/internal/services/crawler/random_delay"
	rh "direwolf/internal/services/crawler/random_headers"
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

type engineRepo interface {
	Insert(ctx context.Context, entity map[string]interface{}) error
	Updated(ctx context.Context, url, md5hash string) (bool, error)
	Exists(ctx context.Context, url string) (bool, error)
	Update(ctx context.Context, entity map[string]interface{}) error
}

type engineLogger interface {
	Printf(format string, ii ...interface{})
	Error(err error, msg string, keysAndValues ...interface{})
	Fatal(err error, msg string, keysAndValues ...interface{})
	Debug(msg string, keysAndValues ...interface{})
}

// CollyEngine is an implementation of scraping engine for crawler service.
// Built with gocolly/colly/v2 under the hood
type CollyEngine struct {
	queue                     *Queue
	repo                      engineRepo
	engine                    *colly.Collector
	logger                    engineLogger
	htmlParser                parser.HTMLParser
	workersNum                int
	crawlerRules              map[string]*colly.LimitRule
	isTor                     bool
	torGate                   string
	randomHTTPHeaderGenerator *rh.RandomHTTPHeaderGenerator
	randomDelayGenerator      *rd.RandomDelayGenerator
	sync.RWMutex
}

// NewCollyEngine returns new instance of CollyEngine
func NewCollyEngine(
	isTor bool,
	repo engineRepo,
	parser parser.HTMLParser,
	config CollyConfig,
	logger engineLogger,
) *CollyEngine {
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
		logger:     logger,
		repo:       repo,
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
		cen.logger.Fatal(err, "cannot connect to Tor network: ")
		os.Exit(1)
	}

	cen.engine.SetProxyFunc(ps)
}

func (cen *CollyEngine) SetHTMLParser(p interface{}) {
	if v, ok := p.(parser.HTMLParser); ok {
		cen.htmlParser = v
	}
}

func (cen *CollyEngine) Visit(ctx context.Context, urls ...string) {
	var (
		addRequestToQueue = func(ctx context.Context, someUrl string, arg ...map[string]interface{}) error {
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
		cen.logger.Error(err, "failed with response  ", r.Request.URL)
	})

	cen.engine.OnRequest(func(r *colly.Request) {
		cen.logger.Printf("visiting %s", r.URL)
	})

	// get full HTML page and parse it to host.Host{}
	cen.engine.OnHTML("html", func(e *colly.HTMLElement) {
		cen.Lock()

		h, err := cen.htmlParser.ParseHTML(e.Response.Body)
		if err != nil {
			cen.logger.Error(err, "cannot parse HTML ", e.Request.URL)
		}

		cen.Unlock()

		h["url"] = e.Request.AbsoluteURL(e.Request.URL.String())
		cen.logger.Debug(
			"CollyEngine: e.Request.AbsoluteURL(e.Request.URL.String()) = ",
			e.Request.AbsoluteURL(e.Request.URL.String()),
		)
		//err = cen.repo.Insert(ctx, h)
		err = cen.SaveHost(ctx, h)
		if err != nil {
			cen.logger.Error(err, "cannot save host ", h)
		}

	})

	// On every <a> element which has href attribute call callback
	cen.engine.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		l := map[string]interface{}{
			"from":    e.Request.AbsoluteURL(e.Request.URL.String()),
			"body":    href,
			"snippet": e.Text,
			"is_v3":   true,
		}

		err := cen.SaveLink(ctx, l)
		if err != nil {
			cen.logger.Error(err, "cannot save link ", l)
		}

		if !cen.hasVisited(e.Request.AbsoluteURL(href)) {
			err = addRequestToQueue(ctx, e.Request.AbsoluteURL(href), l)
			if err != nil {
				cen.logger.Error(err, "cannot add url to queue ", e.Request.AbsoluteURL(href))
			}
		}
	})

	for _, someUrl := range urls {
		if !cen.hasVisited(someUrl) {
			err := addRequestToQueue(ctx, someUrl)
			if err != nil {
				cen.logger.Error(err, "cannot add url to queue ", someUrl)
			}
		}
	}

	err := cen.queue.Run(cen.engine)
	if err != nil {
		cen.logger.Error(err, "cannot run queue ")
		return
	}

	cen.engine.Wait()
}

// hasVisited checks if url has been visited in current task.
//
// since colly.CollectorHasVisited always returns nil as the second value,
// we will not propagate this error.
// https://github.com/gocolly/colly/blob/bbf3f10c37205136e9d4f46fe8118205cc505a67/colly.go#L450
func (cen *CollyEngine) hasVisited(url string) bool {
	ok, err := cen.engine.HasVisited(url)
	if err != nil {
		cen.logger.Error(err, "cannot check if url has been visited ", url)
	}

	return ok
}

func (cen *CollyEngine) Init() { // TODO: ???????
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

func (cen *CollyEngine) SaveLink(ctx context.Context, l map[string]interface{}) error {
	err := cen.repo.Insert(ctx, l)
	if err != nil {
		return err
	}

	return nil
}

func (cen *CollyEngine) SaveHost(ctx context.Context, h map[string]interface{}) error {
	var (
		exists    bool
		err       error
		url, body string
	)

	if v, ok := h["url"]; ok {
		if s, ok := v.(string); ok {
			url = s
			exists, err = cen.repo.Exists(ctx, s)
			if err != nil {
				return err
			}
		}
	}

	if v, ok := h["body"]; ok {
		if s, ok := v.(string); ok {
			body = s
		}
	}

	if !exists {
		err := cen.repo.Insert(ctx, h)
		if err != nil {
			return err
		}
	} else {
		updated, err := cen.repo.Updated(ctx, url, helpers.GetMd5(body))
		if err != nil {
			return err
		}

		if updated {
			err := cen.repo.Update(ctx, h)
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
