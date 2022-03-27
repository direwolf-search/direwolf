package links

import (
	"regexp"
)

// onion v2 link: http://expyuzz4wqqyqhjn.onion/
// onion v3 link: http://2gzyxa5ihm7nsggfxnu52rck2vv4rvmdlkiu3zzui5du4xyclen53wid.onion

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

func GetOnionV2URLPattern() *regexp.Regexp {
	return onionV2URLPattern
}

func GetOnionV3URLPatternString() string {
	return onionV3URLPatternString
}
