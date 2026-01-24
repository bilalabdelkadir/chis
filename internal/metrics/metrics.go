package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_ms",
			Help:    "HTTP request duration in milliseconds",
			Buckets: []float64{5, 10, 25, 50, 100, 250, 500, 1000},
		},
		[]string{"method", "path"},
	)

	WebhooksDeliveredTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "webhooks_delivered_total",
			Help: "Total webhooks delivered",
		},
		[]string{"status"},
	)

	WebhookDeliveryDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "webhook_delivery_duration_ms",
			Help:    "Webhook delivery duration in milliseconds",
			Buckets: []float64{100, 250, 500, 1000, 2500, 5000, 10000},
		},
	)
)
