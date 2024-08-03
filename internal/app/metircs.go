package app

import (
	"context"
	"log"
	"runtime"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

// Metrics represents the metrics of the CRI agent.
// It implements prometheus.Collector.
type Metrics struct {
	runtimeVersion   *prometheus.Desc
	runtimeStatus    *prometheus.Desc
	runtimeConfig    *prometheus.Desc
	runtimePodsCount *prometheus.Desc
	imagesCount      *prometheus.Desc
	containersCount  *prometheus.Desc
	fsInfoUsedBytes  *prometheus.Desc
	fsInfoUsedInodes *prometheus.Desc
	agent            Agent
}

// NewMetrics creates a new metrics collector.
func NewMetrics(agent Agent, prefix string) *Metrics {
	if prefix != "" {
		prefix = "cri"
	}
	prefix = strings.TrimPrefix(prefix, "_")
	name := func(name string) string {
		return prefix + "_" + name
	}
	return &Metrics{
		runtimeVersion: prometheus.NewDesc(
			name("runtime_version"),
			"The version of the container runtime",
			[]string{"version", "runtime_name", "runtime_version", "runtime_api_version"}, nil),
		runtimeStatus: prometheus.NewDesc(
			name("runtime_status"),
			"The status of the container runtime",
			[]string{"type", "reason"}, nil),
		runtimeConfig: prometheus.NewDesc(
			name("runtime_config"),
			"The configuration of the container runtime",
			[]string{"cgroup_driver"}, nil),
		runtimePodsCount: prometheus.NewDesc(
			name("runtime_pods_count"),
			"The number of pods in the container runtime",
			nil, nil),
		imagesCount: prometheus.NewDesc(
			name("images_count"),
			"The number of images in the container runtime",
			nil, nil),
		containersCount: prometheus.NewDesc(
			name("containers_count"),
			"The number of containers in the container runtime",
			nil, nil),
		fsInfoUsedBytes: prometheus.NewDesc(
			name("fs_info_used_bytes"),
			"The number of used bytes in the filesystem of mountpoint",
			[]string{"mountpoint"}, nil),
		fsInfoUsedInodes: prometheus.NewDesc(
			name("fs_info_used_inodes"),
			"The number of used inodes in the filesystem of mountpoint",
			[]string{"mountpoint"}, nil),
		agent: agent,
	}
}

// Describe sends the super-set of all possible descriptors of metrics.
func (m *Metrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.runtimeVersion
	ch <- m.runtimeStatus
	ch <- m.runtimeConfig
	ch <- m.runtimePodsCount
	ch <- m.imagesCount
	ch <- m.containersCount
	ch <- m.fsInfoUsedBytes
	ch <- m.fsInfoUsedBytes
}

// Collect is called by the Prometheus registry when collecting metrics.
func (m *Metrics) Collect(ch chan<- prometheus.Metric) {
	ctx := context.Background()

	version, err := m.agent.GetRuntimeVersion(ctx)
	if err != nil {
		log.Fatal(err)
	}

	ch <- prometheus.MustNewConstMetric(m.runtimeVersion, prometheus.GaugeValue, 1,
		version.Version, version.RuntimeName, version.RuntimeVersion, version.RuntimeApiVersion)

	status, err := m.agent.GetStatus(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, condition := range status.GetStatus().GetConditions() {
		status := float64(0)
		if condition.GetStatus() {
			status = float64(1)
		}
		ch <- prometheus.MustNewConstMetric(m.runtimeStatus, prometheus.GaugeValue, status, condition.Type, condition.Reason)
	}

	if runtime.GOOS == "linux" {
		config, err1 := m.agent.GetRuntimeConfig(ctx)
		if err1 == nil { // cri-dockerd doesn't support RuntimeConfig
			ch <- prometheus.MustNewConstMetric(m.runtimeConfig, prometheus.GaugeValue, 1, config.GetLinux().GetCgroupDriver().String())
		}
	}

	pods, err := m.agent.ListPods(ctx)
	if err != nil {
		log.Fatal(err)
	}
	ch <- prometheus.MustNewConstMetric(m.runtimePodsCount, prometheus.GaugeValue, float64(len(pods.Items)))

	image, err := m.agent.ListImages(ctx)
	if err != nil {
		log.Fatal(err)
	}

	ch <- prometheus.MustNewConstMetric(m.imagesCount, prometheus.GaugeValue, float64(len(image.Images)))

	containers, err := m.agent.ListContainers(ctx)
	if err != nil {
		log.Fatal(err)
	}

	ch <- prometheus.MustNewConstMetric(m.containersCount, prometheus.GaugeValue, float64(len(containers.Containers)))

	fs, err := m.agent.FsInfo(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, fsInfo := range fs.ImageFilesystems {
		mountpoint := fsInfo.GetFsId().GetMountpoint()
		usedBytes := fsInfo.GetUsedBytes().GetValue()
		ch <- prometheus.MustNewConstMetric(m.fsInfoUsedBytes, prometheus.GaugeValue, float64(usedBytes), mountpoint)
		inodesUsed := fsInfo.GetInodesUsed().GetValue()
		ch <- prometheus.MustNewConstMetric(m.fsInfoUsedInodes, prometheus.GaugeValue, float64(inodesUsed), mountpoint)
	}
}
