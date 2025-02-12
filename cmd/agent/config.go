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
	Datasources    []Datasource                 `json:"datasources"`
	Algorithmers   []algochecks.AlgorithmerMeta `json:"algorithmers"`
	Actioners      []actions.ActionerMeta       `json:"actioners"`
	Checks         []algochecks.Check           `json:"checks"`
	BaseWorkingDir string                       `json:"base_working_dir"`
	DatabaseFile   string                       `json:"database_file"`
	APIListenAddr  string                       `json:"api_listen_addr"`
}

func fetchDatasourceByName(c *Config, name string) *Datasource {
	for _, d := range c.Datasources {
		if d.Name == name {
			return &d
		}
	}
	return nil
}
