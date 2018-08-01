package discovery

import (
	"sync"
	"errors"
	"fmt"
)

type (
	// The individual metric
	Metric struct {
		src      string
		dst      string
		amount   float64
		duration int
		rx       int64
		tx       int64
	}
	// The synchronized handler
	MetricHandler struct {
		metrics map[string]*Metric
		lock    sync.Mutex
	}
)

// Create an instance of the metric handler
func NewMetricHandler() *MetricHandler {
	return &MetricHandler{
		metrics: make(map[string]*Metric),
	}
}

// Creates a new metric
func (handler *MetricHandler) NewMetricBuilder() *Metric {
	return &Metric{
	}
}

// Add the metric
//
// Compose the key as follows: source, destination, destinatio_port
// If the entry for the key exists, accumulate the amounts and duration,
// if it doesn't, set the new one for the key.
func (handler *MetricHandler) Add(metric *Metric) error {
	handler.lock.Lock()
	defer handler.lock.Unlock()
	// Add the metric
	if metric == nil {
		return errors.New("no metric supplied")
	}
	// Add
	key := metric.src + metric.dst
	if value, present := handler.metrics[key]; !present {
		handler.metrics[key] = metric
	} else {
		value.amount += metric.amount
		value.duration += metric.duration
		value.rx += metric.rx
		value.tx += metric.tx
	}
	return nil
}

// Returns the accumulated metrics and resets the stored one.
func (handler *MetricHandler) GetMetrics() map[string]*Metric {
	handler.lock.Lock()
	defer handler.lock.Unlock()
	// Retrieve and clear the map
	value := handler.metrics
	handler.metrics = make(map[string]*Metric)
	return value
}

// Sets the source
func (metric *Metric) WithSource(src string) *Metric {
	metric.src = src
	return metric
}

// Sets the destination
func (metric *Metric) WithDestination(dst string) *Metric {
	metric.dst = dst
	return metric
}

// Sets the flow amount
func (metric *Metric) WithFlowAmount(amount float64) *Metric {
	metric.amount = amount
	return metric
}

// Sets the flow duration
func (metric *Metric) WithDuration(duration int) *Metric {
	metric.duration = duration
	return metric
}

// Sets the received amount
func (metric *Metric) WithReceivedAmount(amount int64) *Metric {
	metric.rx = amount
	return metric
}

// Sets the transmitted amount
func (metric *Metric) WithTransmittedAmount(amount int64) *Metric {
	metric.tx = amount
	return metric
}

// Creates the metric
func (metric *Metric) Create() (*Metric, error) {
	if len(metric.src) == 0 || len(metric.dst) == 0 {
		return nil, errors.New(fmt.Sprintf("both source and destination must be present: %s -> %s", metric.src, metric.dst))
	} else {
		metric.amount = float64(metric.rx + metric.tx)
		return metric, nil
	}
}
