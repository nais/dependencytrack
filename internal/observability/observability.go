package observability

import (
	"github.com/prometheus/client_golang/prometheus"
)

var DependencytrackTotalProjects = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "dependencytrack_projects_total",
		Help: "Total number of projects in DependencyTrack",
	},
	[]string{"cluster", "team", "workload", "has_sbom"},
)

var DependencytrackTotalPlatformProjects = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "dependencytrack_projects_platform_total",
		Help: "Total number of platform projects in DependencyTrack",
	},
	[]string{"cluster", "team", "workload", "platform"},
)

func init() {
	prometheus.MustRegister(DependencytrackTotalProjects)
	prometheus.MustRegister(DependencytrackTotalPlatformProjects)
}
