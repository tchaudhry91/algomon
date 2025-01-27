package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	countProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "algoprom_count_processed_total",
		Help: "The total number of measurements checked performed",
	})

	countSuccess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "algoprom_count_success_total",
		Help: "The total number of measurements that succeeded",
	})

	countFail = promauto.NewCounter(prometheus.CounterOpts{
		Name: "algoprom_count_fail_total",
		Help: "The total number of measurements that failed",
	})
)
