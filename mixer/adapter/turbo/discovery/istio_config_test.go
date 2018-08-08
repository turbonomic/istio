package discovery

import (
	"testing"
	"fmt"
	"github.com/magiconair/properties/assert"
	"os"
)

func TestBadInput(t *testing.T) {
	svc, err := ParseTurboCommunicationConfig("hello")
	if svc != nil {
		t.Errorf("There should not have been a service")
	}
	if err == nil {
		t.Errorf("There should have been an error")
	}
}

func TestBadServer(t *testing.T) {
	cfgMap := fmt.Sprintf(`{"communicationConfig": {
                                        "serverMeta": {"version": "1", "turboServer": "local"}, 
                                        "restAPIConfig": {"opsManagerUserName": "user", "opsManagerPassword": "pswd"}
                                   }}`)
	svc, err := ParseTurboCommunicationConfig(cfgMap)
	if svc != nil || err == nil {
		t.Errorf("Should have been an error due to a wrong input")
	}
}

func TestGoodInput(t *testing.T) {
	cfgMap := fmt.Sprintf(`{"communicationConfig": {
                                        "serverMeta": {"version": "1", "turboServer": "https://localhost"}, 
                                        "restAPIConfig": {"opsManagerUserName": "user", "opsManagerPassword": "pswd"}
                                   }}`)
	svc, err := ParseTurboCommunicationConfig(cfgMap)
	if svc == nil || err != nil {
		t.Errorf("Should have been an error due to a wrong input")
	}
	// Operations Manager part.
	assert.Equal(t, svc.OpsManagerUsername, "user")
	assert.Equal(t, svc.OpsManagerPassword, "pswd")
	assert.Equal(t, svc.TurboServer, "https://localhost")
	assert.Equal(t, svc.Version, "1")
	// Target config
	assert.Equal(t, svc.TargetUsername, defaultUsername)
	assert.Equal(t, svc.TargetPassword, defaultPassword)
	assert.Equal(t, svc.ProbeCategory, ProbeCategory)
	assert.Equal(t, svc.TargetType, TargetType)
	host, errHost := os.Hostname()
	if errHost != nil {
		t.Errorf("Error obtaining the host name")
	}
	targetId := TargetType + "-ncm-" + host
	assert.Equal(t, svc.TargetIdentifier, targetId)
}
