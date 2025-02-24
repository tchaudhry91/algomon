package algochecks

import (
	"github.com/tchaudhry91/algomon/actions"
	"github.com/tchaudhry91/algomon/measure"
)

type Check struct {
	Name            string                `json:"name"`
	Inputs          []measure.Measurement `json:"inputs"`
	AlgorithmerType string                `json:"algorithmer_type"`
	Algorithm       string                `json:"algorithm"`
	AlgorithmParams map[string]string     `json:"algorithm_params"`
	Actions         []actions.ActionMeta  `json:"actions"`
	Interval        measure.Duration      `json:"interval"`
	Immediate       bool                  `json:"immediate"`
	Debug           bool                  `json:"debug"`
}
