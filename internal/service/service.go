package service

import (
	"strconv"

	models "github.com/maliven1/metrics/internal/model"
	"github.com/maliven1/metrics/internal/repository"
)

func CheckPath(pathSplit []string, memStorage *repository.MemStorage) int {

	if len(pathSplit) != 5 {
		return models.StatusNotFound
	}
	if float, err := strconv.ParseFloat(pathSplit[4], 64); pathSplit[2] == models.Gauge && err == nil && float != 0 {
		memStorage.SetGauge(pathSplit[3], float)
		return models.StatusOK
	} else if count, err := strconv.Atoi(pathSplit[4]); pathSplit[2] == models.Counter && err == nil && count != 0 {
		if memStorage.CheckCounter(pathSplit[3]) {
			memStorage.AddCounter(pathSplit[3], int64(count))
			return models.StatusOK
		}
		memStorage.SetCounter(pathSplit[3], int64(count))
		return models.StatusOK
	} else {
		return models.StatusBadRequest
	}
}
