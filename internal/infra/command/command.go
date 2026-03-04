package command

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"sync"
)

type Command struct {
	Command string
	Args    []string
	Dir     string
}

func NewDirCommand(cmd string, args []string, dir string) Command {
	return Command{
		Command: cmd,
		Dir:     dir,
		Args:    args,
	}
}

func NewCommand(cmd string, args []string) Command {
	return Command{
		Command: cmd,
		Dir:     "",
		Args:    args,
	}
}

func (c *Command) Exec() ([]byte, error) {
	cmd := exec.Command(c.Command, c.Args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return out, err
	}
	return out, nil
}

func (c *Command) StreamExec(ctx context.Context, ch chan<- []byte) error {
	cmd := exec.CommandContext(ctx, c.Command, c.Args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(2)

	processPipe := func(r io.Reader) {
		defer wg.Done()
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			line := scanner.Bytes()

			b := make([]byte, len(line))
			copy(b, line)

			select {
			case <-ctx.Done():
				return
			case ch <- b:
			}
		}
	}

	go processPipe(stdout)
	go processPipe(stderr)

	go func() {
		_ = cmd.Wait()
		wg.Wait()
		close(ch)
	}()

	return nil
}
