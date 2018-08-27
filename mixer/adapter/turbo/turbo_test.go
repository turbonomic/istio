package turbo

import (
	"net"
	"testing"
	"github.com/magiconair/properties/assert"
	"istio.io/istio/mixer/adapter/turbo/discovery"
	"time"
	"istio.io/istio/mixer/template/metric"
	"istio.io/istio/mixer/adapter/turbo/config"
	"istio.io/istio/mixer/pkg/adapter/test"
)

func TestGetInfo(t *testing.T) {
	info := GetInfo()
	assert.Equal(t, info.Name, "turbo")
}

func instanceOf(name string, value interface{}) *metric.Instance {
	return &metric.Instance{
		Name:  name,
		Value: value,
	}
}

func TestHandler_HandleMetric(t *testing.T) {
	h := &handler{
		bld: &builder{metricHandler: discovery.NewMetricHandler()},
	}
	insts := make([]*metric.Instance, 0, 5)
	insts = append(insts, instanceOf("srcip.metric.istio-system", net.ParseIP("10.10.1.1")))
	insts = append(insts, instanceOf("dstip.metric.istio-system", net.ParseIP("10.10.2.1")))
	insts = append(insts, instanceOf("reqsize.metric.istio-system", int64(10)))
	insts = append(insts, instanceOf("respsize.metric.istio-system", int64(16)))
	insts = append(insts, instanceOf("latency.metric.istio-system", time.Minute))
	err := h.HandleMetric(nil, insts)
	assert.Equal(t, err, nil)
}

func ensureMetricType(t *testing.T, dimension string, value interface{}) {
	h := &handler{
		bld:    &builder{metricHandler: discovery.NewMetricHandler()},
		logger: &nilLogger{},
	}
	// Set the working defaults
	insts := make([]*metric.Instance, 0, 5)
	insts = append(insts, instanceOf("srcip.metric.istio-system", net.ParseIP("10.10.1.1")))
	insts = append(insts, instanceOf("dstip.metric.istio-system", net.ParseIP("10.10.2.1")))
	insts = append(insts, instanceOf("reqsize.metric.istio-system", int64(10)))
	insts = append(insts, instanceOf("respsize.metric.istio-system", int64(16)))
	insts = append(insts, instanceOf("latency.metric.istio-system", time.Minute))
	err := h.HandleMetric(nil, insts)
	if err == nil {
		t.Errorf("Failed test for " + dimension)
	}
}

func TestHandler_HandleMetric_Error(t *testing.T) {
	// Test all parts
	insts := make([]*metric.Instance, 0, 5)
	insts = append(insts, instanceOf("srcip.metric.istio-system", 100))
	insts = append(insts, instanceOf("dstip.metric.istio-system", 100))
	insts = append(insts, instanceOf("reqsize.metric.istio-system", "100"))
	insts = append(insts, instanceOf("respsize.metric.istio-system", "100"))
	insts = append(insts, instanceOf("latency.metric.istio-system",100))
}

func TestBuilder_SetAdapterConfig(t *testing.T) {
	b := &builder{metricHandler: discovery.NewMetricHandler()}
	cfg := new(config.Params)
	b.SetAdapterConfig(cfg)
	assert.Equal(t, b.cfg, cfg)
}

func TestBuilder_Validate(t *testing.T) {
	b := &builder{metricHandler: discovery.NewMetricHandler()}
	rc := b.Validate()
	if rc != nil {
		t.Errorf("Must return nil")
	}
}

func TestBuilder_SetMetricTypes(t *testing.T) {
	b := &builder{metricHandler: discovery.NewMetricHandler()}
	metricTypes := make(map[string]*metric.Type)
	b.SetMetricTypes(metricTypes)
}

func TestParseConfig(t *testing.T) {
	b := &builder{metricHandler: discovery.NewMetricHandler()}
	cfg := new(config.Params)
	cfg.Url = "https://localhost"
	cfg.User = "user"
	cfg.Password = "Passwd"
	cfg.TargetVersion = "1.0.0"
	b.SetAdapterConfig(cfg)
	_, err := b.parseConfig()
	assert.Equal(t, err, nil)
}

func TestBuilder_Build(t *testing.T) {
	b := &builder{metricHandler: discovery.NewMetricHandler()}
	cfg := new(config.Params)
	cfg.Url = "https://abcdef"
	cfg.User = "user"
	cfg.Password = "Passwd"
	cfg.TargetVersion = "1.0.0"
	b.SetAdapterConfig(cfg)
	env := test.NewEnv(t)
	_, err := b.Build(nil, env)
	assert.Equal(t, err, nil)
}

func TestBuilder_Build_Bad_Param(t *testing.T) {
	b := &builder{metricHandler: discovery.NewMetricHandler()}
	cfg := new(config.Params)
	cfg.Url = "abcdef"
	cfg.User = "user"
	cfg.Password = "Passwd"
	cfg.TargetVersion = "1.0.0"
	b.SetAdapterConfig(cfg)
	env := test.NewEnv(t)
	_, err := b.Build(nil, env)
	if err == nil {
		t.Errorf("Should have been an error")
	}
}

func TestBuilder_Build_Bad_Host(t *testing.T) {
	b := &builder{metricHandler: discovery.NewMetricHandler()}
	cfg := new(config.Params)
	cfg.Url = "https://abcdef"
	cfg.User = "user"
	cfg.Password = "Passwd"
	cfg.TargetVersion = "1.0.0"
	b.SetAdapterConfig(cfg)
	env := test.NewEnv(t)
	allowedHost = ":::No:::Such:::Host"
	_, err := b.Build(nil, env)
	if err == nil {
		t.Errorf("Should have been an error")
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
