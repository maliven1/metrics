// Package agent
package agent

import (
	"log"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/maliven1/metrics/internal/config"
	models "github.com/maliven1/metrics/internal/model"
	"github.com/shirou/gopsutil/v4/cpu"
	gopsutil "github.com/shirou/gopsutil/v4/mem"

	_ "net/http/pprof"
)

type MemStorage interface {
	SetGauge(key string, value float64)
	SetCounter(key string, value int64)
	GetGauge() map[string]float64
	GetCounter() map[string]int64
	AddCounter(key string, value int64) bool
}

type Agent struct {
	memStorage MemStorage
	cfg        *config.AgentConfig
	m          *sync.Mutex
}

func NewAgent(m MemStorage, cfg *config.AgentConfig) *Agent {
	var mutex sync.Mutex
	return &Agent{memStorage: m, cfg: cfg, m: &mutex}
}

func (a Agent) GetMetrics() (map[string]float64, map[string]int64) {

	return a.memStorage.GetGauge(), a.memStorage.GetCounter()
}

func (a Agent) CollectMetrics() {
	for {
		a.addMetrics()
		err := a.gopsutilAddMetrics()
		if err != nil {
			log.Printf("error: %v", err)
		}
		time.Sleep(time.Duration(a.cfg.PollInterval) * time.Second)
	}

}
func (a Agent) gopsutilAddMetrics() error {
	a.m.Lock()
	defer a.m.Unlock()
	v, err := gopsutil.VirtualMemory()
	if err != nil {
		log.Printf("error: %v", err)
		return err
	}

	a.memStorage.SetGauge(models.TotalMemory, float64(v.Total))
	a.memStorage.SetGauge(models.FreeMemory, float64(v.Free))
	info, err := cpu.Info()
	if err != nil {
		log.Printf("error: %v", err)
		return err
	}
	numCPU := len(info)
	a.memStorage.SetGauge(models.CPUutilization1, float64(numCPU))
	return nil
}
func (a Agent) addMetrics() {
	a.m.Lock()
	defer a.m.Unlock()

	const count int64 = 1

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	a.memStorage.SetGauge(models.Alloc, float64(mem.Alloc))
	a.memStorage.SetGauge(models.BuckHashSys, float64(mem.BuckHashSys))
	a.memStorage.SetGauge(models.Frees, float64(mem.Frees))
	a.memStorage.SetGauge(models.GCCPUFraction, mem.GCCPUFraction)
	a.memStorage.SetGauge(models.GCSys, float64(mem.GCSys))
	a.memStorage.SetGauge(models.HeapAlloc, float64(mem.HeapAlloc))
	a.memStorage.SetGauge(models.HeapIdle, float64(mem.HeapIdle))
	a.memStorage.SetGauge(models.HeapInuse, float64(mem.HeapInuse))
	a.memStorage.SetGauge(models.HeapObjects, float64(mem.HeapObjects))
	a.memStorage.SetGauge(models.HeapReleased, float64(mem.HeapReleased))
	a.memStorage.SetGauge(models.HeapSys, float64(mem.HeapSys))
	a.memStorage.SetGauge(models.LastGC, float64(mem.LastGC))
	a.memStorage.SetGauge(models.Lookups, float64(mem.Lookups))
	a.memStorage.SetGauge(models.MCacheInuse, float64(mem.MCacheInuse))
	a.memStorage.SetGauge(models.MCacheSys, float64(mem.MCacheSys))
	a.memStorage.SetGauge(models.MSpanInuse, float64(mem.MSpanInuse))
	a.memStorage.SetGauge(models.MSpanSys, float64(mem.MSpanSys))
	a.memStorage.SetGauge(models.Mallocs, float64(mem.Mallocs))
	a.memStorage.SetGauge(models.NextGC, float64(mem.NextGC))
	a.memStorage.SetGauge(models.NumForcedGC, float64(mem.NumForcedGC))
	a.memStorage.SetGauge(models.NumGC, float64(mem.NumGC))
	a.memStorage.SetGauge(models.OtherSys, float64(mem.OtherSys))
	a.memStorage.SetGauge(models.PauseTotalNs, float64(mem.PauseTotalNs))
	a.memStorage.SetGauge(models.StackInuse, float64(mem.StackInuse))
	a.memStorage.SetGauge(models.StackSys, float64(mem.StackSys))
	a.memStorage.SetGauge(models.Sys, float64(mem.Sys))
	a.memStorage.SetGauge(models.TotalAlloc, float64(mem.TotalAlloc))
	a.memStorage.SetGauge(models.RandomValue, rand.Float64())
	if a.memStorage.AddCounter(models.PollCount, count) {
		return
	}
	a.memStorage.SetCounter(models.PollCount, count)
}
