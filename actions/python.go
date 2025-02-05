package actions

import (
	"context"
	"log"
	"time"
)

// Basic Python Actioner
type PythonActioner struct {
	VEnv        string            `json:"venv"`
	Directory   string            `json:"directory"`
	EnvOverride map[string]string `json:"env_override"`
	logger      *log.Logger
}

func (pa *PythonActioner) Action(ctx context.Context, action string, params map[string]string) (Output, error) {
	out := Output{
		RC:          -1,
		CombinedOut: "",
		Timestamp:   time.Now().UTC(),
	}
	return out, nil
}
