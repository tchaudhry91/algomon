package main

import (
	"github.com/tchaudhry91/algoprom/algochecks"
)

type Datasource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Config struct {
	Datasources       []Datasource                 `json:"datasources"`
	Algorithmers      []algochecks.AlgorithmerMeta `json:"algorithmers"`
	Checks            []algochecks.Check           `json:"checks"`
	MetricsListenAddr string                       `json:"metrics_listen_addr"`
	BaseWorkingDir    string                       `json:"base_working_dir"`
}
