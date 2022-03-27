package torproxy

const (
	torProxy = "socks5://127.0.0.1:9150"
)

func GetTorProxy() string {
	return torProxy
}
