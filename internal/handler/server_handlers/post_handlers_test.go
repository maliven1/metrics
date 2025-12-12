package serverhandlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	mock_serverhandlers "github.com/maliven1/metrics/internal/handler/server_handlers/mocks"
	models "github.com/maliven1/metrics/internal/model"
	"go.uber.org/zap"
)

func TestHandler_PostMetricsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_serverhandlers.NewMockService(ctrl)
	handler := NewHandler(mockService, nil)

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugaredLogger := logger.Sugar()

	tests := []struct {
		name           string
		metrics        []models.Metrics
		expectedStatus int
		setupMock      func()
	}{
		{
			name: "valid metrics",
			metrics: []models.Metrics{
				{ID: "TestMetric1", MType: models.Gauge, Value: func() *float64 { v := 123.45; return &v }()},
				{ID: "TestMetric2", MType: models.Counter, Delta: func() *int64 { v := int64(678); return &v }()},
			},
			expectedStatus: http.StatusOK,
			setupMock: func() {
				mockService.EXPECT().SetMetrics(gomock.Any()).Times(1)
			},
		},
		{
			name:           "nil metrics",
			metrics:        nil,
			expectedStatus: http.StatusBadRequest,
			setupMock:      func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			metricsJSON, _ := json.Marshal(tt.metrics)
			req := httptest.NewRequest(http.MethodPost, "/updates/", bytes.NewBuffer(metricsJSON))
			w := httptest.NewRecorder()

			handler.PostMetricsHandler(sugaredLogger)(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestHandler_PostBodyHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_serverhandlers.NewMockService(ctrl)
	handler := NewHandler(mockService, nil)

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugaredLogger := logger.Sugar()

	tests := []struct {
		name           string
		metric         models.Metrics
		expectedStatus int
		setupMock      func()
	}{
		{
			name:           "valid gauge metric",
			metric:         models.Metrics{ID: "TestMetric", MType: models.Gauge, Value: func() *float64 { v := 123.45; return &v }()},
			expectedStatus: http.StatusOK,
			setupMock: func() {
				mockService.EXPECT().AddStructMetric(gomock.Any()).Return(nil).Times(1)
			},
		},
		{
			name:           "valid counter metric",
			metric:         models.Metrics{ID: "TestCounter", MType: models.Counter, Delta: func() *int64 { v := int64(678); return &v }()},
			expectedStatus: http.StatusOK,
			setupMock: func() {
				mockService.EXPECT().AddStructMetric(gomock.Any()).Return(nil).Times(1)
			},
		},
		{
			name:           "invalid metric type",
			metric:         models.Metrics{ID: "TestMetric", MType: "invalid", Value: func() *float64 { v := 123.45; return &v }()},
			expectedStatus: http.StatusBadRequest,
			setupMock: func() {
				mockService.EXPECT().AddStructMetric(gomock.Any()).Return(fmt.Errorf("error")).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			metricJSON, _ := json.Marshal(tt.metric)
			req := httptest.NewRequest(http.MethodPost, "/update/", bytes.NewBuffer(metricJSON))
			w := httptest.NewRecorder()

			handler.PostBodyHandler(sugaredLogger)(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestHandler_PostURLHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_serverhandlers.NewMockService(ctrl)
	handler := NewHandler(mockService, nil)

	tests := []struct {
		name           string
		url            string
		expectedStatus int
		setupMock      func()
	}{
		{
			name:           "valid gauge metric update",
			url:            "/update/gauge/TestMetric/123.45",
			expectedStatus: http.StatusOK,
			setupMock: func() {
				mockService.EXPECT().CheckAddPath(gomock.Any()).Return(nil).Times(1)
			},
		},
		{
			name:           "valid counter metric update",
			url:            "/update/counter/TestCounter/678",
			expectedStatus: http.StatusOK,
			setupMock: func() {
				mockService.EXPECT().CheckAddPath(gomock.Any()).Return(nil).Times(1)
			},
		},
		{
			name:           "invalid path length",
			url:            "/update/gauge/TestMetric",
			expectedStatus: http.StatusNotFound,
			setupMock:      func() {},
		},
		{
			name:           "invalid metric type",
			url:            "/update/invalid/TestMetric/123.45",
			expectedStatus: http.StatusBadRequest,
			setupMock: func() {
				mockService.EXPECT().CheckAddPath(gomock.Any()).Return(fmt.Errorf("error")).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			req := httptest.NewRequest(http.MethodPost, tt.url, nil)
			w := httptest.NewRecorder()

			handler.PostURLHandler()(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
