package measure

import (
	"context"
	"log"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/tchaudhry91/algoprom/alerts"
)

type Measurement struct {
	Name          string           `json:"name"`
	Datasource    string           `json:"datasource"`
	Interval      Duration         `json:"interval"`
	Query         string           `json:"query"`
	Algorithms    []string         `json:"algorithms"`
	AlertChannels []alerts.Alerter `json:"alert_channels"`
	Immediate     bool             `json:"immediate"`
	Debug         bool             `json:"debug"`
}

func (m *Measurement) Measure(ctx context.Context, logger *log.Logger, datasourceURL, tempWorkDir string) error {
	apiClient, err := getPromAPIClient(datasourceURL)
	if err != nil {
		return err
	}
	res, _, err := apiClient.Query(ctx, m.Query, time.Now())
	if err != nil {
		return err
	}
	switch res.Type() {
	case model.ValVector:
		vector := res.(model.Vector)
		for _, sample := range vector {
			logger.Printf("%s-%s", sample.Metric, sample.Value)
		}
	}
	return nil
}

func getPromAPIClient(datasourceURL string) (v1.API, error) {
	client, err := api.NewClient(api.Config{Address: datasourceURL})
	if err != nil {
		return nil, err
	}
	api := v1.NewAPI(client)
	return api, nil
}
