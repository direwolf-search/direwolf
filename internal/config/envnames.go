package config

import (
	"fmt"
	"strings"
)

const (
	mainPrefix = "DIREWOLF"
)

type environmentEnvName struct {
	prefix           string
	environmentName  string
	concreteNodeName string
	credentialsKey   string
}

func NewEnvironmentEnvName(fields ...string) *environmentEnvName {
	return &environmentEnvName{
		prefix:           mainPrefix,
		environmentName:  fields[0],
		concreteNodeName: fields[1],
		credentialsKey:   fields[2],
	}
}

func (en *environmentEnvName) String() string {
	return fmt.Sprintf("%s_%s_%s_%s",
		en.prefix,
		strings.ToUpper(en.environmentName),
		strings.ToUpper(en.concreteNodeName),
		strings.ToUpper(en.credentialsKey),
	)
}

type serviceEnvName struct {
	prefix      string
	serviceName string
	nodeName    string
	nodeName2   string
}

func NewServiceEnvName(fields ...string) *serviceEnvName {
	return &serviceEnvName{
		prefix:      mainPrefix,
		serviceName: fields[0],
		nodeName:    fields[1],
		nodeName2:   fields[2],
	}
}

func (en *serviceEnvName) String() string {
	return fmt.Sprintf("%s_%s_%s_%s",
		en.prefix,
		strings.ToUpper(en.serviceName),
		strings.ToUpper(strings.TrimSuffix(en.nodeName, "s")),
		strings.ToUpper(en.nodeName2),
	)
}
