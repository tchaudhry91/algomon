package main

import (
	"github.com/tchaudhry91/algoprom/actions"
	"github.com/tchaudhry91/algoprom/algochecks"
)

type Datasource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Config struct {
	Datasources       []Datasource                 `json:"datasources"`
	Algorithmers      []algochecks.AlgorithmerMeta `json:"algorithmers"`
	Actioners         []actions.ActionerMeta       `json:"actioners"`
	Checks            []algochecks.Check           `json:"checks"`
	MetricsListenAddr string                       `json:"metrics_listen_addr"`
	BaseWorkingDir    string                       `json:"base_working_dir"`
}

func fetchDatasourceByName(c *Config, name string) *Datasource {
	for _, d := range c.Datasources {
		if d.Name == name {
			return &d
		}
	}
	return nil
}
