package turbo

import "istio.io/istio/mixer/adapter/turbo/discovery"

// The Istio probe-side configuration. We pass parameters through it.
type Config struct {
	tapSpec              *discovery.IstioTAPServiceSpec
	DiscoveryIntervalSec int
	metricHandler        *discovery.MetricHandler
}

func NewVMTConfig() *Config {
	return &Config{}
}

func (c *Config) WithTapSpec(spec *discovery.IstioTAPServiceSpec) *Config {
	c.tapSpec = spec
	return c
}

func (c *Config) WithDiscoveryInterval(di int) *Config {
	c.DiscoveryIntervalSec = di
	return c
}

func (c *Config) WithMetricHandler(metricHandler *discovery.MetricHandler) *Config {
	c.metricHandler = metricHandler
	return c
}
