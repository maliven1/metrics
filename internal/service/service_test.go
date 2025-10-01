package service

import (
	"testing"

	"github.com/golang/mock/gomock"
	models "github.com/maliven1/metrics/internal/model"
	mock_service "github.com/maliven1/metrics/internal/service/mocks"
)

func TestService_CheckAddPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_service.NewMockMemStorage(ctrl)
	s := Service{memStorage: m}

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
			want: models.StatusOK,
		}, {
			name: "ok counter test",
			args: args{pathSplit: []string{"localhost:8080", "update", models.Counter, "SomeCounter", "123"}},
			want: models.StatusOK,
		}, {
			name: "StatusBadRequest test - invalid metric type",
			args: args{pathSplit: []string{"localhost:8080", "update", "d", "SomeMetcrics", "321"}},
			want: models.StatusBadRequest,
		}, {
			name: "StatusNotFound test - wrong path length",
			args: args{pathSplit: []string{"localhost:8080", "update", models.Gauge, "SomeMetcrics"}},
			want: models.StatusNotFound,
		}, {
			name: "StatusNotFound test - wrong path length",
			args: args{pathSplit: []string{"localhost:8080", "update", models.Gauge, "321"}},
			want: models.StatusNotFound,
		}, {
			name: "StatusBadRequest test - invalid float value",
			args: args{pathSplit: []string{"localhost:8080", "update", models.Gauge, "SomeMetcrics", "invalid"}},
			want: models.StatusBadRequest,
		}, {
			name: "StatusBadRequest test - invalid counter value",
			args: args{pathSplit: []string{"localhost:8080", "update", models.Counter, "SomeCounter", "invalid"}},
			want: models.StatusBadRequest,
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
