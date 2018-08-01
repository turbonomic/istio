/*
	The registration client is responsible for the registering the k8s master info with the
	Turbo.
<<<<<<< HEAD
 */
package registration

import (
	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"github.com/golang/glog"
=======
*/
package registration

import (
	"github.com/golang/glog"
	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
>>>>>>> First edition of the Turbonomic Istio Mixer adapter.
)

const (
	TargetIdentifierField string = "targetIdentifier"
	Username              string = "username"
	Password              string = "password"
	propertyId            string = "id"
)

/*
	The registration client.
<<<<<<< HEAD
 */
=======
*/
>>>>>>> First edition of the Turbonomic Istio Mixer adapter.
type IstioRegistrationClient struct {
}

func NewIstioRegistrationClient() *IstioRegistrationClient {
<<<<<<< HEAD
	return &IstioRegistrationClient{
	}
=======
	return &IstioRegistrationClient{}
>>>>>>> First edition of the Turbonomic Istio Mixer adapter.
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

func (regClient *IstioRegistrationClient) GetAccountDefinition() []*proto.AccountDefEntry {
	var acctDefProps []*proto.AccountDefEntry
	// target ID
	targetIDAcctDefEntry := builder.NewAccountDefEntryBuilder(TargetIdentifierField, "Address",
		"Istio Mixer identifier", ".*", false, false).Create()
	acctDefProps = append(acctDefProps, targetIDAcctDefEntry)
	// username
	usernameAcctDefEntry := builder.NewAccountDefEntryBuilder(Username, "Username",
		"Username of the Kubernetes master", ".*", false, false).Create()
	acctDefProps = append(acctDefProps, usernameAcctDefEntry)
	// password
	passwordAcctDefEntry := builder.NewAccountDefEntryBuilder(Password, "Password",
		"Password of the Kubernetes master", ".*", false, true).Create()
	acctDefProps = append(acctDefProps, passwordAcctDefEntry)
	return acctDefProps
}

func (regClient *IstioRegistrationClient) GetIdentifyingFields() string {
	return TargetIdentifierField
}

func (regClient *IstioRegistrationClient) GetActionPolicy() []*proto.ActionPolicyDTO {
	actionPolicyBuilder := builder.NewActionPolicyBuilder()
	//1. containerPod: move, provision; not resize;
	pod := proto.EntityDTO_CONTAINER_POD
	podPolicy := make(map[proto.ActionItemDTO_ActionType]proto.ActionPolicyDTO_ActionCapability)
	podPolicy[proto.ActionItemDTO_MOVE] = proto.ActionPolicyDTO_SUPPORTED
	podPolicy[proto.ActionItemDTO_PROVISION] = proto.ActionPolicyDTO_SUPPORTED
	podPolicy[proto.ActionItemDTO_RIGHT_SIZE] = proto.ActionPolicyDTO_NOT_SUPPORTED
	regClient.addActionPolicy(actionPolicyBuilder, pod, podPolicy)
	return actionPolicyBuilder.Create()
}

func (regClient *IstioRegistrationClient) addActionPolicy(ab *builder.ActionPolicyBuilder,
	entity proto.EntityDTO_EntityType,
	policies map[proto.ActionItemDTO_ActionType]proto.ActionPolicyDTO_ActionCapability) {

	for action, policy := range policies {
		ab.WithEntityActions(entity, action, policy)
	}
}

func (regClient *IstioRegistrationClient) GetEntityMetadata() []*proto.EntityIdentityMetadata {
	result := []*proto.EntityIdentityMetadata{}
<<<<<<< HEAD
	entities := []proto.EntityDTO_EntityType{
	}
=======
	entities := []proto.EntityDTO_EntityType{}
>>>>>>> First edition of the Turbonomic Istio Mixer adapter.
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
