package observability

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	Namespace = "vuln"
)

var Labels = []string{"cluster", "team", "has_sbom", "image", "version"}

var WorkloadRegistered = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "workloads_registered",
	},
	[]string{"cluster", "team", "workload", "has_sbom"},
)

var WorkloadRiskscore = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "workloads_riskscore",
	},
	Labels,
)

var WorkloadCritical = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "workloads_critical",
	},
	Labels,
)

func init() {
	prometheus.MustRegister(WorkloadRegistered)
	prometheus.MustRegister(WorkloadRiskscore)
	prometheus.MustRegister(WorkloadCritical)
}
