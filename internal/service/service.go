package service

import (
	"strconv"

	models "github.com/maliven1/metrics/internal/model"
)

func CheckPath(pathSplit []string, memStorage models.MemStorage) int {

	if len(pathSplit) != 5 {
		return models.StatusNotFound
	}
	if float, err := strconv.ParseFloat(pathSplit[4], 64); pathSplit[2] == models.Gauge && err == nil && float != 0 {
		memStorage.Gauge[pathSplit[3]] = float
		return models.StatusOK
	} else if count, err := strconv.Atoi(pathSplit[4]); pathSplit[2] == models.Counter && err == nil && count != 0 {
		_, ok := memStorage.Counter[pathSplit[3]]
		if ok {
			memStorage.Counter[pathSplit[3]] += int64(count)
			return models.StatusOK
		}
		memStorage.Counter[pathSplit[3]] = int64(count)
		return models.StatusOK
	} else {
		return models.StatusBadRequest
	}
}
