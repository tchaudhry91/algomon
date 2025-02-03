package algochecks

import (
	"context"
	"log"

	"github.com/tchaudhry91/algoprom/measure"
)

// Basic Python Algorithm Apply
type PythonAlgorithmer struct {
	VEnv      string `json:"venv"`
	Directory string `json:"directory"`
	logger    *log.Logger
}

func (pa *PythonAlgorithmer) ApplyAlgorithm(ctx context.Context, algorithm string, algorithmParams map[string]string, inputs []measure.Measurement, workingDir string) error {
	return nil
}
