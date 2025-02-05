package actions

import (
	"context"
	"log"
	"time"
)

type Output struct {
	Timestamp   time.Time `json:"timestamp"`
	CombinedOut string    `json:"combined_out"`
	RC          int       `json:"rc"`
	Error       error     `json:"error"`
}

type ActionerMeta struct {
	Type        string            `json:"type"`
	Params      map[string]string `json:"params"`
	EnvOverride map[string]string `json:"env_override"`
}

type Actioner interface {
	Action(ctx context.Context, action string, params map[string]string) (Output, error)
}

func Build(meta ActionerMeta, logger *log.Logger) Actioner {
	if meta.Type == "python" {
		return &PythonActioner{
			VEnv:        meta.Params["venv"],
			Directory:   meta.Params["directory"],
			EnvOverride: meta.EnvOverride,
			logger:      logger,
		}
	}
	return nil
}
