package algochecks

import (
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
	ApplyAlgorithm(algorithm string, algorithmParams map[string]string, inputs []measure.Measurement, workingDir string) error
}

func Build(aType string, params map[string]string) Algorithmer {
	if aType == "python" {
		return &PythonAlgorithmer{
			VEnv:      params["venv"],
			Directory: params["directory"],
		}
	}
	return nil
}

// Basic Python Algorithm Apply
type PythonAlgorithmer struct {
	VEnv      string `json:"venv"`
	Directory string `json:"directory"`
}

func (pa *PythonAlgorithmer) ApplyAlgorithm(algorithm string, algorithmParams map[string]string, inputs []measure.Measurement, workingDir string) error {
	return nil
}
