package actions

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/charmbracelet/log"
	"os"
	"os/exec"
	"path"
	"time"
)

// Basic Python Actioner
type PythonActioner struct {
	VEnv        string            `json:"venv"`
	Directory   string            `json:"directory"`
	EnvOverride map[string]string `json:"env_override"`
	logger      *log.Logger
}

func (pa *PythonActioner) Action(ctx context.Context, action string, input string, params map[string]string, workingDir string) (Output, error) {
	out := Output{
		RC:          -1,
		CombinedOut: "",
		Timestamp:   time.Now().UTC(),
	}
	cmdWrap := []string{"-c"}
	pythonCmd := ""
	if pa.VEnv != "" {
		pythonCmd = fmt.Sprintf("source %s/bin/activate;", pa.VEnv)
	}
	// Write Inputs and Params
	paramsData, err := json.Marshal(params)
	if err != nil {
		return out, fmt.Errorf("Error Marshalling Params to JSON: %v", err)
	}
	err = os.WriteFile(path.Join(workingDir, "inputs.json"), []byte(input), 0644)
	if err != nil {
		return out, fmt.Errorf("Error writing inputs file: %v", err)
	}

	err = os.WriteFile(path.Join(workingDir, "params.json"), paramsData, 0644)
	if err != nil {
		return out, fmt.Errorf("Error writing params file: %v", err)
	}

	pythonCmd = fmt.Sprintf("%s python %s --inputs %s --params %s", pythonCmd, path.Join(pa.Directory, action+".py"), "inputs.json", "params.json")
	cmdWrap = append(cmdWrap, pythonCmd)
	cmd := exec.CommandContext(ctx, "sh", cmdWrap...)
	cmd.Dir = workingDir
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, envMapToSlice(pa.EnvOverride)...)

	combined, err := cmd.CombinedOutput()
	if err != nil {
		if exiterror, ok := err.(*exec.ExitError); ok {
			out.RC = exiterror.ProcessState.ExitCode()
		}
		out.Error = err
	}
	out.RC = cmd.ProcessState.ExitCode()
	out.CombinedOut = string(combined)

	return out, nil
}

func envMapToSlice(env map[string]string) []string {
	envs := []string{}
	for k, v := range env {
		envs = append(envs, fmt.Sprintf("%s=%s", k, v))
	}
	return envs
}
