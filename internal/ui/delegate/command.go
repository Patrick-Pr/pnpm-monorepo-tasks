package delegate

import (
	"io"
	"os/exec"
)

// scriptCommand waits for the script's output streams to close before returning
// the terminal to Bubble Tea. Package managers can exit before their children
// finish shutting down and restoring the terminal's input mode.
type scriptCommand struct {
	*exec.Cmd
}

func (c *scriptCommand) SetStdin(r io.Reader) {
	if c.Stdin == nil {
		c.Stdin = r
	}
}

func (c *scriptCommand) SetStdout(w io.Writer) {
	if c.Stdout == nil {
		c.Stdout = scriptOutput{w}
	}
}

func (c *scriptCommand) SetStderr(w io.Writer) {
	if c.Stderr == nil {
		c.Stderr = scriptOutput{w}
	}
}

// Hide the underlying *os.File so os/exec uses a pipe and waits for its copy
// goroutine, including output inherited by descendant processes.
type scriptOutput struct {
	io.Writer
}
