package onePassword

import (
	"github.com/1Password/connect-sdk-go/connect"
	"github.com/1Password/connect-sdk-go/onepassword"
)

type ConnectAPI struct {
	CClient   connect.Client
	AuthToken string
	HostName  string
}

type LogMessage struct {
	Name     string `json:"name"`
	URL      string `json:"URL"`
	Title    string `json:"title"`
	Action   string `json:"action"`
	Email    string `json:"email"`
	Category string `json:"category"`
}

func NewConnectAPI(authToken, host string) *ConnectAPI {
	client := connect.NewClient(host, authToken)
	return &ConnectAPI{
		CClient:   client,
		AuthToken: authToken,
		HostName:  host,
	}
}

func (c *ConnectAPI) RetrieveItemByTitle(itemUUID, vaultUUID string) (*onepassword.Item, error) {
	item, err := c.CClient.GetItem(itemUUID, vaultUUID)
	if err != nil {
		return nil, err
	}
	return item, err
}
