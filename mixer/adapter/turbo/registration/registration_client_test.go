package registration

import (
	"testing"
	"github.com/magiconair/properties/assert"
)

func TestIstioRegistrationClient_GetSupplyChainDefinition(t *testing.T) {
	client := NewIstioRegistrationClient()
	chain := client.GetSupplyChainDefinition()
	assert.Equal(t, len(chain), 1)
	assert.Equal(t, len(chain[0].CommBoughtOrSet), 0)
}

func TestIstioRegistrationClient_GetIdentifyingFields(t *testing.T) {
	client := NewIstioRegistrationClient()
	assert.Equal(t, client.GetIdentifyingFields(), TargetIdentifierField)
}

func TestIstioRegistrationClient_GetAccountDefinition(t *testing.T) {
	client := NewIstioRegistrationClient()
	accounts := client.GetAccountDefinition()
	assert.Equal(t, len(accounts), 3)
	assert.Equal(t, *accounts[0].GetCustomDefinition().DisplayName, "Address")
	assert.Equal(t, *accounts[1].GetCustomDefinition().DisplayName, "Username")
	assert.Equal(t, *accounts[2].GetCustomDefinition().DisplayName, "Password")
}

func TestIstioRegistrationClient_GetActionPolicy(t *testing.T) {
	client := NewIstioRegistrationClient()
	policy := client.GetActionPolicy()
	assert.Equal(t, len(policy), 0)
}

func TestIstioRegistrationClient_GetEntityMetadata(t *testing.T) {
	client := NewIstioRegistrationClient()
	metadata := client.GetEntityMetadata()
	assert.Equal(t, len(metadata), 0)
}