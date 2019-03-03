package rclient

import "github.com/mtulio/go-rundeck/pkg/rundeck"

// RConf keeps config
type RConf struct {
	Base       *rundeck.ClientConfig
	EnableAPI  bool
	EnableHTTP bool
}

// NewConfig Returns configurator
func NewConfig() (*RConf, error) {
	return &RConf{
		Base: &rundeck.ClientConfig{},
	}, nil
}
