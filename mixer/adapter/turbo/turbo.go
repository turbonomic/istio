//go:generate $GOPATH/src/istio.io/istio/bin/mixer_codegen.sh -f mixer/adapter/turbo/config/config.proto -x "-n turbo -t metric"
package turbo

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"istio.io/istio/mixer/adapter/turbo/config"
	"istio.io/istio/mixer/adapter/turbo/discovery"
	"istio.io/istio/mixer/pkg/adapter"
	"istio.io/istio/mixer/template/metric"
	"time"
	"strings"
	"github.com/pkg/errors"
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
	// Disconnect from Turbonomic server when mixer is shutdown
	handleExit(func() { tapSvc.DisconnectFromTurbo() })
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

// handleExit disconnects the tap service from Turbo service when Istio is shutdown
func handleExit(disconnectFunc disconnectFromTurboFunc) {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGHUP)

	go func() {
		select {
		case <-sigChan:
			disconnectFunc()
		}
	}()
}

// Parses the TAP service spec.
// Constructs the JSON, and then parses it.
func (b *builder) parseConfig() (*discovery.IstioTAPServiceSpec, error) {
	cfgMap := fmt.Sprintf(`{"communicationConfig": {
                                        "serverMeta": {"version": "%s", "turboServer": "%s"}, 
                                        "restAPIConfig": {"opsManagerUserName": "%s", "opsManagerPassword": "%s"}
                                   }}`, b.cfg.TargetVersion, b.cfg.Url, b.cfg.User, b.cfg.Password)
	return discovery.ParseTurboCommunicationConfig(cfgMap)
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
func (h *handler) buildMetric(dimensions map[string]interface{}) *discovery.Metric {
	builder := h.bld.metricHandler.NewMetricBuilder()
	for key, value := range dimensions {
		switch key {
		case "source_ip":
			v, ok := value.(string)
			if !ok {
				return h.bld.metricHandler.NewMetricBuilder()
			}
			builder = builder.WithSource(v)
		case "destination_ip":
			v, ok := value.(string)
			if !ok {
				return h.bld.metricHandler.NewMetricBuilder()
			}
			builder = builder.WithDestination(v)
		case "req_size":
			v, ok := value.(int64)
			if !ok {
				return h.bld.metricHandler.NewMetricBuilder()
			}
			builder = builder.WithTransmittedAmount(v)
		case "resp_size":
			v, ok := value.(int64)
			if !ok {
				return h.bld.metricHandler.NewMetricBuilder()
			}
			builder = builder.WithReceivedAmount(v)
		case "latency":
			duration, ok := value.(time.Duration)
			if !ok {
				return h.bld.metricHandler.NewMetricBuilder()
			}
			builder = builder.WithDuration(int(duration.Nanoseconds() / 1000))
		}
	}
	return builder
}

// metric.Handler#HandleMetric
func (h *handler) HandleMetric(ctx context.Context, insts []*metric.Instance) error {
	for _, inst := range insts {
		if m, err := h.buildMetric(inst.Dimensions).Create(); err != nil {
			h.logger.Errorf("Error building metric %s", err)
			return err
		} else {
			h.bld.metricHandler.Add(m)
		}
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
