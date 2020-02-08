package cosmic

import "github.com/MissionCriticalCloud/go-cosmic/v6/cosmic"

// Config is the configuration structure used to instantiate a
// new Cosmic client.
type Config struct {
	APIURL      string
	APIKey      string
	SecretKey   string
	ZoneName    string
	HTTPGETOnly bool
	Timeout     int64
}

// NewClient returns a new Cosmic client.
func (c *Config) NewClient() (*CosmicClient, error) {
	client := cosmic.NewAsyncClient(c.APIURL, c.APIKey, c.SecretKey, nil, 120)
	client.HTTPGETOnly = c.HTTPGETOnly
	client.AsyncTimeout(c.Timeout)
	return &CosmicClient{CosmicClient: client, ZoneName: c.ZoneName}, nil
}
