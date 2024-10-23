package observability

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	Namespace = "dtrack"
)

var labels = []string{"cluster", "workload_namespace", "project", "workload_type", "workload", "has_sbom"}

var WorkloadRiskscore = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "workloads_riskscore",
	},
	labels,
)

var WorkloadCritical = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "workloads_critical",
	},
	labels,
)

func init() {
	prometheus.MustRegister(WorkloadRiskscore)
	prometheus.MustRegister(WorkloadCritical)
}
