package notify

import (
	"github.com/badkaktus/gorocket"
	"net"
	"os"
	"strings"
)

type Params struct {
	Url string
	Token string
}

type Notifier struct {
	*Params
	*gorocket.Client
}

func NewNotifier(params *Params) *Notifier {
	client := gorocket.NewClient(params.Url)
	return &Notifier{params, client}
}

func (n *Notifier) Notify(text string) error {
	ipAddresses, err := localIpAddresses()
	if err != nil {
		ipAddresses = []string{"IP adresi belirlenemedi"}
	}
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "Hostname belirlenemedi"
	}
	text = "*[ " + hostname + "* | _" + strings.Join(ipAddresses, " - ") + "_ *]*\n\n" + text
	text = text + "\n\n*Komut:* `" + strings.Join(os.Args, " ") + "`"

	message := gorocket.HookMessage{
		Text: text,
	}
	_, err = n.Hooks(&message, n.Params.Token)
	return err
}

func localIpAddresses() ([]string, error) {
	ipAddresses := make([]string, 0)
	ifAddresses, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	for _, address := range ifAddresses {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipAddresses = append(ipAddresses, ipnet.IP.String())
			}
		}
	}
	return ipAddresses, nil
}
