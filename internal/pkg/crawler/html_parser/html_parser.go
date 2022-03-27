package html_parser

import (
	"strings"

	"github.com/antchfx/htmlquery"
	strip "github.com/grokify/html-strip-tags-go"
	"golang.org/x/net/html"

	"direwolf/internal/domain/model/host"
	"direwolf/internal/domain/model/link"
	"direwolf/internal/pkg/helpers"
	"direwolf/internal/pkg/links"
)

type parser struct{}

func NewParser() *parser {
	return &parser{}
}

func (p *parser) GetLinks(node *html.Node, url string) []*link.Link {
	var links = make([]*link.Link, 0)

	hrefs := htmlquery.Find(node, "//a")

	if len(hrefs) != 0 {
		for _, a := range hrefs {
			href := htmlquery.FindOne(a, "/@href")
			if p.IsOnionLink(htmlquery.InnerText(href)) {
				l := link.NewLink(url, htmlquery.InnerText(href), htmlquery.InnerText(a), true)
				links = append(links, l)
			}
		}
	}

	return links
}

func (p *parser) IsOnionLink(href string) bool {
	return links.GetOnionV3URLPattern().MatchString(href)
}

func (p *parser) getH1(node *html.Node) string {
	return htmlquery.InnerText(htmlquery.FindOne(node, "//h1"))
}

func (p *parser) getTitle(node *html.Node) string {
	return htmlquery.InnerText(htmlquery.FindOne(node, "//title"))
}

// trimTags removes all html tags from given string
func (p *parser) trimTags(body string) string {
	return strip.StripTags(body)
}

func (p *parser) ParseHTML(url string, body []byte) (*host.Host, error) {
	var (
		h host.Host
	)

	stringReader := strings.NewReader(string(body))
	doc, err := htmlquery.Parse(stringReader)
	if err != nil {
		return nil, err
	}

	//h.Links = p.getLinks(doc, url)
	h.Body = string(body)
	h.Text = p.trimTags(string(body)) // TODO:
	h.Status = true
	h.URL = url
	h.Title = p.getTitle(doc)
	h.H1 = p.getH1(doc)
	h.Hash = helpers.GetMd5(h.Body)

	return &h, nil
}

type HTMLParser interface {
	ParseHTML(url string, body []byte) (*host.Host, error)
}

func NewHTMLParser() HTMLParser {
	return NewParser()
}
