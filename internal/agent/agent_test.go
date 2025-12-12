package agent

import (
	"reflect"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/maliven1/metrics/internal/agent/mocks"
	"github.com/maliven1/metrics/internal/config"
)

func TestAgent_GetMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockMemStorage(ctrl)
	agent := Agent{
		memStorage: mockStorage,
		cfg:        &config.AgentConfig{PollInterval: 1},
		m:          &sync.Mutex{},
	}

	// Set up expectations
	gaugeMap := map[string]float64{"test": 1.0}
	counterMap := map[string]int64{"test": 1}
	mockStorage.EXPECT().GetGauge().Return(gaugeMap)
	mockStorage.EXPECT().GetCounter().Return(counterMap)

	// Call the function
	gauges, counters := agent.GetMetrics()

	// Check the results
	if !reflect.DeepEqual(gauges, gaugeMap) {
		t.Errorf("GetMetrics() gauges = %v, want %v", gauges, gaugeMap)
	}
	if !reflect.DeepEqual(counters, counterMap) {
		t.Errorf("GetMetrics() counters = %v, want %v", counters, counterMap)
	}
}
