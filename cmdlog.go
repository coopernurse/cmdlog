package cmdlog

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"time"
)

func Run(name string, cmd *exec.Cmd) *Result {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now().UnixNano()
	err := cmd.Run()
	end := time.Now().UnixNano()

	return newResult(name, start, end, cmd, err, stdout.Bytes(), stderr.Bytes())
}

type Result struct {
	Name      string
	StartNano int64
	EndNano   int64

	Host string
	Dir  string
	Path string
	Args []string

	Error   error
	ExitStr string
	Stdout  []byte
	Stderr  []byte
}

func (r *Result) ElapsedSeconds() float64 {
	return float64(r.EndNano-r.StartNano) / float64(1e9)
}

func (r *Result) StartDate(format string) string {
	t := time.Unix(0, r.StartNano)
	return t.Format(format)
}

func newResult(name string, start int64, end int64, cmd *exec.Cmd, runErr error, stdout []byte, stderr []byte) *Result {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
		log.Println("cmdlog: os.Hostname() failed:", err)
	}
	return &Result{
		Name:      name,
		StartNano: start,
		EndNano:   end,
		Host:      hostname,
		Dir:       cmd.Dir,
		Path:      cmd.Path,
		Args:      cmd.Args,
		Error:     runErr,
		ExitStr:   cmd.ProcessState.String(),
		Stdout:    stdout,
		Stderr:    stderr,
	}
}
