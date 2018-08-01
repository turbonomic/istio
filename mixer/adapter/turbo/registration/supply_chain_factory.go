/*
	The supply chain factory is responsible for creation of the supply chain to be passed onto to the Turbo.
*/
package registration

import (
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"github.com/turbonomic/turbo-go-sdk/pkg/supplychain"
)

var (
	vFlowType         = proto.CommodityDTO_FLOW
	vFlowTemplateComm = &proto.TemplateCommodity{CommodityType: &vFlowType}
)

type SupplyChainFactory struct {
	vmTemplateType proto.TemplateDTO_TemplateType
}

func NewSupplyChainFactory() *SupplyChainFactory {
	tmpType := proto.TemplateDTO_EXTENSION
	return &SupplyChainFactory{
		vmTemplateType: tmpType,
	}
}

func (f *SupplyChainFactory) createSupplyChain() ([]*proto.TemplateDTO, error) {
	podSupplyChainNode, err := f.buildPodSupplyBuilder()
	if err != nil {
		return nil, err
	}
	return supplychain.NewSupplyChainBuilder().Top(podSupplyChainNode).Create()
}

// ContainerPod link in the supply chain.
// There is no supply chain per-say here, as we do not discover any service entities.
func (f *SupplyChainFactory) buildPodSupplyBuilder() (*proto.TemplateDTO, error) {
	podSupplyChainNodeBuilder := supplychain.NewSupplyChainNodeBuilder(proto.EntityDTO_CONTAINER_POD)
	return podSupplyChainNodeBuilder.Create()
}
