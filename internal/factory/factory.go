package factory

import (
	"reflect"
	"strings"

	"direwolf/internal/factory/app"
)

// AppFactory ...
type AppFactory interface {
	BuildApp(components []interface{}) app.App
}

func GetName(component interface{}) string {
	path := reflect.TypeOf(component).Elem().PkgPath()
	sp := strings.Split(path, "/")

	return sp[len(sp)-1]
}
