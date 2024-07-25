package command

import (
	"fmt"
	"github.com/creack/pty"
	"golang.org/x/term"
	"io"
	"os"
)

type CmdWrapperPty struct {
	*CmdWrapper
	ptmx *os.File
}

func (c *CmdWrapperPty) Start() (err error) {
	if c.ptmx != nil {
		return nil
	}
	c.ptmx, err = pty.Start(c.cmd)
	// No need to capture/restore old state, because we close the PTY when we're done.
	_, err = term.MakeRaw(int(c.ptmx.Fd()))
	if err != nil {
		return fmt.Errorf("error putting PTY into raw mode: %w", err)
	}
	return err
}

func (c *CmdWrapperPty) Wait() error {
	defer c.ptmx.Close()
	return c.cmd.Wait()
}

func (c *CmdWrapperPty) StdoutPipe() (io.ReadCloser, error) {
	if c.ptmx == nil {
		c.Start()
	}
	return c.ptmx, nil
}

func (c *CmdWrapperPty) StderrPipe() (io.ReadCloser, error) {
	return nil, fmt.Errorf("not supported in PTY")
}
func (c *CmdWrapperPty) StdinPipe() (io.WriteCloser, error) {
	if c.ptmx == nil {
		c.Start()
	}
	return c.ptmx, nil
}

func (c *CmdWrapperPty) SetCmdArgs() {
	//c.cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}
