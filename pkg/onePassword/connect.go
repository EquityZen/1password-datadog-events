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

func NewConnectAPI(authToken, host string) *ConnectAPI {
	client := connect.NewClient(host, authToken)
	return &ConnectAPI{
		CClient:   client,
		AuthToken: authToken,
		HostName:  host,
	}
}

func (c *ConnectAPI) RetrieveItemByTitle(itemUUID, vaultUUID string) (error, *onepassword.Item) {
	item, err := c.CClient.GetItem(itemUUID, vaultUUID)
	if err != nil {
		return err, nil
	}
	return nil, item
}
