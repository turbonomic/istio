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
		metricTypes   map[string]*metric.Type
		tapSvc        *IstioTAPService
		cfg           *config.Params
		metricHandler *discovery.MetricHandler
	}

	handler struct {
		bld         *builder
		metricTypes map[string]*metric.Type
		logger      adapter.Logger
	}
)

type disconnectFromTurboFunc func()

var (
	_ metric.HandlerBuilder = &builder{}
	_ metric.Handler        = &handler{}
)

// adapter.HandlerBuilder#Build
func (b *builder) Build(ctx context.Context, env adapter.Env) (adapter.Handler, error) {
	host, err := os.Hostname()
	if err != nil {
		env.Logger().Errorf("Error retrieving host name: %s\n", err)
		return nil, err
	}
	if !strings.Contains(host, "istio-telemetry") {
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
	tapSvc.ConnectToTurbo()
	return &handler{
		bld:         b,
		metricTypes: b.metricTypes,
		logger:      env.Logger(),
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
	b.metricTypes = types
}

////////////////// Request-time Methods //////////////////////////

// Builds the metric
func (h *handler) buildMetric(dimensions map[string]interface{}) *discovery.Metric {
	builder := h.bld.metricHandler.NewMetricBuilder()
	for key, value := range dimensions {
		switch key {
		case "source_ip":
			builder = builder.WithSource(value.(string))
		case "destination_ip":
			builder = builder.WithDestination(value.(string))
		case "req_size":
			builder = builder.WithTransmittedAmount(value.(int64))
		case "resp_size":
			builder = builder.WithReceivedAmount(value.(int64))
		case "latency":
			duration := value.(time.Duration)
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
