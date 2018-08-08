package turbo

import (
	"testing"
	"istio.io/istio/mixer/adapter/turbo/discovery"
	"github.com/magiconair/properties/assert"
	"fmt"
)

func TestNewVMTConfig(t *testing.T) {
	tapSpec := &discovery.IstioTAPServiceSpec{}
	cfg := NewVMTConfig().WithMetricHandler(discovery.NewMetricHandler()).
		WithDiscoveryInterval(100).
		WithTapSpec(tapSpec)
	assert.Equal(t, cfg.tapSpec, tapSpec)
	assert.Equal(t, cfg.DiscoveryIntervalSec, 100)
	assert.Equal(t, cfg.metricHandler, discovery.NewMetricHandler())
}

func TestNewIstioTAPService(t *testing.T) {
	cfgMap := fmt.Sprintf(`{"communicationConfig": {
                                        "serverMeta": {"version": "1", "turboServer": "https://localhost"}, 
                                        "restAPIConfig": {"opsManagerUserName": "user", "opsManagerPassword": "pswd"}
                                   }}`)
	tapSpec, err := discovery.ParseTurboCommunicationConfig(cfgMap)
	cfg := NewVMTConfig().WithMetricHandler(discovery.NewMetricHandler()).
		WithDiscoveryInterval(100).
		WithTapSpec(tapSpec)
	svc , err := NewIstioTAPService(cfg)
	assert.Equal(t, err, nil)
	if svc == nil {
		t.Errorf("Service must not be nil")
	}
}

func TestNewIstioTAPService_Error(t *testing.T) {
	svc , _ := NewIstioTAPService(nil)
	if svc != nil {
		t.Errorf("The service must be nil due to the nil config being passed in")
	}
}