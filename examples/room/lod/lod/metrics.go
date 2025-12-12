// © 2019-2025 Diarkis Inc. All rights reserved.

package lod

import (
	"github.com/Diarkis/diarkis"
	"github.com/Diarkis/diarkis/metrics"
)

var (
	// MetricActiveUsersGaugeOpts tracks the number of users managed by LOD manager
	MetricActiveUsersGaugeOpts = metrics.GaugeVecOpts{
		GaugeOpts: metrics.GaugeOpts{
			Name:      "lod_manager_active_users",
			Help:      "Number of user entities currently managed by the LOD manager.",
			NodeRoles: []diarkis.NodeRole{diarkis.UDP, diarkis.TCP},
		},
		LabelNames: []string{"room_id"},
	}

	// MetricLoopDurationOpts measures LOD loop processing time
	MetricLoopDurationOpts = metrics.HistogramOpts{
		Name:      "lod_manager_loop_duration_seconds",
		Help:      "Duration of processAllUsers execution in seconds.",
		NodeRoles: []diarkis.NodeRole{diarkis.UDP, diarkis.TCP},
		Buckets:   []float64{0.001, 0.005, 0.010, 0.025, 0.050, 0.100, 0.250, 0.500, 1.0},
	}

	// MetricUpdatesSentOpts counts updates sent to users
	MetricUpdatesSentOpts = metrics.CounterVecOpts{
		CounterOpts: metrics.CounterOpts{
			Name:      "lod_manager_updates_sent_total",
			Help:      "Total number of LOD updates sent to users.",
			NodeRoles: []diarkis.NodeRole{diarkis.UDP, diarkis.TCP},
		},
		LabelNames: []string{"distance_category", "reason"},
	}

	// MetricBufferDropsOpts counts messages dropped due to full buffer
	MetricBufferDropsOpts = metrics.CounterOpts{
		Name:      "lod_manager_send_buffer_drops_total",
		Help:      "Total number of messages dropped due to full send buffer.",
		NodeRoles: []diarkis.NodeRole{diarkis.UDP, diarkis.TCP},
	}

	// MetricBufferSizeGaugeOpts tracks current send buffer size
	MetricBufferSizeGaugeOpts = metrics.GaugeOpts{
		Name:      "lod_manager_send_buffer_size",
		Help:      "Current number of messages queued in the send buffer.",
		NodeRoles: []diarkis.NodeRole{diarkis.UDP, diarkis.TCP},
	}

	// MetricUpdatesSkippedOpts counts updates that were skipped
	MetricUpdatesSkippedOpts = metrics.CounterVecOpts{
		CounterOpts: metrics.CounterOpts{
			Name:      "lod_manager_updates_skipped_total",
			Help:      "Total number of LOD updates skipped.",
			NodeRoles: []diarkis.NodeRole{diarkis.UDP, diarkis.TCP},
		},
		LabelNames: []string{"reason"},
	}

	// MetricUserDistancesOpts tracks distribution of user distances
	MetricUserDistancesOpts = metrics.HistogramOpts{
		Name:      "lod_manager_user_distances",
		Help:      "Distribution of distances between users.",
		NodeRoles: []diarkis.NodeRole{diarkis.UDP, diarkis.TCP},
		Buckets:   []float64{100, 500, 1000, 2000, 5000, 10000, 20000, 50000},
	}

	// MetricMessagesSentOpts counts messages actually sent by sendWorker
	MetricMessagesSentOpts = metrics.CounterOpts{
		Name:      "lod_manager_messages_sent_total",
		Help:      "Total number of messages successfully sent to users by sendWorker.",
		NodeRoles: []diarkis.NodeRole{diarkis.UDP, diarkis.TCP},
	}

	// MetricInstancesGaugeOpts tracks number of active LOD manager instances
	MetricInstancesGaugeOpts = metrics.GaugeOpts{
		Name:      "lod_manager_instances",
		Help:      "Number of active LOD manager instances.",
		NodeRoles: []diarkis.NodeRole{diarkis.UDP, diarkis.TCP},
	}

	// MetricErrorsOpts counts errors by type
	MetricErrorsOpts = metrics.CounterVecOpts{
		CounterOpts: metrics.CounterOpts{
			Name:      "lod_manager_errors_total",
			Help:      "Total number of errors encountered by LOD manager.",
			NodeRoles: []diarkis.NodeRole{diarkis.UDP, diarkis.TCP},
		},
		LabelNames: []string{"error_type"},
	}
)

var (
	activeUsersGauge       *metrics.GaugeVec
	loopDurationHistogram  *metrics.Histogram
	updatesSentCounter     *metrics.CounterVec
	bufferDropsCounter     *metrics.Counter
	bufferSizeGauge        *metrics.Gauge
	updatesSkippedCounter  *metrics.CounterVec
	userDistancesHistogram *metrics.Histogram
	messagesSentCounter    *metrics.Counter
	instancesGauge         *metrics.Gauge
	errorsCounter          *metrics.CounterVec
)

func SetupLodMetrics() {
	activeUsersGauge = metrics.NewGaugeVec(MetricActiveUsersGaugeOpts)
	loopDurationHistogram = metrics.NewHistogram(MetricLoopDurationOpts)
	updatesSentCounter = metrics.NewCounterVec(MetricUpdatesSentOpts)
	bufferDropsCounter = metrics.NewCounter(MetricBufferDropsOpts)
	bufferSizeGauge = metrics.NewGauge(MetricBufferSizeGaugeOpts)
	updatesSkippedCounter = metrics.NewCounterVec(MetricUpdatesSkippedOpts)
	userDistancesHistogram = metrics.NewHistogram(MetricUserDistancesOpts)
	messagesSentCounter = metrics.NewCounter(MetricMessagesSentOpts)
	instancesGauge = metrics.NewGauge(MetricInstancesGaugeOpts)
	errorsCounter = metrics.NewCounterVec(MetricErrorsOpts)
}
