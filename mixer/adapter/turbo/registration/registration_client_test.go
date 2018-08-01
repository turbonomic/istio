package registration

import (
	"fmt"
	"github.com/turbonomic/kubeturbo/pkg/discovery/stitching"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"testing"
)

func xcheck(expected map[proto.ActionItemDTO_ActionType]proto.ActionPolicyDTO_ActionCapability,
	elements []*proto.ActionPolicyDTO_ActionPolicyElement) error {

	if len(expected) != len(elements) {
		return fmt.Errorf("length not equal: %d Vs. %d", len(expected), len(elements))
	}

	for _, e := range elements {
		action := e.GetActionType()
		capability := e.GetActionCapability()
		p, exist := expected[action]
		if !exist {
			return fmt.Errorf("action type(%v) not exist", action)
		}

		if p != capability {
			return fmt.Errorf("action(%v) policy mismatch %v Vs %v", action, capability, p)
		}
	}

	return nil
}

func TestK8sRegistrationClient_GetActionPolicy(t *testing.T) {
	conf := NewRegistrationClientConfig(stitching.UUID, 0, true)
	reg := NewK8sRegistrationClient(conf)

	supported := proto.ActionPolicyDTO_SUPPORTED
	recommend := proto.ActionPolicyDTO_NOT_EXECUTABLE
	notSupported := proto.ActionPolicyDTO_NOT_SUPPORTED

	pod := proto.EntityDTO_CONTAINER_POD
	container := proto.EntityDTO_CONTAINER
	app := proto.EntityDTO_APPLICATION

	move := proto.ActionItemDTO_MOVE
	resize := proto.ActionItemDTO_RIGHT_SIZE
	provision := proto.ActionItemDTO_PROVISION

	expectedPod := make(map[proto.ActionItemDTO_ActionType]proto.ActionPolicyDTO_ActionCapability)
	expectedPod[move] = supported
	expectedPod[resize] = notSupported
	expectedPod[provision] = supported

	expectedContainer := make(map[proto.ActionItemDTO_ActionType]proto.ActionPolicyDTO_ActionCapability)
	expectedContainer[move] = notSupported
	expectedContainer[resize] = supported
	expectedContainer[provision] = recommend

	expectedApp := make(map[proto.ActionItemDTO_ActionType]proto.ActionPolicyDTO_ActionCapability)
	expectedApp[move] = notSupported
	expectedApp[resize] = notSupported
	expectedApp[provision] = recommend

	policies := reg.GetActionPolicy()

	for _, item := range policies {
		entity := item.GetEntityType()
		expected := expectedPod

		if entity == pod {
			expected = expectedPod
		} else if entity == container {
			expected = expectedContainer
		} else if entity == app {
			expected = expectedApp
		} else {
			t.Errorf("Unknown entity type: %v", entity)
			continue
		}

		if err := xcheck(expected, item.GetPolicyElement()); err != nil {
			t.Errorf("Failed action policy check for entity(%v) %v", entity, err)
		}
	}
}

func TestK8sRegistrationClient_GetEntityMetadata(t *testing.T) {
	conf := NewRegistrationClientConfig(stitching.UUID, 0, true)
	reg := NewK8sRegistrationClient(conf)

	//1. all the entity types
	entities := []proto.EntityDTO_EntityType{
		proto.EntityDTO_VIRTUAL_DATACENTER,
		proto.EntityDTO_VIRTUAL_MACHINE,
		proto.EntityDTO_CONTAINER_POD,
		proto.EntityDTO_CONTAINER,
		proto.EntityDTO_APPLICATION,
		proto.EntityDTO_VIRTUAL_APPLICATION,
	}
	entitySet := make(map[proto.EntityDTO_EntityType]struct{})

	for _, etype := range entities {
		entitySet[etype] = struct{}{}
	}

	//2. verify all the entity MetaData
	metaData := reg.GetEntityMetadata()
	if len(metaData) != len(entitySet) {
		t.Errorf("EntityMetadata count dis-match: %d vs %d", len(metaData), len(entitySet))
	}

	for _, meta := range metaData {
		etype := meta.GetEntityType()
		if _, exist := entitySet[etype]; !exist {
			t.Errorf("Unexpected EntityType: %v", etype)
		}

		properties := meta.GetNonVolatileProperties()
		if len(properties) != 1 {
			t.Errorf("Number of NonVolatieProperties should be 1 Vs. %v", len(properties))
		}

		if properties[0].GetName() != propertyId {
			t.Errorf("Property name should be : %v Vs. %v", propertyId, properties[0].GetName())
		}
	}

}
