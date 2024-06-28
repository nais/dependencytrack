package observability

import (
	"github.com/prometheus/client_golang/prometheus"
)

var DependencytrackTotalProjects = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "dependencytrack_total_projects",
		Help: "Total number of projects in DependencyTrack",
	},
	[]string{"cluster", "namespace"},
)

func init() {
	prometheus.MustRegister(DependencytrackTotalProjects)
}
