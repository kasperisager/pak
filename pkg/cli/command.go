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
	name        string
	handler     Handler
	flag        flag.FlagSet
	usage       string
	parent      *Command
	subcommands map[string]Command
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

func (cmd *Command) AddCommand(name string, help string, init func(*Command)) {
	if cmd.subcommands == nil {
		cmd.subcommands = make(map[string]Command)
	}

	command := Command{name: name, parent: cmd}
	init(&command)
	cmd.subcommands[name] = command
}

func (cmd *Command) Run(args []string) {
	cmd.flag.Usage = func() { cmd.printUsage() }
	cmd.flag.Parse(args)

	args = cmd.flag.Args()

	switch {
	case cmd.handler != nil:
		cmd.handler.Handle(args)

	case len(args) > 0:
		if subcmd, ok := cmd.subcommands[args[0]]; ok {
			subcmd.Run(args[1:])
		} else {
			cmd.Fatalf("'%s' is not a command", args[0])
		}

	default:
		cmd.printUsage()
	}
}

func (cmd *Command) Error(message string) {
	fmt.Fprintf(
		os.Stderr,
		"%s: %s\n",
		cmd.invocation(),
		message,
	)
	ExitCode(1)
}

func (cmd *Command) Errorf(format string, args ...interface{}) {
	cmd.Error(fmt.Sprintf(format, args...))
}

func (cmd *Command) Fatal(message string) {
	cmd.Error(message)
	Exit()
}

func (cmd *Command) Fatalf(format string, args ...interface{}) {
	cmd.Fatal(fmt.Sprintf(format, args...))
}

func (cmd *Command) invocation() string {
	if cmd.parent == nil {
		return cmd.name
	} else {
		return cmd.parent.invocation() + " " + cmd.name
	}
}

func (cmd *Command) printUsage() {
	fmt.Fprintf(
		os.Stderr,
		"usage: %s %s\n",
		cmd.invocation(),
		cmd.usage,
	)
	cmd.flag.PrintDefaults()
	ExitCode(2)
	Exit()
}
