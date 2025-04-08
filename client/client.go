package client

type Client struct {
	ID        string
	MessageCh chan string
}

func NewClient(clientID string) *Client {
	return &Client{
		ID:        clientID,
		MessageCh: make(chan string, 100),
	}
}
