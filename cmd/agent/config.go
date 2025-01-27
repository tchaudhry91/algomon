package main

import "github.com/tchaudhry91/algoprom/measure"

type Datasource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Config struct {
	Datasources       []Datasource          `json:"datasources"`
	Algorithms        string                `json:"algorithms"`
	Measurements      []measure.Measurement `json:"measurements"`
	MetricsListenAddr string                `json:"metrics_listen_addr"`
	BaseWorkingDir    string                `json:"base_working_dir"`
}
