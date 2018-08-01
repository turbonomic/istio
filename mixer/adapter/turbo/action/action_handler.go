package action

import (
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"

	sdkprobe "github.com/turbonomic/turbo-go-sdk/pkg/probe"
)

type ActionHandler struct {
}

func NewIstioActionHandler() *ActionHandler {
	return &ActionHandler{}
}

func (h *ActionHandler) unsupported() *proto.ActionResult {

	state := proto.ActionResponseState_DISABLED
	progress := int32(100)
	msg := "Unsupported"

	res := &proto.ActionResponse{
		ActionResponseState: &state,
		Progress:            &progress,
		ResponseDescription: &msg,
	}

	return &proto.ActionResult{
		Response: res,
	}
}

// Implement ActionExecutorClient interface defined in Go SDK.
// Execute the current action and return the action result to SDK.
func (h *ActionHandler) ExecuteAction(actionExecutionDTO *proto.ActionExecutionDTO,
	accountValues []*proto.AccountValue, progressTracker sdkprobe.ActionProgressTracker) (*proto.ActionResult, error) {
	return h.unsupported(), nil
}
