package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	gauge   = "gauge"
	counter = "counter"
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	if r.Method == http.MethodPost {
		memStorage := MemStorage{gauge: make(map[string]float64), counter: make(map[string]int64)}
		pathSplit := strings.Split(r.URL.Path, "/")
		if float, err := strconv.ParseFloat(pathSplit[4], 64); pathSplit[2] == gauge && err == nil && float != 0 && pathSplit[3] != "" {
			memStorage.gauge[pathSplit[3]] = float
			w.WriteHeader(http.StatusOK)
			return
		} else if count, err := strconv.Atoi(pathSplit[4]); pathSplit[2] == counter && err == nil && count != 0 && pathSplit[3] != "" {
			_, ok := memStorage.counter[pathSplit[3]]
			if ok {
				memStorage.counter[pathSplit[3]] += int64(count)
				w.WriteHeader(http.StatusOK)
				return
			}
			memStorage.counter[pathSplit[3]] = int64(count)
			w.WriteHeader(http.StatusOK)
			return
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	} else {

		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, mainPage)
	log.Println("serv start")
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
