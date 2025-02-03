package algochecks

import (
	"context"
	"fmt"
	"log"
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

func (pa *PythonAlgorithmer) ApplyAlgorithm(ctx context.Context, algorithm string, algorithmParams map[string]string, inputs []measure.Measurement, workingDir string) (Output, error) {
	cmdWrap := []string{"-c"}
	pythonCmd := ""
	if pa.VEnv != "" {
		pythonCmd = fmt.Sprintf("source %s/bin/activate;", pa.VEnv)
	}
	// Write out the inputs and params to files

	pythonCmd = fmt.Sprintf("%s python %s", pythonCmd, path.Join(pa.Directory, algorithm+".py"))
	cmdWrap = append(cmdWrap, pythonCmd)
	cmd := exec.CommandContext(ctx, "sh", cmdWrap...)
	cmd.Dir = workingDir

	out := Output{
		RC:          -1,
		CombinedOut: "",
		Timestamp:   time.Now().UTC(),
	}

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
