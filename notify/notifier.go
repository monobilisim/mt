package notify

import "github.com/badkaktus/gorocket"

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
	message := gorocket.HookMessage{Text: text}
	_, err := n.Hooks(&message, n.Params.Token)
	return err
}
