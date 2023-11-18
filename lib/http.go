package lib

import (
	"fmt"
	"log"

	"github.com/valyala/fasthttp"
)

type HttpClient struct {
	client  *fasthttp.Client
	baseUrl string
	auth    string
}

func NewHttpClient(baseUrl string, auth string) *HttpClient {
	return &HttpClient{client: &fasthttp.Client{}, auth: auth, baseUrl: baseUrl}
}

func (c *HttpClient) MigrateDevices(lds []DeviceLegacy) ([]Device, error) {
	ret := make([]Device, 0, len(lds))
	for _, d := range lds {
		new_d, err := d.Migrate(c)
		if err != nil {
			return nil, fmt.Errorf("error migrating device %s: %s", d.Name, err)
		}
		log.Printf("New device: %+v\n", new_d)
		if new_d.Name == "" {
			// something's wrong i don't know what
			return nil, fmt.Errorf("error migrating device %s: empty name", d.Name)
		}
		ret = append(ret, *new_d)
	}
	return ret, nil
}
