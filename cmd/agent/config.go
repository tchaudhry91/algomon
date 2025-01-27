package main

import (
	"github.com/tchaudhry91/algoprom/alerts"
)

type Measurement struct {
	Name          string           `json:"name"`
	Datasource    string           `json:"datasource"`
	Interval      Duration         `json:"interval"`
	Query         string           `json:"query"`
	Algorithms    string           `json:"algorithms"`
	AlertChannels []alerts.Alerter `json:"alert_channels"`
	Immediate     bool             `json:"immediate"`
}

type Datasource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Config struct {
	Datasources       []Datasource  `json:"datasources"`
	Algorithms        string        `json:"algorithms"`
	Measurements      []Measurement `json:"measurements"`
	MetricsListenAddr string        `json:"metrics_listen_addr"`
}
