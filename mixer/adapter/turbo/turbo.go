//go:generate $GOPATH/src/istio.io/istio/bin/mixer_codegen.sh -f mixer/adapter/turbo/config/config.proto -x "-n turbo -t metric"
package turbo

import (
	"context"
	"fmt"
	"os"
	"istio.io/istio/mixer/adapter/turbo/config"
	"istio.io/istio/mixer/adapter/turbo/discovery"
	"istio.io/istio/mixer/pkg/adapter"
	"istio.io/istio/mixer/template/metric"
	"time"
	"strings"
	"github.com/pkg/errors"
	"net"
)

const (
	discoveryInterval = 600
)

type (
	builder struct {
		// maps instance_name to collector.
		tapSvc        *IstioTAPService
		cfg           *config.Params
		metricHandler *discovery.MetricHandler
	}

	handler struct {
		bld    *builder
		logger adapter.Logger
	}
)

type disconnectFromTurboFunc func()

var (
	_ metric.HandlerBuilder = &builder{}
	_ metric.Handler        = &handler{}
	// Make it a variable for the unit testing purposes.
	// In Istio 1.0, the Mixer component gets deployed twice, yet we can only work with the telemetry part.
	// The only way I found to do it was to query and filter by the host name.
	allowedHost = "istio-telemetry"
)

// adapter.HandlerBuilder#Build
func (b *builder) Build(ctx context.Context, env adapter.Env) (adapter.Handler, error) {
	host, err := os.Hostname()
	if err != nil {
		env.Logger().Errorf("Error retrieving host name: %s\n", err)
		return nil, err
	}
	if !strings.Contains(host, allowedHost) {
		env.Logger().Errorf("Unsupported istio mixer: %s\n", host)
		return nil, errors.New("Unsupported istio mixer")
	}
	// Initialize the probe part.
	// Convert to JSON, so that we can parse that:
	tapSpec, err := b.parseConfig()
	if err != nil {
		return nil, err
	}
	b.metricHandler = discovery.NewMetricHandler()
	vmtConfig := NewVMTConfig()
	vmtConfig.WithDiscoveryInterval(discoveryInterval).
		WithTapSpec(tapSpec).WithMetricHandler(b.metricHandler)
	tapSvc, tapErr := NewIstioTAPService(vmtConfig)
	if tapErr != nil {
		return nil, err
	}
	b.tapSvc = tapSvc
	// Connect
	go tapSvc.ConnectToTurbo()
	return &handler{
		bld:    b,
		logger: env.Logger(),
	}, nil
}

// adapter.HandlerBuilder#SetAdapterConfig
func (b *builder) SetAdapterConfig(cfg adapter.Config) {
	b.cfg = cfg.(*config.Params)
}

// Parses the TAP service spec.
// Constructs the JSON, and then parses it.
func (b *builder) parseConfig() (*discovery.IstioTAPServiceSpec, error) {
	cfgMap := fmt.Sprintf(`{"communicationConfig": {
                                        "serverMeta": {"version": "%s", "turboServer": "%s"}, 
                                        "restAPIConfig": {"opsManagerUserName": "%s", "opsManagerPassword": "%s"}
                                   }}`, b.cfg.TargetVersion, b.cfg.Url, b.cfg.User, b.cfg.Password)
	return discovery.ParseTurboCommunicationConfig(cfgMap, b.cfg.AdapterId)
}

// adapter.HandlerBuilder#Validate
func (b *builder) Validate() (ce *adapter.ConfigErrors) {
	return
}

// metric.HandlerBuilder#SetMetricTypes
func (b *builder) SetMetricTypes(types map[string]*metric.Type) {
}

////////////////// Request-time Methods //////////////////////////

// Builds the metric
func (h *handler) buildMetric(name string, value interface{}, builder *discovery.Metric) *discovery.Metric {
	switch name {
	case "srcip.metric.istio-system":
		v, ok := value.(net.IP)
		if !ok {
			return h.bld.metricHandler.NewMetricBuilder()
		}
		builder = builder.WithSource(v.String())
	case "dstip.metric.istio-system":
		v, ok := value.(net.IP)
		if !ok {
			return h.bld.metricHandler.NewMetricBuilder()
		}
		builder = builder.WithDestination(v.String())
	case "reqsize.metric.istio-system":
		v, ok := value.(int64)
		if !ok {
			return h.bld.metricHandler.NewMetricBuilder()
		}
		builder = builder.WithTransmittedAmount(v)
	case "respsize.metric.istio-system":
		v, ok := value.(int64)
		if !ok {
			return h.bld.metricHandler.NewMetricBuilder()
		}
		builder = builder.WithReceivedAmount(v)
	case "latency.metric.istio-system":
		duration, ok := value.(time.Duration)
		if !ok {
			return h.bld.metricHandler.NewMetricBuilder()
		}
		builder = builder.WithDuration(int(duration.Nanoseconds() / 1000000))
	}
	return builder
}

// metric.Handler#HandleMetric
func (h *handler) HandleMetric(ctx context.Context, insts []*metric.Instance) error {
	builder := h.bld.metricHandler.NewMetricBuilder()
	for _, inst := range insts {
		builder = h.buildMetric(inst.Name, inst.Value, builder)
	}
	if m, err := builder.Create(); err != nil {
		h.logger.Errorf("Error building metric %s", err)
	} else {
		h.bld.metricHandler.Add(m)
	}
	return nil
}

// adapter.Handler#Close
func (h *handler) Close() error {
	go h.bld.tapSvc.DisconnectFromTurbo()
	return nil
}

////////////////// Bootstrap //////////////////////////
// GetInfo returns the adapter.Info specific to this adapter.
func GetInfo() adapter.Info {
	return adapter.Info{
		Name:        "turbo",
		Impl:        "istio.io/istio/mixer/adapter/turbo",
		Description: "Sends the communication metrics to a Turbo instance.",
		SupportedTemplates: []string{
			metric.TemplateName,
		},
		NewBuilder:    func() adapter.HandlerBuilder { return &builder{} },
		DefaultConfig: &config.Params{},
	}
}
