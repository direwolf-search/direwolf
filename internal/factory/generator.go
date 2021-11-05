package factory

import (
	"reflect"

	"direwolf/internal/domain/service/crawler"
	"direwolf/internal/factory/factories/crawler_factory"
)

var (
	crawlerReflection = reflect.TypeOf((*crawler.Crawler)(nil)).Elem()
)

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (fg *Generator) NewFactory(component interface{}) AppFactory {
	if reflect.TypeOf(component).Implements(crawlerReflection) {
		return crawler_factory.NewCrawlerFactory()
	}

	return nil
}

// https://go.dev/play/p/WMgP9QgJJbQ
