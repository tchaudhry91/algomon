package algochecks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/tchaudhry91/algoprom/measure"
)

// Basic Python Algorithm Apply
type PythonAlgorithmer struct {
	VEnv      string `json:"venv"`
	Directory string `json:"directory"`
	logger    *log.Logger
}

func (pa *PythonAlgorithmer) ApplyAlgorithm(ctx context.Context, algorithm string, algorithmParams map[string]string, inputs map[string]measure.Result, workingDir string) (Output, error) {
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
	inputsData, err := json.Marshal(inputs)
	if err != nil {
		return out, fmt.Errorf("Error Marshalling Inputs to JSON: %v", err)
	}
	paramsData, err := json.Marshal(algorithmParams)
	if err != nil {
		return out, fmt.Errorf("Error Marshalling Params to JSON: %v", err)
	}
	err = os.WriteFile(path.Join(workingDir, "inputs.json"), inputsData, 0644)
	if err != nil {
		return out, fmt.Errorf("Error writing inputs file: %v", err)
	}

	err = os.WriteFile(path.Join(workingDir, "params.json"), paramsData, 0644)
	if err != nil {
		return out, fmt.Errorf("Error writing params file: %v", err)
	}

	pythonCmd = fmt.Sprintf("%s python %s --inputs %s --params %s", pythonCmd, path.Join(pa.Directory, algorithm+".py"), "inputs.json", "params.json")
	cmdWrap = append(cmdWrap, pythonCmd)
	cmd := exec.CommandContext(ctx, "sh", cmdWrap...)
	cmd.Dir = workingDir

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
