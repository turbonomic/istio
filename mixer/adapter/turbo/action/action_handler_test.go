// Tests the action_handler
package action

import (
	"testing"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"github.com/magiconair/properties/assert"
)

func TestBasic(t *testing.T) {
	handler := NewIstioActionHandler()
	result := handler.unsupported()
	assert.Equal(t, *result.Response.ActionResponseState, proto.ActionResponseState_DISABLED)
}

func TestActionHandler_ExecuteAction(t *testing.T) {
	handler := NewIstioActionHandler()
	result, err := handler.ExecuteAction(nil, nil, nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, *result.Response.ActionResponseState, proto.ActionResponseState_DISABLED)
}
