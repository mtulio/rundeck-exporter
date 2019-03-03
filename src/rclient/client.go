package rclient

import (
	"github.com/mtulio/go-rundeck/pkg/rundeck"
)

// RClient keeps rundeck client to iteract with API or web app.
type RClient struct {
	API  *rundeck.Client
	SOAP *rundeck.Client
}

// NewClient return the client
func NewClient(rconf *RConf) (*RClient, error) {

	rc := RClient{}
	if rconf.EnableAPI {
		cfg := &rundeck.ClientConfig{
			BaseURL:    rconf.Base.BaseURL,
			APIVersion: rconf.Base.APIVersion,
			Token:      rconf.Base.Token,
		}
		c, err := rundeck.NewClient(cfg)
		if err != nil {
			return nil, err
		}
		rc.API = c
	}
	if rconf.EnableHTTP {
		cfg := &rundeck.ClientConfig{
			BaseURL:      rconf.Base.BaseURL,
			APIVersion:   rconf.Base.APIVersion,
			AuthMethod:   rconf.Base.AuthMethod,
			Username:     rconf.Base.Username,
			Password:     rconf.Base.Password,
			OverridePath: true,
		}
		c, err := rundeck.NewClient(cfg)
		if err != nil {
			return nil, err
		}
		rc.SOAP = c
	}
	return &rc, nil
}
