package daemon

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/makii42/gottaw/output"
)

// Daemon is the interface to a background process that is
// started after a successful pipeline execution
type Daemon interface {
	// Start is the trigger to kick of the daemon. Returns an
	// error if the Start fails or if it is running already
	Start() error
	// Stop is the trigger to stop the daemon. Returns an
	// error if the Stop fails or if nothing is running already
	Stop() error
}

type daemon struct {
	l       *output.Logger
	cmdName string
	args    []string
	cmd     *exec.Cmd
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewDaemon creates a new Daemon that can be started and
// stopped using the interface methods.
func NewDaemon(cmdStr string) Daemon {
	elements := strings.Split(cmdStr, " ")
	command, elements := elements[0], elements[1:]
	return &daemon{
		cmdName: command,
		args:    elements,
		cmd:     nil,
	}
}

func (d *daemon) Start() error {
	if d.cmd != nil || d.ctx != nil || d.cancel != nil {
		return fmt.Errorf("command already present")
	}
	d.ctx, d.cancel = context.WithCancel(context.Background())
	d.cmd = exec.CommandContext(d.ctx, d.cmdName, d.args...)
	d.cmd.Stdout = os.Stdout
	d.cmd.Stderr = os.Stderr

	if err := d.cmd.Start(); err != nil {
		return err
	}
	return nil
}

func (d *daemon) Stop() error {
	if d.cmd == nil || d.ctx == nil || d.cancel == nil {
		return nil
	}
	d.cancel()
	d.cmd, d.ctx, d.cancel = nil, nil, nil
	return nil
}
