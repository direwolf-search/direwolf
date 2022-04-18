package factory

import (
	"reflect"
	"strings"

	"direwolf/internal/domain"
)

// AppFactory ...
type AppFactory interface {
	BuildApp(components []interface{}) domain.App
}

func GetName(component interface{}) string {
	path := reflect.TypeOf(component).Elem().PkgPath()
	sp := strings.Split(path, "/")

	return sp[len(sp)-1]
}
