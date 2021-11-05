package torproxy

import "regexp"

const (
	torProxy = "socks5://127.0.0.1:9150"
)

func GetTorProxy() string {
	return torProxy
}

var (
	onionV3URLPattern = regexp.MustCompile(`\b(http://)?[2-7a-z]{56}\.onion\b`)
	onionV2URLPattern = regexp.MustCompile(`\b(http://)?[2-7a-z]{16}\.onion\b`)
)

var (
	onionV3URLPatternString = `\b(http://)?[2-7a-z]{56}\.onion\b`
	onionV2URLPatternString = `\b(http://)?[2-7a-z]{16}\.onion\b`
)

func GetOnionV3URLPattern() *regexp.Regexp {
	return onionV3URLPattern
}

func GetOnionV3URLPatternString() string {
	return onionV3URLPatternString
}
