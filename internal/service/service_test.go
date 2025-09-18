package service

import (
	"testing"

	models "github.com/maliven1/metrics/internal/model"
	"github.com/maliven1/metrics/internal/repository"
	"github.com/maliven1/metrics/internal/storage"
)

func TestService_CheckAddPath(t *testing.T) {
	memStorage := storage.NewMemStorage()
	cache := repository.NewCache(memStorage)
	s := NewService(cache)
	type args struct {
		pathSplit []string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "ok test",
			args: args{pathSplit: []string{"localhost:8080", "update", models.Gauge, "SomeMetcrics", "321"}},
			want: models.StatusOK,
		}, {
			name: "StatusBadRequest test",
			args: args{pathSplit: []string{"localhost:8080", "update", "d", "SomeMetcrics", "321"}},
			want: models.StatusBadRequest,
		}, {
			name: "StatusNotFound test",
			args: args{pathSplit: []string{"localhost:8080", "update", models.Gauge, "SomeMetcrics"}},
			want: models.StatusNotFound,
		}, {
			name: "StatusNotFound test",
			args: args{pathSplit: []string{"localhost:8080", "update", models.Gauge, "321"}},
			want: models.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := s.CheckAddPath(tt.args.pathSplit); got != tt.want {
				t.Errorf("Service.CheckPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
