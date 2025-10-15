package service

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	models "github.com/maliven1/metrics/internal/model"
	mock_service "github.com/maliven1/metrics/internal/service/mocks"
)

func TestService_CheckAddPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_service.NewMockMemStorage(ctrl)
	s := MemService{memStorage: m}

	m.EXPECT().SetGauge("SomeMetcrics", 321.0).Times(1)

	m.EXPECT().AddCounter("SomeCounter", int64(123)).Return(true).Times(1)

	type args struct {
		pathSplit []string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "ok gauge test",
			args: args{pathSplit: []string{"localhost:8080", "update", models.Gauge, "SomeMetcrics", "321"}},
			want: http.StatusOK,
		}, {
			name: "ok counter test",
			args: args{pathSplit: []string{"localhost:8080", "update", models.Counter, "SomeCounter", "123"}},
			want: http.StatusOK,
		}, {
			name: "StatusBadRequest test - invalid metric type",
			args: args{pathSplit: []string{"localhost:8080", "update", "d", "SomeMetcrics", "321"}},
			want: http.StatusBadRequest,
		}, {
			name: "StatusNotFound test - wrong path length",
			args: args{pathSplit: []string{"localhost:8080", "update", models.Gauge, "SomeMetcrics"}},
			want: http.StatusNotFound,
		}, {
			name: "StatusNotFound test - wrong path length",
			args: args{pathSplit: []string{"localhost:8080", "update", models.Gauge, "321"}},
			want: http.StatusNotFound,
		}, {
			name: "StatusBadRequest test - invalid float value",
			args: args{pathSplit: []string{"localhost:8080", "update", models.Gauge, "SomeMetcrics", "invalid"}},
			want: http.StatusBadRequest,
		}, {
			name: "StatusBadRequest test - invalid counter value",
			args: args{pathSplit: []string{"localhost:8080", "update", models.Counter, "SomeCounter", "invalid"}},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := s.CheckAddPath(tt.args.pathSplit); got != tt.want {
				t.Errorf("Service.CheckAddPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestService_GetMetric(t *testing.T) {
	type args struct {
		pathSplit []string
	}
	tests := []struct {
		name       string
		args       args
		setupMock  func(*mock_service.MockMemStorage)
		want       string
		wantStatus int
	}{
		{
			name: "ok gauge test",
			args: args{pathSplit: []string{"localhost:8080", "value", models.Gauge, "SomeGauge"}},
			setupMock: func(m *mock_service.MockMemStorage) {
				m.EXPECT().GetItemGauge("SomeGauge").Return("SomeGauge", 123.45)
			},
			want:       "123.45",
			wantStatus: http.StatusOK,
		}, {
			name: "ok counter test",
			args: args{pathSplit: []string{"localhost:8080", "value", models.Counter, "SomeCounter"}},
			setupMock: func(m *mock_service.MockMemStorage) {

				m.EXPECT().GetItemGauge("SomeCounter").Return("", 0.0)

				m.EXPECT().GetItemCounter("SomeCounter").Return("SomeCounter", int64(456))
			},
			want:       "456",
			wantStatus: http.StatusOK,
		}, {
			name: "StatusNotFound test - invalid metric type",
			args: args{pathSplit: []string{"localhost:8080", "value", "invalid", "SomeGauge"}},
			setupMock: func(m *mock_service.MockMemStorage) {

				m.EXPECT().GetItemGauge("SomeGauge").Return("", 0.0)

				m.EXPECT().GetItemCounter("SomeGauge").Return("", int64(0))
			},
			want:       "",
			wantStatus: http.StatusNotFound,
		}, {
			name: "StatusNotFound test - wrong path length",
			args: args{pathSplit: []string{"localhost:8080", "value", models.Gauge}},
			setupMock: func(m *mock_service.MockMemStorage) {

			},
			want:       "",
			wantStatus: http.StatusNotFound,
		}, {
			name: "StatusNotFound test - gauge not found",
			args: args{pathSplit: []string{"localhost:8080", "value", models.Gauge, "NotFoundGauge"}},
			setupMock: func(m *mock_service.MockMemStorage) {

				m.EXPECT().GetItemGauge("NotFoundGauge").Return("", 0.0)

				m.EXPECT().GetItemCounter("NotFoundGauge").Return("", int64(0))
			},
			want:       "",
			wantStatus: http.StatusNotFound,
		}, {
			name: "StatusNotFound test - counter not found",
			args: args{pathSplit: []string{"localhost:8080", "value", models.Counter, "NotFoundCounter"}},
			setupMock: func(m *mock_service.MockMemStorage) {

				m.EXPECT().GetItemGauge("NotFoundCounter").Return("", 0.0)

				m.EXPECT().GetItemCounter("NotFoundCounter").Return("", int64(0))
			},
			want:       "",
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_service.NewMockMemStorage(ctrl)
			s := MemService{memStorage: m}

			tt.setupMock(m)

			got, status := s.GetMetric(tt.args.pathSplit)
			if got != tt.want {
				t.Errorf("Service.GetMetric() got = %v, want %v", got, tt.want)
			}
			if status != tt.wantStatus {
				t.Errorf("Service.GetMetric() status = %v, wantStatus %v", status, tt.wantStatus)
			}
		})
	}
}
