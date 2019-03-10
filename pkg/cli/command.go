package cli

import (
	"flag"
	"fmt"
	"os"
)

type Handler interface {
	Handle(args []string)
}

type HandlerFunc func(args []string)

func (handle HandlerFunc) Handle(args []string) {
	handle(args)
}

type Command struct {
	name     string
	handler  Handler
	flag     flag.FlagSet
	usage    string
	parent   *Command
	commands map[string]Command
}

func New(name string, usage string) *Command {
	return &Command{name: name, usage: usage}
}

func (cmd *Command) Usage(usage string) {
	cmd.usage = usage
}

func (cmd *Command) Flag() *flag.FlagSet {
	return &cmd.flag
}

func (cmd *Command) Handle(handler Handler) {
	cmd.handler = handler
}

func (cmd *Command) HandleFunc(handler HandlerFunc) {
	cmd.handler = handler
}

func (cmd *Command) Command(name string, help string, init func(*Command)) {
	if cmd.commands == nil {
		cmd.commands = make(map[string]Command)
	}

	command := Command{name: name, parent: cmd}
	init(&command)
	cmd.commands[name] = command
}

func (cmd *Command) Run(args []string) {
	cmd.flag.Usage = func() { cmd.printUsage() }
	cmd.flag.Parse(args)

	args = cmd.flag.Args()

	if cmd.handler != nil {
		cmd.handler.Handle(args)
	} else if len(args) > 0 {
		if cmd, ok := cmd.commands[args[0]]; ok {
			cmd.Run(args[1:])
		} else {
			fmt.Fprintf(flag.CommandLine.Output(), "%s: '%s' is not a command\n", cmd.invocation(), args[0])
			os.Exit(2)
		}
	} else {
		cmd.printUsage()
	}
}

func (cmd *Command) Errorf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s: %s", cmd.invocation(), fmt.Sprintf(format, args...))
	exitCode = 2
}

func (cmd *Command) Fatalf(format string, args ...interface{}) {
	cmd.Errorf(format, args...)
	exitCode = 1
	exit()
}

func (cmd *Command) invocation() string {
	if cmd.parent == nil {
		return cmd.name
	} else {
		return cmd.parent.invocation() + " " + cmd.name
	}
}

func (cmd *Command) printUsage() {
	fmt.Fprintf(cmd.flag.Output(), "usage: %s %s\n", cmd.invocation(), cmd.usage)
	cmd.flag.PrintDefaults()
	os.Exit(2)
}
