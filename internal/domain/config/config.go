package config

type Config interface {
	RandomHTTPHeaderTypeName() string
	RandomDelayRangeName() string
	WorkersNum() int
	TorGate() string
}

// CI test
