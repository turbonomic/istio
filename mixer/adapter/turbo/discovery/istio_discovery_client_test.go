package discovery

import (
	"math"
	"testing"
	"github.com/magiconair/properties/assert"
	"fmt"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"istio.io/istio/mixer/adapter/turbo/registration"
)

func TestNewIstioDiscoveryClient(t *testing.T) {
	client := NewIstioDiscoveryClient(nil, nil)
	if client == nil {
		t.Errorf("Client must not be nil")
	}
}

func TestIstioDiscoveryClient_GetAccountValues(t *testing.T) {
	cfgMap := fmt.Sprintf(`{"communicationConfig": {
                                        "serverMeta": {"version": "1", "turboServer": "https://localhost"}, 
                                        "restAPIConfig": {"opsManagerUserName": "user", "opsManagerPassword": "pswd"}
                                   }}`)
	svc, _ := ParseTurboCommunicationConfig(cfgMap, "id")
	client := NewIstioDiscoveryClient(svc, nil)
	if client == nil {
		t.Errorf("Client must not be nil")
	}
	v := client.GetAccountValues()
	if v == nil {
		t.Errorf("Values must not be nil")
	}
	assert.Equal(t, v.TargetType(), svc.TargetType)
	assert.Equal(t, v.TargetCategory(), svc.ProbeCategory)
	assert.Equal(t, v.TargetIdentifierField(), registration.TargetIdentifierField)
}

func TestIstioDiscoveryClient_Validate(t *testing.T) {
	client := NewIstioDiscoveryClient(nil, nil)
	result, err := client.Validate(nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, result, &proto.ValidationResponse{})
}

func constructMetrics(h *MetricHandler) {
	b := h.NewMetricBuilder()
	// Continue
	b.WithSource("src")
	b.WithDestination("dst")
	b.WithReceivedAmount(1)
	b.WithTransmittedAmount(1)
	b.WithFlowAmount(2)
	b.WithDuration(10)
	m, _ := b.Create()
	h.Add(m)
	b = h.NewMetricBuilder()
	b.WithSource("src")
	b.WithDestination("dst")
	b.WithReceivedAmount(2)
	b.WithTransmittedAmount(2)
	b.WithFlowAmount(4)
	b.WithDuration(10)
	m1, _ := b.Create()
	h.Add(m1)
}

func TestIstioDiscoveryClient_Discover(t *testing.T) {
	cfgMap := fmt.Sprintf(`{"communicationConfig": {
                                        "serverMeta": {"version": "1", "turboServer": "https://localhost"}, 
                                        "restAPIConfig": {"opsManagerUserName": "user", "opsManagerPassword": "pswd"}
                                   }}`)
	svc, _ := ParseTurboCommunicationConfig(cfgMap, "id")
	client := NewIstioDiscoveryClient(svc, NewMetricHandler())
	constructMetrics(client.metricHandler)
	var accountValues []*proto.AccountValue
	response, err := client.Discover(accountValues)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(response.FlowDTO), 1)
}

func TestAmount(t *testing.T) {
	amount := 10000.
	duration := 120.

	kbps1 := (amount / 1024.) / (duration / 1000.)
	kbps2 := (amount / duration) / (1024. / 1000.)
	kbps3 := (amount / duration) / (KBPS)
	assert.Equal(t, true, math.Abs(kbps1 - kbps2) < 0.0000000001)
	assert.Equal(t, kbps2, kbps3)
}