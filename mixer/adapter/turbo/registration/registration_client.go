/*
	The registration client is responsible for the registering the k8s master info with the
	Turbo.
*/
package registration

import (
	"github.com/golang/glog"
	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

const (
	TargetIdentifierField string = "targetIdentifier"
	Username              string = "username"
	Password              string = "password"
	propertyId            string = "id"
)

/*
	The registration client.
*/
type IstioRegistrationClient struct {
}

func NewIstioRegistrationClient() *IstioRegistrationClient {
	return &IstioRegistrationClient{}
}

func (regClient *IstioRegistrationClient) GetSupplyChainDefinition() []*proto.TemplateDTO {
	supplyChainFactory := NewSupplyChainFactory()
	supplyChain, err := supplyChainFactory.createSupplyChain()
	if err != nil {
		glog.Errorf("Failed to create supply chain: %v", err)
		// TODO error handling
	}
	return supplyChain
}

// Returns the account definitions.
// The username/password fields are dummy ones at the moment.
func (regClient *IstioRegistrationClient) GetAccountDefinition() []*proto.AccountDefEntry {
	var acctDefProps []*proto.AccountDefEntry
	// target ID
	targetIDAcctDefEntry := builder.NewAccountDefEntryBuilder(TargetIdentifierField, "Address",
		"Istio Mixer identifier", ".*", false, false).Create()
	acctDefProps = append(acctDefProps, targetIDAcctDefEntry)
	// username
	usernameAcctDefEntry := builder.NewAccountDefEntryBuilder(Username, "Username",
		"Istio Mixer adapter username", ".*", false, false).Create()
	acctDefProps = append(acctDefProps, usernameAcctDefEntry)
	// password
	passwordAcctDefEntry := builder.NewAccountDefEntryBuilder(Password, "Password",
		"Istio Mixer adapter password", ".*", false, true).Create()
	acctDefProps = append(acctDefProps, passwordAcctDefEntry)
	return acctDefProps
}

// Returns the single identifying field to be used by the Turbonomic server.
func (regClient *IstioRegistrationClient) GetIdentifyingFields() string {
	return TargetIdentifierField
}

// Required by the Turbonomic GO SDK, but not needed here.
// Returning empty one.
func (regClient *IstioRegistrationClient) GetActionPolicy() []*proto.ActionPolicyDTO {
	return builder.NewActionPolicyBuilder().Create()
}

// We only need 'id' field for the entity identity metadata.
// So, return this one.
func (regClient *IstioRegistrationClient) GetEntityMetadata() []*proto.EntityIdentityMetadata {
	result := []*proto.EntityIdentityMetadata{}
	entities := []proto.EntityDTO_EntityType{}
	for _, etype := range entities {
		meta := regClient.newIdMetaData(etype, []string{propertyId})
		result = append(result, meta)
	}
	return result
}

func (regClient *IstioRegistrationClient) newIdMetaData(etype proto.EntityDTO_EntityType, names []string) *proto.EntityIdentityMetadata {
	data := make([]*proto.EntityIdentityMetadata_PropertyMetadata, 0, 100)
	for _, name := range names {
		dat := &proto.EntityIdentityMetadata_PropertyMetadata{
			Name: &name,
		}
		data = append(data, dat)
	}
	result := &proto.EntityIdentityMetadata{
		EntityType:            &etype,
		NonVolatileProperties: data,
	}
	return result
}
