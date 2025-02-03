package algochecks

import (
	"context"
	"log"

	"github.com/tchaudhry91/algoprom/measure"
)

type Output struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	RC     int    `json:"rc"`
}

type AlgorithmerMeta struct {
	Type   string            `json:"type"`
	Params map[string]string `json:"params"`
}

type Algorithmer interface {
	ApplyAlgorithm(ctx context.Context, algorithm string, algorithmParams map[string]string, inputs []measure.Measurement, workingDir string) error
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
