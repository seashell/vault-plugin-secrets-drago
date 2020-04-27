package drago

import (
	"github.com/go-resty/resty/v2"
)

const (
	aclBootstrapEndpoint = "/acl/boostrap"
)

type ClientConfig struct {
	BaseUrl string
	Token   string
}

type Drago struct {
	config ClientConfig
	client *resty.Client
}

func New() (*Drago, error) {
	return &Drago{
		client: resty.New(),
	}, nil
}

func (d *Drago) DoSomething() {

	url := "http://" + d.config.BaseUrl + "..."

	body := map[string]interface{}{
		"key": "value",
	}

	resp, err := d.client.R().
		SetHeader("X-Drago-Token", d.config.Token).
		SetBody(body).
		Post(url)

}
