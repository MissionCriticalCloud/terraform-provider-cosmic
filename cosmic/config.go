package cosmic

import "github.com/MissionCriticalCloud/go-cosmic/v6/cosmic"

// Config is the configuration structure used to instantiate a
// new Cosmic client.
type Config struct {
	APIURL      string
	APIKey      string
	SecretKey   string
	HTTPGETOnly bool
	Timeout     int64
}

// NewClient returns a new Cosmic client.
func (c *Config) NewClient() (*cosmic.CosmicClient, error) {
	cs := cosmic.NewAsyncClient(c.APIURL, c.APIKey, c.SecretKey, nil, 120)
	cs.HTTPGETOnly = c.HTTPGETOnly
	cs.AsyncTimeout(c.Timeout)
	return cs, nil
}
