package turbo

import (
	"errors"
	"fmt"

	"github.com/turbonomic/turbo-go-sdk/pkg/probe"
	"github.com/turbonomic/turbo-go-sdk/pkg/service"
	"istio.io/istio/mixer/adapter/turbo/action"
	"istio.io/istio/mixer/adapter/turbo/discovery"
	"istio.io/istio/mixer/adapter/turbo/registration"
)

type IstioTAPService struct {
	*service.TAPService
}

// Creates new adapter <-> Turbonomic server bridge
func NewIstioTAPService(config *Config) (*IstioTAPService, error) {
	if config == nil || config.tapSpec == nil {
		return nil, errors.New("invalid IstioTAPServiceConfig")
	}

	// Istio Probe Registration Client
	registrationClient := registration.NewIstioRegistrationClient()

	// Istio Probe Discovery Client
	discoveryClient := discovery.NewIstioDiscoveryClient(config.tapSpec, config.metricHandler)

	// Istio Probe Action Execution Client
	actionHandler := action.NewIstioActionHandler()

	// The Istio TAP Service that will register the istio mixer adapter target with the
	// Turbonomic server and await for validation, discovery, action execution requests
	tapService, err :=
		service.NewTAPServiceBuilder().
			WithTurboCommunicator(config.tapSpec.TurboCommunicationConfig).
			WithTurboProbe(probe.NewProbeBuilder(config.tapSpec.TargetType, config.tapSpec.ProbeCategory).
				WithDiscoveryOptions(probe.FullRediscoveryIntervalSecondsOption(int32(config.DiscoveryIntervalSec))).
				RegisteredBy(registrationClient).
				WithEntityMetadata(registrationClient).
				DiscoversTarget(config.tapSpec.TargetIdentifier, discoveryClient).
				ExecutesActionsBy(actionHandler)).
			Create()
	if err != nil {
		return nil, fmt.Errorf("error when creating Istio TAPService: %s", err)
	}
	return &IstioTAPService{tapService}, nil
}
