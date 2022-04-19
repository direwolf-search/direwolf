package html_parser

import (
	"strings"

	"github.com/antchfx/htmlquery"
	strip "github.com/grokify/html-strip-tags-go"
	"golang.org/x/net/html"

	"direwolf/internal/pkg/helpers"
	"direwolf/internal/pkg/links"
)

type parser struct{}

func NewParser() *parser {
	return &parser{}
}

func (p *parser) GetLinks(node *html.Node, url string) []map[string]interface{} {
	var ls = make([]map[string]interface{}, 0)

	hrefs := htmlquery.Find(node, "//a")

	if len(hrefs) > 0 {
		for _, a := range hrefs {
			href := htmlquery.FindOne(a, "/@href")
			if p.IsOnionLink(htmlquery.InnerText(href)) {
				l := map[string]interface{}{
					"from":    url,
					"body":    htmlquery.InnerText(href),
					"snippet": htmlquery.InnerText(a),
					"is_v3":   true,
				}

				ls = append(ls, l)
			}
		}
	}

	return ls
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

func (p *parser) ParseHTML(body []byte, url string) (map[string]interface{}, error) {
	stringReader := strings.NewReader(string(body))
	doc, err := htmlquery.Parse(stringReader)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"body":    string(body),
		"text":    p.trimTags(string(body)),
		"status":  true,
		"title":   p.getTitle(doc),
		"h1":      p.getH1(doc),
		"links":   p.GetLinks(doc, url),
		"md5hash": helpers.GetMd5(string(body)),
	}, nil
}

type HTMLParser interface {
	ParseHTML(body []byte, url string) (map[string]interface{}, error)
}

func NewHTMLParser() HTMLParser {
	return NewParser()
}
