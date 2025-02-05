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
	Type        string            `json:"type"`
	Params      map[string]string `json:"params"`
	EnvOverride map[string]string `json:"env_override"`
}

type Algorithmer interface {
	ApplyAlgorithm(ctx context.Context, algorithm string, algorithmParams map[string]string, inputs map[string]measure.Result, workingDir string) (Output, error)
}

func Build(meta AlgorithmerMeta, logger *log.Logger) Algorithmer {
	if meta.Type == "python" {
		return &PythonAlgorithmer{
			VEnv:        meta.Params["venv"],
			Directory:   meta.Params["directory"],
			EnvOverride: meta.EnvOverride,
			logger:      logger,
		}
	}
	return nil
}
