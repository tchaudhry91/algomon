package measure

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type Measurement struct {
	Name       string `json:"name"`
	Datasource string `json:"datasource"`
	Query      string `json:"query"`
}

func (m *Measurement) Measure(ctx context.Context) (map[string]string, error) {
	apiClient, err := getPromAPIClient(m.Datasource)
	if err != nil {
		return nil, err
	}
	res, _, err := apiClient.Query(ctx, m.Query, time.Now())
	if err != nil {
		return nil, err
	}
	results := map[string]string{}
	switch res.Type() {
	case model.ValVector:
		vector := res.(model.Vector)
		for _, sample := range vector {
			results[sample.Metric.String()] = sample.Value.String()
		}
	}
	return results, nil
}

func getPromAPIClient(datasourceURL string) (v1.API, error) {
	client, err := api.NewClient(api.Config{Address: datasourceURL})
	if err != nil {
		return nil, err
	}
	api := v1.NewAPI(client)
	return api, nil
}
