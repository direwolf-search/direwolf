package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

/*
 https://play.golang.org/p/MMSSmKsX_tB
*/

const (
	yamlConfigFilePath = "config.yaml"
)

var (
	ErrInvalidConfiguration = errors.New("app configuration is invalid")
	ErrNoService            = errors.New("no such service registered")
)

type Config struct {
	Services     []Service     `yaml:"services"`
	Environments []Environment `yaml:"environments"`
}

// NewConfig ...
func NewConfig() Config {
	f, err := os.Open(yamlConfigFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var cfg = Config{}
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}

func (c Config) getEnvironments() []Environment {
	return c.Environments
}

func (c Config) getServices() []Service {
	return c.Services
}

func (c Config) makeEnvironmentsEnvVars() map[string]string {
	var (
		resultMap = make(map[string]string)
	)

	envs := c.getEnvironments()
	for _, env := range envs {
		environmentName := env.getName()
		concreteNodes := env.getConcreteNodes()
		for _, node := range concreteNodes {
			nodeName := node.getName()
			creds := node.getCredentials()
			for credKey, credVal := range creds {
				envEnvName := NewEnvironmentEnvName(
					environmentName,
					nodeName,
					credKey,
				)
				resultMap[envEnvName.String()] = credVal
			}
		}
	}

	return resultMap
}

func (c Config) makeServicesEnvVars() map[string]string {
	var (
		resultMap = make(map[string]string)
	)

	srvs := c.getServices()
	for _, srv := range srvs {
		serviceName := srv.GetName()
		csNodes := srv.GetNodes()
		for csNodeKey, csNodeVal := range csNodes {
			for _, val := range csNodeVal.([]interface{}) {
				serviceMap := val.(map[interface{}]interface{})
				if serviceMap["default"].(int) == 1 {
					for k, v := range serviceMap {
						if k != "default" {
							name := NewServiceEnvName(serviceName, csNodeKey, k.(string))
							resultMap[name.String()] = fmt.Sprintf("%v", v)
						}
					}
				}
			}
		}
	}

	return resultMap
}

func (c Config) ReadAndExport() error {
	environmentsMap := c.makeEnvironmentsEnvVars()
	for envName, envVal := range environmentsMap {
		err := os.Setenv(envName, envVal)
		if err != nil {
			return err
		}
	}
	servicesMap := c.makeServicesEnvVars()
	for envName, envVal := range servicesMap {
		err := os.Setenv(envName, envVal)
		if err != nil {
			return err
		}
	}

	return nil
}

type Service struct {
	Name  string                 `yaml:"name"`
	Nodes map[string]interface{} `yaml:"nodes,omitempty"`
}

func (cs Service) GetName() string {
	return cs.Name
}

func (cs Service) GetNodes() map[string]interface{} {
	return cs.Nodes
}

type Environment struct {
	Name          string         `yaml:"name"`
	ConcreteNodes []ConcreteNode `yaml:"concrete_nodes"`
}

func (env Environment) getName() string {
	return env.Name
}

func (env Environment) getConcreteNodes() []ConcreteNode {
	return env.ConcreteNodes
}

type ConcreteNode struct {
	Name        string            `yaml:"name"`
	Credentials map[string]string `yaml:"credentials"`
}

func (cn ConcreteNode) getName() string {
	return cn.Name
}

func (cn ConcreteNode) getCredentials() map[string]string {
	return cn.Credentials
}
