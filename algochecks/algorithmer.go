package algochecks

import (
	"context"
	"log"
	"time"

	"github.com/tchaudhry91/algoprom/measure"
)

type Output struct {
	Timestamp   time.Time `json:"timestamp"`
	CombinedOut string    `json:"combined_out"`
	RC          int       `json:"rc"`
	Error       error     `json:"error"`
}

type AlgorithmerMeta struct {
	Type   string            `json:"type"`
	Params map[string]string `json:"params"`
}

type Algorithmer interface {
	ApplyAlgorithm(ctx context.Context, algorithm string, algorithmParams map[string]string, inputs []measure.Measurement, workingDir string) (Output, error)
}

func Build(aType string, params map[string]string, logger *log.Logger) Algorithmer {
	if aType == "python" {
		return &PythonAlgorithmer{
			VEnv:      params["venv"],
			Directory: params["directory"],
			logger:    logger,
		}
	}
	return nil
}
