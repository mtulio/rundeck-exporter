package rclient

import "github.com/mtulio/go-rundeck/pkg/rundeck"

// RConf keeps config
type RConf struct {
	Base       *rundeck.ClientConfig
	EnableAPI  bool
	EnableHTTP bool
	verifySSL  bool
}

// NewConfig Returns configurator
func NewConfig() (*RConf, error) {
	return &RConf{
		Base:      &rundeck.ClientConfig{},
		verifySSL: true,
	}, nil
}

// DisableVerifySSL disables the HTTP SSL check attribute.
func (rc *RConf) DisableVerifySSL() bool {
	rc.verifySSL = false
	return true
}
