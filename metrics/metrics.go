package metrics

import (
	"docker-stats-exporter/docker"
	"fmt"
	"net/http"
	"sync"
)

func cpuPercent(stats *docker.Stats) float64 {
	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage - stats.PreCPUStats.CPUUsage.TotalUsage)
	sysDelta := float64(stats.CPUStats.SystemUsage - stats.PreCPUStats.SystemUsage)
	if sysDelta > 0 && cpuDelta > 0 {
		return cpuDelta / sysDelta
	}
	return 0
}

func memoryPercent(stats *docker.Stats) float64 {
	return float64(stats.MemoryStats.Usage) / float64(stats.MemoryStats.Limit)
}

func format(stats *docker.Stats, container docker.Container) string {
	name := container.Names[0][1:]

	metrics := ""
	metrics += fmt.Sprintf(
		"docker_container_memory_usage_bytes{name=%q} %d\n",
		name, stats.MemoryStats.Usage)

	metrics += fmt.Sprintf(
		"docker_container_memory_limit_bytes{name=%q} %d\n",
		name, stats.MemoryStats.Limit)

	metrics += fmt.Sprintf(
		"docker_container_memory_percent{name=%q} %.2f\n",
		name, memoryPercent(stats))

	metrics += fmt.Sprintf(
		"docker_container_cpu_percent{name=%q} %.2f\n",
		name, cpuPercent(stats))

	return metrics
}

func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	client := docker.DockerClient()
	containers, err := docker.GetContainers(client)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	metrics := make(chan string, len(containers))
	var wg sync.WaitGroup

	for _, c := range containers {
		wg.Add(1)
		go func(container docker.Container) {
			defer wg.Done()

			stats, err := docker.GetStats(client, container.ID)
			if err != nil {
				return
			}

			metrics <- format(stats, container)
		}(c)
	}

	wg.Wait()
	close(metrics)

	for metric := range metrics {
		fmt.Fprint(w, metric)
	}
}
