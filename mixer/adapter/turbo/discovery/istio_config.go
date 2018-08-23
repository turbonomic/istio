package discovery

import (
	"encoding/json"
	"fmt"

	"github.com/turbonomic/turbo-go-sdk/pkg/service"
)

const (
	defaultUsername = "defaultUser"
	defaultPassword = "defaultPassword"
	ProbeCategory   = "Network"
	TargetType      = "IstioMixer"
)

type IstioTargetConfig struct {
	ProbeCategory    string `json:"probeCategory,omitempty"`
	TargetType       string `json:"targetType,omitempty"`
	TargetIdentifier string `json:"targetName,omitempty"`
	TargetUsername   string `json:"-"`
	TargetPassword   string `json:"-"`
}

type IstioTAPServiceSpec struct {
	*service.TurboCommunicationConfig `json:"communicationConfig,omitempty"`
	*IstioTargetConfig                `json:"targetConfig,omitempty"`
}

func ParseTurboCommunicationConfig(configFile string, id string) (*IstioTAPServiceSpec, error) {
	// load the config
	turboCommConfig, err := readTurboCommunicationConfig(configFile)
	if err != nil {
		return nil, err
	}
	if err := turboCommConfig.ValidateTurboCommunicationConfig(); err != nil {
		return nil, err
	}
	// The target config.
	targetId := TargetType + "-ncm-" + id
	turboCommConfig.IstioTargetConfig = &IstioTargetConfig{
		TargetIdentifier: targetId,
		TargetType:       TargetType,
		ProbeCategory:    ProbeCategory,
		TargetUsername:   defaultUsername,
		TargetPassword:   defaultPassword,
	}
	return turboCommConfig, nil
}

func readTurboCommunicationConfig(path string) (*IstioTAPServiceSpec, error) {
	var spec IstioTAPServiceSpec
	err := json.Unmarshal([]byte(path), &spec)
	if err != nil {
		return nil, fmt.Errorf("parsing error :%v", err.Error())
	}
	return &spec, nil
}
