package metrics

import (
	"context"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/logger"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
)

func Test_addCounter(t *testing.T) {
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
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("AddCounter() panicked unexpectedly")
				}
			}()

			c := &config.Config{Logger: config.Logger{Level: "info"}}
			log := logger.NewAPILogger(c)
			log.InitLogger()

			customMetrics, err := NewCustomMetrics(log)
			if err != nil {
				t.Errorf("Failed to create CustomMetrics: %v", err)
				return
			}

			customMetrics.AddCounter(context.Background(), tt.uniqueIdentifier, tt.attributeMap)
		})
	}
}

func TestCreateMetrics(t *testing.T) {
	tests := []struct {
		name     string
		counters []MetricParams
		wantErr  bool
	}{
		{
			name: "Test case 1: Single counter",
			counters: []MetricParams{
				{Name: "counter1", Desc: "desc1"},
			},
			wantErr: false,
		},
		{
			name: "Test case 2: Multiple counters",
			counters: []MetricParams{
				{Name: "counter1", Desc: "desc1"},
				{Name: "counter2", Desc: "desc2"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meter := otel.GetMeterProvider().Meter("bff-service")
			got, err := createMetrics(meter, tt.counters)

			if (err != nil) != tt.wantErr {
				t.Errorf("createMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, counter := range tt.counters {
				if _, ok := got[counter.Name]; !ok {
					t.Errorf("createMetrics() = %v, want %v", got, tt.counters)
				}
			}
		})
	}
}

func TestCreateCounter(t *testing.T) {
	tests := []struct {
		name   string
		params MetricParams
	}{
		{
			name:   "Params does not contain field",
			params: MetricParams{},
		},
		{
			name: "Params contain Name field",
			params: MetricParams{
				Name: "counter1",
			},
		},
		{
			name: "Params contain Desc field",
			params: MetricParams{
				Desc: "desc1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meter := otel.GetMeterProvider().Meter("bff-service")
			_, err := createCounter(meter, tt.params)
			assert.NoError(t, err)
		})
	}
}

func TestNewCustomMetrics(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Test case 1: Valid logger",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &config.Config{Logger: config.Logger{Level: "info"}}
			log := logger.NewAPILogger(c)
			log.InitLogger()

			got, err := NewCustomMetrics(log)
			if err != nil {
				t.Errorf("NewCustomMetrics() error = %v, want nil", err)
				return
			}

			if got == nil {
				t.Errorf("NewCustomMetrics() = nil, want non-nil")
			}
		})
	}
}
