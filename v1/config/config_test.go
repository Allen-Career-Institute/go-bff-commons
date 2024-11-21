package config

import (
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
	"time"
)

func TestGetCircuitBreakerClientConfigs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock the dynamic config
	dc := NewMockDynamicConfig(ctrl)

	// Define arguments struct for the test cases
	type args struct {
		client string
		cnf    *Config
	}

	// Define test cases
	tests := []struct {
		name string
		args args
		want CircuitBreakerClientConfig
		mock []*gomock.Call
	}{
		{
			name: "Local config read when DynamicConfig is nil",
			args: args{
				client: "test-client",
				cnf:    &Config{DynamicConfig: nil},
			},
			want: CircuitBreakerClientConfig{
				FailurePercentageThresholdWithinTimePeriod:   DefaultCircuitBreakerPercentageThreshold,
				FailureMinExecutionThresholdWithinTimePeriod: DefaultCircuitMinExecutionThreshold,
				FailurePeriodThreshold:                       DefaultCircuitBreakerFailurePeriodThresholdInSeconds * time.Second,
				SuccessThreshold:                             DefaultCircuitBreakerSuccessThreshold,
				Delay:                                        DefaultCircuitBreakerDelayMs * time.Millisecond,
			},
			mock: []*gomock.Call{
				// No mocks for DynamicConfig since it's nil
			},
		},
		{
			name: "Dynamic config read from AWS App Config",
			args: args{
				client: "test-client",
				cnf:    &Config{DynamicConfig: dc},
			},
			want: CircuitBreakerClientConfig{
				FailurePercentageThresholdWithinTimePeriod:   50,
				FailureMinExecutionThresholdWithinTimePeriod: 10,
				FailurePeriodThreshold:                       60 * time.Second,
				SuccessThreshold:                             10,
				Delay:                                        500 * time.Millisecond,
			},
			mock: []*gomock.Call{
				dc.EXPECT().Get("test-client"+CircuitBreakerFailurePercentageThresholdSuffix).Return("50", nil).Times(1),
				dc.EXPECT().Get("test-client"+CircuitBreakerMinExecutionThresholdSuffix).Return("10", nil).Times(1),
				dc.EXPECT().Get("test-client"+CircuitBreakerFailurePeriodThresholdSuffix).Return("60", nil).Times(1),
				dc.EXPECT().Get("test-client"+CircuitBreakerSuccessThresholdSuffix).Return("10", nil).Times(1),
				dc.EXPECT().Get("test-client"+CircuitBreakerDelaySuffix).Return("500", nil).Times(1),
			},
		},
		{
			name: "Dynamic config with parsing errors",
			args: args{
				client: "test-client",
				cnf:    &Config{DynamicConfig: dc},
			},
			want: CircuitBreakerClientConfig{
				FailurePercentageThresholdWithinTimePeriod:   DefaultCircuitBreakerPercentageThreshold,
				FailureMinExecutionThresholdWithinTimePeriod: DefaultCircuitMinExecutionThreshold,
				FailurePeriodThreshold:                       DefaultCircuitBreakerFailurePeriodThresholdInSeconds * time.Second,
				SuccessThreshold:                             DefaultCircuitBreakerSuccessThreshold,           // Use the default value
				Delay:                                        DefaultCircuitBreakerDelayMs * time.Millisecond, // Use the default value
			},
			mock: []*gomock.Call{
				dc.EXPECT().Get("test-client"+CircuitBreakerFailurePercentageThresholdSuffix).Return("invalid", nil),
				dc.EXPECT().Get("test-client"+CircuitBreakerMinExecutionThresholdSuffix).Return("invalid", nil),
				dc.EXPECT().Get("test-client"+CircuitBreakerFailurePeriodThresholdSuffix).Return("invalid", nil),
				dc.EXPECT().Get("test-client"+CircuitBreakerSuccessThresholdSuffix).Return("invalid", nil),
				dc.EXPECT().Get("test-client"+CircuitBreakerDelaySuffix).Return("invalid", nil),
			},
		},
	}

	// Loop over each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function with the arguments
			got := GetCircuitBreakerClientConfigs(tt.args.client, tt.args.cnf)

			// Compare the actual result with the expected result
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCircuitBreakerClientConfigs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRetryClientConfigs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock the dynamic config
	dc := NewMockDynamicConfig(ctrl)

	// Define arguments struct for the test cases
	type args struct {
		client string
		cnf    *Config
	}

	// Define test cases
	tests := []struct {
		name string
		args args
		want RetryClientConfig
		mock []*gomock.Call
	}{
		{
			name: "Local config read when DynamicConfig is nil",
			args: args{
				client: "test-client",
				cnf:    &Config{DynamicConfig: nil},
			},
			want: RetryClientConfig{
				MaxRetries: DefaultRetryMaxRetries,
				Delay:      DefaultRetryDelayMs * time.Millisecond,
			},
			mock: []*gomock.Call{
				// No mocks for DynamicConfig since it's nil
			},
		},
		{
			name: "Dynamic config read from AWS App Config",
			args: args{
				client: "test-client",
				cnf:    &Config{DynamicConfig: dc},
			},
			want: RetryClientConfig{
				MaxRetries: 3,
				Delay:      1 * time.Second,
			},
			mock: []*gomock.Call{
				dc.EXPECT().Get("test-client"+RetryMaxRetriesSuffix).Return("3", nil),
				dc.EXPECT().Get("test-client"+RetryDelaySuffix).Return("1000", nil),
			},
		},
		{
			name: "Dynamic config with parsing errors",
			args: args{
				client: "test-client",
				cnf:    &Config{DynamicConfig: dc},
			},
			want: RetryClientConfig{
				MaxRetries: DefaultRetryMaxRetries,
				Delay:      DefaultRetryDelayMs * time.Millisecond,
			},
			mock: []*gomock.Call{
				dc.EXPECT().Get("test-client"+RetryMaxRetriesSuffix).Return("invalid", nil),
				dc.EXPECT().Get("test-client"+RetryDelaySuffix).Return("invalid", nil),
			},
		},
	}

	// Loop over each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute the mock expectations
			for _, m := range tt.mock {
				m.Times(1)
			}

			// Call the function with the arguments
			got := GetRetryClientConfigs(tt.args.client, tt.args.cnf)

			// Compare the actual result with the expected result
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRetryClientConfigs() = %v, want %v", got, tt.want)
			}
		})
	}
}
