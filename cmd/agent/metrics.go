package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	countProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "algoprom_count_processed_total",
		Help: "The total number of measurements checked performed",
	}, []string{"measurement"})

	countSuccess = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "algoprom_count_success_total",
		Help: "The total number of measurements that succeeded",
	}, []string{"measurement"})

	countFail = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "algoprom_count_fail_total",
		Help: "The total number of measurements that failed",
	}, []string{"measurement"})
)
