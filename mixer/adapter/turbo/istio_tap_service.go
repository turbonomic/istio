package turbo

import (
	"errors"
	"fmt"

<<<<<<< HEAD
	"github.com/turbonomic/turbo-go-sdk/pkg/service"
	"github.com/turbonomic/turbo-go-sdk/pkg/probe"
	"istio.io/istio/mixer/adapter/turbo/registration"
	"istio.io/istio/mixer/adapter/turbo/action"
	"istio.io/istio/mixer/adapter/turbo/discovery"
=======
	"github.com/turbonomic/turbo-go-sdk/pkg/probe"
	"github.com/turbonomic/turbo-go-sdk/pkg/service"
	"istio.io/istio/mixer/adapter/turbo/action"
	"istio.io/istio/mixer/adapter/turbo/discovery"
	"istio.io/istio/mixer/adapter/turbo/registration"
>>>>>>> First edition of the Turbonomic Istio Mixer adapter.
)

type IstioTAPService struct {
	*service.TAPService
}

func NewKubernetesTAPService(config *Config) (*IstioTAPService, error) {
	if config == nil || config.tapSpec == nil {
		return nil, errors.New("invalid IstioTAPServiceConfig")
	}

	// Kubernetes Probe Registration Client
	registrationClient := registration.NewIstioRegistrationClient()

	// Istio Probe Discovery Client
	discoveryClient := discovery.NewIstioDiscoveryClient(config.tapSpec, config.metricHandler)

	// Kubernetes Probe Action Execution Client
	actionHandler := action.NewIstioActionHandler()

	// The KubeTurbo TAP Service that will register the kubernetes target with the
	// Turbonomic server and await for validation, discovery, action execution requests
	tapService, err :=
		service.NewTAPServiceBuilder().
			WithTurboCommunicator(config.tapSpec.TurboCommunicationConfig).
			WithTurboProbe(probe.NewProbeBuilder(config.tapSpec.TargetType, config.tapSpec.ProbeCategory).
				WithDiscoveryOptions(probe.FullRediscoveryIntervalSecondsOption(int32(config.DiscoveryIntervalSec))).
				RegisteredBy(registrationClient).
				WithActionPolicies(registrationClient).
				WithEntityMetadata(registrationClient).
				DiscoversTarget(config.tapSpec.TargetIdentifier, discoveryClient).
				ExecutesActionsBy(actionHandler)).
			Create()
	if err != nil {
		return nil, fmt.Errorf("error when creating Istio TAPService: %s", err)
	}
	return &IstioTAPService{tapService}, nil
}
