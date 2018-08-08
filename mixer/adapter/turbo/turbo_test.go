package turbo

import (
	"testing"
	"github.com/magiconair/properties/assert"
	"istio.io/istio/mixer/adapter/turbo/discovery"
	"time"
	"istio.io/istio/mixer/template/metric"
)

func TestGetInfo(t *testing.T) {
	info := GetInfo()
	assert.Equal(t, info.Name, "turbo")
}

func TestHandler_HandleMetric(t *testing.T) {
	h := &handler{
		bld: &builder{metricHandler: discovery.NewMetricHandler()},
	}
	var dimensions map[string]interface{}
	dimensions = make(map[string]interface{})
	dimensions["source_ip"] = "10.10.1.1"
	dimensions["destination_ip"] = "10.10.2.1"
	dimensions["req_size"] = int64(10)
	dimensions["resp_size"] = int64(16)
	dimensions["latency"] = time.Minute
	insts := make([]*metric.Instance, 1, 1)
	insts[0] = &metric.Instance{
		Dimensions: dimensions,
	}
	err := h.HandleMetric(nil, insts)
	assert.Equal(t, err, nil)
}

func TestHandler_HandleMetric_Error(t *testing.T) {
	h := &handler{
		bld:    &builder{metricHandler: discovery.NewMetricHandler()},
		logger: &nilLogger{},
	}
	var dimensions map[string]interface{}
	dimensions = make(map[string]interface{})
	dimensions["source_ip"] = "10.10.1.1"
	dimensions["destination_ip"] = "10.10.2.1"
	dimensions["req_size"] = int64(10)
	dimensions["resp_size"] = int64(16)
	dimensions["latency"] = time.Now()
	insts := make([]*metric.Instance, 1)
	insts[0] = &metric.Instance{
		Dimensions: dimensions,
	}
	err := h.HandleMetric(nil, insts)
	if err == nil {
		t.Errorf("The metric must not be nil")
	}
}

// Borrows Mock logger.
type nilLogger struct{}

func (m nilLogger) Infof(format string, args ...interface{}) {}

func (m nilLogger) Warningf(format string, args ...interface{}) {}

func (m nilLogger) Errorf(format string, args ...interface{}) error { return nil }

func (m nilLogger) Debugf(format string, args ...interface{}) {}

func (m nilLogger) InfoEnabled() bool { return false }

func (m nilLogger) WarnEnabled() bool { return false }

func (m nilLogger) ErrorEnabled() bool { return false }

func (m nilLogger) DebugEnabled() bool { return false }
