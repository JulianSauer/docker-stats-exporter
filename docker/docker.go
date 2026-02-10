package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

type Container struct {
	ID    string   `json:"Id"`
	Names []string `json:"Names"`
}

type Stats struct {
	CPUStats struct {
		CPUUsage struct {
			TotalUsage uint64 `json:"total_usage"`
		} `json:"cpu_usage"`
		SystemUsage uint64 `json:"system_cpu_usage"`
	} `json:"cpu_stats"`

	PreCPUStats struct {
		CPUUsage struct {
			TotalUsage uint64 `json:"total_usage"`
		} `json:"cpu_usage"`
		SystemUsage uint64 `json:"system_cpu_usage"`
	} `json:"precpu_stats"`

	MemoryStats struct {
		Usage uint64 `json:"usage"`
		Limit uint64 `json:"limit"`
	} `json:"memory_stats"`
}

func DockerClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return net.Dial("unix", "/var/run/docker.sock")
			},
		},
		Timeout: 5 * time.Second,
	}
}

func GetContainers(client *http.Client) ([]Container, error) {
	resp, err := client.Get("http://unix/containers/json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var c []Container
	return c, json.NewDecoder(resp.Body).Decode(&c)
}

func GetStats(client *http.Client, id string) (*Stats, error) {
	url := fmt.Sprintf("http://unix/containers/%s/stats?stream=false", id)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s Stats
	return &s, json.NewDecoder(resp.Body).Decode(&s)
}
