package metrics

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"testing"
	"time"
)

var meter = otel.Meter("test-meter")

func setup() *CircuitBreakerMetrics {
	cbMetric, err := NewCircuitBreakerMetrics(meter)
	if err != nil {
		panic(err)
	}
	return cbMetric
}

var cb = setup()
var ctx = context.Background()

func Test_AddCounter(t *testing.T) {
	tests := []struct {
		name             string
		uniqueIdentifier string
		attributeMap     map[string]string
	}{
		{
			name:             "Test case 1",
			uniqueIdentifier: "unique1",
			attributeMap:     map[string]string{"key1": "value1", "key2": "value2"},
		},
		{
			name:             "Test case 2",
			uniqueIdentifier: "unique2",
			attributeMap:     map[string]string{"key3": "value3", "key4": "value4"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddCounter(cb.StateChangeCount, ctx, tt.uniqueIdentifier, tt.attributeMap)
		})
	}
}

func TestHistogramRecord(t *testing.T) {
	duration := time.Since(time.Now())
	tests := []struct {
		name             string
		uniqueIdentifier string
		attributeMap     map[string]string
	}{
		{
			name:             "Test case 1",
			uniqueIdentifier: "unique1",
			attributeMap:     map[string]string{"key1": "value1", "key2": "value2"},
		},
		{
			name:             "Test case 2",
			uniqueIdentifier: "unique2",
			attributeMap:     map[string]string{"key3": "value3", "key4": "value4"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HistogramRecord(cb.StateDuration, ctx, tt.uniqueIdentifier, tt.attributeMap, duration)
		})
	}
}

func TestNewCircuitBreakerMetrics(t *testing.T) {
	type args struct {
		meter metric.Meter
	}
	tests := []struct {
		name    string
		args    args
		want    *CircuitBreakerMetrics
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
		{
			name: "Test case 1",
			args: args{meter: meter},
			want: &CircuitBreakerMetrics{
				StateChangeCount: cb.StateChangeCount,
				StateDuration:    cb.StateDuration,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCircuitBreakerMetrics(tt.args.meter)
			if !tt.wantErr(t, err, fmt.Sprintf("NewCircuitBreakerMetrics(%v)", tt.args.meter)) {
				return
			}
			assert.Equalf(t, tt.want, got, "NewCircuitBreakerMetrics(%v)", tt.args.meter)
		})
	}
}
