// Package envhelper provides simple functionality for reading
// environment variables from an env file and exporting them
package envhelper

import (
	"os"
	"strings"
)

type env map[string]string

func newEnv() env {
	return make(map[string]string, 0)
}

func (e env) loadEnv() error {
	if len(e) != 0 {
		for k, v := range e {
			if err := os.Setenv(k, v); err != nil {
				return err
			}
		}
	}

	return nil
}

// getEnvContent reads content of the file with given file name
func getEnvContent(envFilePath string) ([]byte, error) {
	file, err := os.ReadFile(envFilePath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// ParseContent parses given []byte end returns env filled with keys and values from content.
// If there is no keys and values parsed, returns empty env
func parseContent(c []byte) env {
	var (
		e = newEnv()
	)

	if len(c) == 0 {
		return e
	}

	sp := strings.Split(string(c), "\n")

	for _, s := range sp {
		kv := strings.Split(s, "=")
		if len(kv) == 2 {
			key := trim(kv[0])
			val := trim(kv[1])

			e[key] = val
		}
	}

	return e
}

// LoadEnv reads file with given file name, parse its content and loads all env vars parsed
func LoadEnv(envFilePath string) error {
	c, err := getEnvContent(envFilePath)
	if err != nil {
		return err
	}

	e := parseContent(c)

	err = e.loadEnv()
	if err != nil {
		return err
	}

	return nil
}

func trim(str string) string {
	return strings.Trim(str, " ")
}
