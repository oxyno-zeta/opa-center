// +build integration

package server

import (
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/metrics"
)

// Generate metrics instance
var metricsCtx = metrics.NewMetricsClient()
