package discovery

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewMetricHandler(t *testing.T) {
	h := NewMetricHandler()
	if h == nil {
		t.Errorf("Nil handler")
	}
}

func TestNewMetricBuilder(t *testing.T) {
	h := NewMetricHandler()
	b := h.NewMetricBuilder()
	if b == nil {
		t.Errorf("Nil handler")
	}
	// Test bad metric
	_, err := b.Create()
	assert.NotNil(t, err)
	// Restart
	b = h.NewMetricBuilder()
	b.WithSource("src")
	b.WithDestination("dst")
	b.WithReceivedAmount(1)
	b.WithTransmittedAmount(1)
	b.WithFlowAmount(2)
	b.WithDuration(10)
	assert.EqualValues(t, "src", b.src)
	assert.EqualValues(t, "dst", b.dst)
	assert.EqualValues(t, 1, b.rx)
	assert.EqualValues(t, 1, b.tx)
	assert.EqualValues(t, 2, b.amount)
	assert.EqualValues(t, 10, b.duration)
}

func TestMetricHandler_Add(t *testing.T) {
	h := NewMetricHandler()
	b := h.NewMetricBuilder()
	// Error
	err := h.Add(nil)
	assert.NotNil(t, err)
	// Continue
	b.WithSource("src")
	b.WithDestination("dst")
	b.WithReceivedAmount(1)
	b.WithTransmittedAmount(1)
	b.WithFlowAmount(2)
	b.WithDuration(10)
	m, _ := b.Create()
	h.Add(m)
	b = h.NewMetricBuilder()
	b.WithSource("src")
	b.WithDestination("dst")
	b.WithReceivedAmount(2)
	b.WithTransmittedAmount(2)
	b.WithFlowAmount(4)
	b.WithDuration(10)
	m1, _ := b.Create()
	h.Add(m1)
	// Check the results
	metrics := h.GetMetrics()
	assert.Equal(t, 1, len(metrics))
	assert.Equal(t, metrics["srcdst"].amount, float64(6))
	assert.Equal(t, metrics["srcdst"].rx, int64(3))
	assert.Equal(t, metrics["srcdst"].tx, int64(3))
}
