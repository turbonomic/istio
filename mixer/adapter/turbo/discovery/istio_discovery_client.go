package discovery

import (
	sdkprobe "github.com/turbonomic/turbo-go-sdk/pkg/probe"
	"github.com/golang/glog"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"istio.io/istio/mixer/adapter/turbo/registration"
	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
)

const (
	// KBPS
	KBPS = 1024 * 60
	// Default destination port
	DEFAULT_DESTINATION_PORT = 8080
)

// Implements the go sdk discovery client interface
type IstioDiscoveryClient struct {
	targetConfig  *IstioTAPServiceSpec
	metricHandler *MetricHandler
}

func NewIstioDiscoveryClient(tapSpec *IstioTAPServiceSpec, metricHandler *MetricHandler) *IstioDiscoveryClient {
	return &IstioDiscoveryClient{
		targetConfig:  tapSpec,
		metricHandler: metricHandler,
	}
}

// Provide the information about the target (as in how the platform talks to the target).
func (client *IstioDiscoveryClient) GetAccountValues() *sdkprobe.TurboTargetInfo {
	var accountValues []*proto.AccountValue
	conf := client.targetConfig
	// Convert all parameters in clientConf to AccountValue list
	targetID := registration.TargetIdentifierField
	accVal := &proto.AccountValue{
		Key:         &targetID,
		StringValue: &conf.TargetIdentifier,
	}
	accountValues = append(accountValues, accVal)
	username := registration.Username
	accVal = &proto.AccountValue{
		Key:         &username,
		StringValue: &conf.TargetUsername,
	}
	accountValues = append(accountValues, accVal)
	password := registration.Password
	accVal = &proto.AccountValue{
		Key:         &password,
		StringValue: &conf.TargetPassword,
	}
	accountValues = append(accountValues, accVal)
	return sdkprobe.NewTurboTargetInfoBuilder(conf.ProbeCategory,
		conf.TargetType, targetID, accountValues).Create()
}

// Validate the Target
func (client *IstioDiscoveryClient) Validate(accountValues []*proto.AccountValue) (*proto.ValidationResponse, error) {
	glog.V(2).Infof("Validating Istio target...")
	validationResponse := &proto.ValidationResponse{}
	return validationResponse, nil
}

// DiscoverTopology receives a discovery request from server and start probing the k8s.
// This is a part of the interface that gets registered with and is invoked asynchronously by the GO SDK Probe.
func (client *IstioDiscoveryClient) Discover(accountValues []*proto.AccountValue) (*proto.DiscoveryResponse, error) {
	newDiscoveryResultDTOs, err := client.doDiscover()
	if err != nil {
		glog.Errorf("Failed to obtain Istio information: %s", err)
	}
	discoveryResponse := &proto.DiscoveryResponse{
		FlowDTO: newDiscoveryResultDTOs,
	}
	return discoveryResponse, nil
}

// Perform the discovery
func (client *IstioDiscoveryClient) doDiscover() ([]*proto.FlowDTO, error) {
	var flows []*proto.FlowDTO
	metrics := client.metricHandler.GetMetrics()
	for _, m := range metrics {
		// We simulate the destination port
		flow, err := builder.NewFlowDTOBuilder().
			Source(m.src).
			Destination(m.dst, DEFAULT_DESTINATION_PORT).
			Protocol(builder.TCP).
			Received(m.rx).
			Transmitted(m.tx).
			FlowAmount((m.amount / float64(m.duration)) / KBPS).Create()
		if err != nil {
			glog.Errorf("failed to build Flow DTOs: %s", err)
			return nil, err
		}
		flows = append(flows, flow)
	}
	return flows, nil
}
