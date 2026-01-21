package monitor

import (
	"OpsGo/internal/infrastructure/redis"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

const (
	MetricCPU    = "metrics:cpu"
	MetricMemory = "metrics:mem"
	MaxPoints    = 100 // Keep last 100 points for real-time VIEW
)

type SystemMetric struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

type MonitorService struct {
	stopChan chan struct{}
}

func NewMonitorService() *MonitorService {
	return &MonitorService{
		stopChan: make(chan struct{}),
	}
}

func (s *MonitorService) StartCollector() {
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.collect()
			case <-s.stopChan:
				return
			}
		}
	}()
	log.Println("Monitor Collector started")
}

func (s *MonitorService) StopCollector() {
	close(s.stopChan)
}

func (s *MonitorService) collect() {
	ctx := context.Background()
	now := time.Now().UnixMilli()

	// CPU
	percent, err := cpu.Percent(0, false)
	if err == nil && len(percent) > 0 {
		s.saveMetric(ctx, MetricCPU, now, percent[0])
	}

	// Memory
	v, err := mem.VirtualMemory()
	if err == nil {
		s.saveMetric(ctx, MetricMemory, now, v.UsedPercent)
	}
}

func (s *MonitorService) saveMetric(ctx context.Context, key string, ts int64, val float64) {
	metric := SystemMetric{
		Timestamp: ts,
		Value:     val,
	}
	data, _ := json.Marshal(metric)

	pipe := redis.Client.Pipeline()
	pipe.RPush(ctx, key, data)
	pipe.LTrim(ctx, key, -MaxPoints, -1) // Keep last N points
	_, err := pipe.Exec(ctx)

	if err != nil {
		log.Printf("Failed to save metric %s: %v", key, err)
	}
}

func (s *MonitorService) GetMetrics(ctx context.Context) (map[string][]SystemMetric, error) {
	pipe := redis.Client.Pipeline()
	cpuCmd := pipe.LRange(ctx, MetricCPU, 0, -1)
	memCmd := pipe.LRange(ctx, MetricMemory, 0, -1)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("redis error: %v", err)
	}

	return map[string][]SystemMetric{
		"cpu":    parseMetrics(cpuCmd.Val()),
		"memory": parseMetrics(memCmd.Val()),
	}, nil
}

func parseMetrics(raw []string) []SystemMetric {
	res := make([]SystemMetric, 0, len(raw))
	for _, s := range raw {
		var m SystemMetric
		if err := json.Unmarshal([]byte(s), &m); err == nil {
			res = append(res, m)
		}
	}
	return res
}
