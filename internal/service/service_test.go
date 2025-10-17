package service

import (
	"testing"

	"github.com/golang/mock/gomock"
	models "github.com/maliven1/metrics/internal/model"
	mock_service "github.com/maliven1/metrics/internal/service/mocks"
)

func TestMemService_AddStructMetric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_service.NewMockMemRepo(ctrl)
	s := NewService(m)

	type args struct {
		metric models.Metrics
	}
	tests := []struct {
		name      string
		args      args
		setupMock func(*mock_service.MockMemRepo)
		wantErr   bool
	}{
		{
			name: "ok gauge test",
			args: args{metric: models.Metrics{ID: "SomeMetric", MType: models.Gauge, Value: func() *float64 { v := 321.0; return &v }()}},
			setupMock: func(m *mock_service.MockMemRepo) {
				m.EXPECT().SetGauge("SomeMetric", 321.0).Times(1)
			},
			wantErr: false,
		},
		{
			name: "ok counter test with existing counter",
			args: args{metric: models.Metrics{ID: "SomeCounter", MType: models.Counter, Delta: func() *int64 { v := int64(123); return &v }()}},
			setupMock: func(m *mock_service.MockMemRepo) {
				m.EXPECT().AddCounter("SomeCounter", int64(123)).Return(true).Times(1)
			},
			wantErr: false,
		},
		{
			name: "ok counter test with new counter",
			args: args{metric: models.Metrics{ID: "SomeCounter", MType: models.Counter, Delta: func() *int64 { v := int64(123); return &v }()}},
			setupMock: func(m *mock_service.MockMemRepo) {
				m.EXPECT().AddCounter("SomeCounter", int64(123)).Return(false).Times(1)
				m.EXPECT().SetCounter("SomeCounter", int64(123)).Times(1)
			},
			wantErr: false,
		},
		{
			name: "error invalid metric type",
			args: args{metric: models.Metrics{ID: "SomeMetric", MType: "invalid", Value: func() *float64 { v := 321.0; return &v }()}},
			setupMock: func(m *mock_service.MockMemRepo) {
			},
			wantErr: true,
		},
		{
			name: "error nil value for gauge",
			args: args{metric: models.Metrics{ID: "SomeMetric", MType: models.Gauge}},
			setupMock: func(m *mock_service.MockMemRepo) {
			},
			wantErr: true,
		},
		{
			name: "error nil delta for counter",
			args: args{metric: models.Metrics{ID: "SomeCounter", MType: models.Counter}},
			setupMock: func(m *mock_service.MockMemRepo) {
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(m)
			if err := s.AddStructMetric(tt.args.metric); (err != nil) != tt.wantErr {
				t.Errorf("MemService.AddStructMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemService_GetStructMetric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_service.NewMockMemRepo(ctrl)
	s := NewService(m)

	type args struct {
		metric models.Metrics
	}
	tests := []struct {
		name      string
		args      args
		setupMock func(*mock_service.MockMemRepo)
		want      models.Metrics
		wantErr   bool
	}{
		{
			name: "ok gauge test",
			args: args{metric: models.Metrics{ID: "SomeGauge", MType: models.Gauge}},
			setupMock: func(m *mock_service.MockMemRepo) {
				m.EXPECT().CheckItemGauge("SomeGauge").Return(true).Times(1)
				m.EXPECT().GetItemGauge("SomeGauge").Return("SomeGauge", 123.45).Times(1)
			},
			want:    models.Metrics{ID: "SomeGauge", MType: models.Gauge, Value: func() *float64 { v := 123.45; return &v }()},
			wantErr: false,
		},
		{
			name: "ok counter test",
			args: args{metric: models.Metrics{ID: "SomeCounter", MType: models.Counter}},
			setupMock: func(m *mock_service.MockMemRepo) {
				m.EXPECT().CheckCounter("SomeCounter").Return(true).Times(1)
				m.EXPECT().GetItemCounter("SomeCounter").Return("SomeCounter", int64(456)).Times(1)
			},
			want:    models.Metrics{ID: "SomeCounter", MType: models.Counter, Delta: func() *int64 { v := int64(456); return &v }()},
			wantErr: false,
		},
		{
			name: "error gauge not found",
			args: args{metric: models.Metrics{ID: "NotFoundGauge", MType: models.Gauge}},
			setupMock: func(m *mock_service.MockMemRepo) {
				m.EXPECT().CheckItemGauge("NotFoundGauge").Return(false).Times(1)
				// For gauge metrics that don't exist, we don't call GetItemGauge
			},
			want:    models.Metrics{ID: "NotFoundGauge", MType: models.Gauge},
			wantErr: true,
		},
		{
			name: "error counter not found",
			args: args{metric: models.Metrics{ID: "NotFoundCounter", MType: models.Counter}},
			setupMock: func(m *mock_service.MockMemRepo) {
				m.EXPECT().CheckCounter("NotFoundCounter").Return(false).Times(1)
				// For counter metrics that don't exist, we don't call GetItemCounter
			},
			want:    models.Metrics{ID: "NotFoundCounter", MType: models.Counter},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(m)
			got, err := s.GetStructMetric(tt.args.metric)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemService.GetStructMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.ID != tt.want.ID || got.MType != tt.want.MType {
				t.Errorf("MemService.GetStructMetric() = %v, want %v", got, tt.want)
			}
			if tt.want.Value != nil && got.Value != nil && *got.Value != *tt.want.Value {
				t.Errorf("MemService.GetStructMetric() = %v, want %v", got, tt.want)
			}
			if tt.want.Delta != nil && got.Delta != nil && *got.Delta != *tt.want.Delta {
				t.Errorf("MemService.GetStructMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemService_CheckAddPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_service.NewMockMemRepo(ctrl)
	s := NewService(m)

	tests := []struct {
		name      string
		pathSplit []string
		setupMock func(*mock_service.MockMemRepo)
		wantErr   bool
	}{
		{
			name:      "ok gauge test",
			pathSplit: []string{"localhost:8080", "update", models.Gauge, "SomeMetric", "321"},
			setupMock: func(m *mock_service.MockMemRepo) {
				m.EXPECT().SetGauge("SomeMetric", 321.0).Times(1)
			},
			wantErr: false,
		},
		{
			name:      "ok counter test with existing counter",
			pathSplit: []string{"localhost:8080", "update", models.Counter, "SomeCounter", "123"},
			setupMock: func(m *mock_service.MockMemRepo) {
				m.EXPECT().AddCounter("SomeCounter", int64(123)).Return(true).Times(1)
			},
			wantErr: false,
		},
		{
			name:      "ok counter test with new counter",
			pathSplit: []string{"localhost:8080", "update", models.Counter, "SomeCounter", "123"},
			setupMock: func(m *mock_service.MockMemRepo) {
				m.EXPECT().AddCounter("SomeCounter", int64(123)).Return(false).Times(1)
				m.EXPECT().SetCounter("SomeCounter", int64(123)).Times(1)
			},
			wantErr: false,
		},
		{
			name:      "error invalid metric type",
			pathSplit: []string{"localhost:8080", "update", "invalid", "SomeMetric", "321"},
			setupMock: func(m *mock_service.MockMemRepo) {
			},
			wantErr: true,
		},
		{
			name:      "error wrong path length",
			pathSplit: []string{"localhost:8080", "update", models.Gauge, "SomeMetric"},
			setupMock: func(m *mock_service.MockMemRepo) {
			},
			wantErr: true,
		},
		{
			name:      "error invalid float value",
			pathSplit: []string{"localhost:8080", "update", models.Gauge, "SomeMetric", "invalid"},
			setupMock: func(m *mock_service.MockMemRepo) {
			},
			wantErr: true,
		},
		{
			name:      "error invalid counter value",
			pathSplit: []string{"localhost:8080", "update", models.Counter, "SomeCounter", "invalid"},
			setupMock: func(m *mock_service.MockMemRepo) {
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(m)
			if err := s.CheckAddPath(tt.pathSplit); (err != nil) != tt.wantErr {
				t.Errorf("MemService.CheckAddPath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemService_GetMetric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_service.NewMockMemRepo(ctrl)
	s := NewService(m)

	tests := []struct {
		name      string
		pathSplit []string
		setupMock func(*mock_service.MockMemRepo)
		want      string
		wantErr   bool
	}{
		{
			name:      "ok gauge test",
			pathSplit: []string{"localhost:8080", "value", models.Gauge, "SomeGauge"},
			setupMock: func(m *mock_service.MockMemRepo) {
				m.EXPECT().CheckItemGauge("SomeGauge").Return(true).Times(1)
				m.EXPECT().GetItemGauge("SomeGauge").Return("SomeGauge", 123.45).Times(1)
			},
			want:    "123.45",
			wantErr: false,
		},
		{
			name:      "ok counter test",
			pathSplit: []string{"localhost:8080", "value", models.Counter, "SomeCounter"},
			setupMock: func(m *mock_service.MockMemRepo) {
				m.EXPECT().CheckCounter("SomeCounter").Return(true).Times(1)
				m.EXPECT().GetItemCounter("SomeCounter").Return("SomeCounter", int64(456)).Times(1)
			},
			want:    "456",
			wantErr: false,
		},
		{
			name:      "error invalid metric type",
			pathSplit: []string{"localhost:8080", "value", "invalid", "SomeGauge"},
			setupMock: func(m *mock_service.MockMemRepo) {
				// For invalid metric type, we don't call any methods
			},
			want:    "",
			wantErr: true,
		},
		{
			name:      "error wrong path length",
			pathSplit: []string{"localhost:8080", "value", models.Gauge},
			setupMock: func(m *mock_service.MockMemRepo) {
				// For wrong path length, we don't call any methods
			},
			want:    "",
			wantErr: true,
		},
		{
			name:      "error gauge not found",
			pathSplit: []string{"localhost:8080", "value", models.Gauge, "NotFoundGauge"},
			setupMock: func(m *mock_service.MockMemRepo) {
				m.EXPECT().CheckItemGauge("NotFoundGauge").Return(false).Times(1)
				// For gauge metrics that don't exist, we don't call GetItemGauge
			},
			want:    "",
			wantErr: true,
		},
		{
			name:      "error counter not found",
			pathSplit: []string{"localhost:8080", "value", models.Counter, "NotFoundCounter"},
			setupMock: func(m *mock_service.MockMemRepo) {
				m.EXPECT().CheckCounter("NotFoundCounter").Return(false).Times(1)
				// For counter metrics that don't exist, we don't call GetItemCounter
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(m)
			got, err := s.GetMetric(tt.pathSplit)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemService.GetMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MemService.GetMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}
