package measure_test

import (
	"context"
	"testing"

	"github.com/tchaudhry91/algoprom/measure"
)

func TestMeasure(t *testing.T) {
	m := measure.Measurement{
		"Test",
		"http://demo.robustperception.io:9090",
		"sum(increase(caddy_http_requests_total[5m]))",
	}

	res, err := m.Measure(context.Background())
	if err != nil {
		t.Errorf("Failed measurement test with: %v", err)
		t.FailNow()
	}
	t.Logf("Results:%v", res)

}
