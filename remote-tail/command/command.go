package command

import (
	"os/exec"
	"bufio"
	"io"
)

type Command struct {
	command *exec.Cmd
	Host    string
	Script  string
	Stdout  io.ReadCloser
	Stderr  io.ReadCloser
}

type Message struct {
	Host    string
	Content string
}

func NewCommand(host string, script string) (cmd *Command, err error) {
	cmd = &Command{
		Host: host,
		Script: script,
	}

	cmd.command = exec.Command("ssh", "-f", host, script)
	stdout, err := cmd.command.StdoutPipe()
	if err != nil {
		return nil, err
	}

	cmd.Stdout = stdout

	stderr, err := cmd.command.StderrPipe()
	if err != nil {
		return nil, err
	}

	cmd.Stderr = stderr

	return
}

func (cmd *Command) Execute(stdout chan Message, stderr chan Message) {
	cmd.command.Start()

	bindOutput(cmd.Host, &cmd.Stdout, stdout)

	cmd.command.Wait()
}

func bindOutput(host string, input *io.ReadCloser, output chan Message) {
	reader := bufio.NewReader(*input)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}

		output <- Message{
			Host: host,
			Content: line,
		}
	}
}