package uconfig

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/omeid/uconfig/plugins"
)

type Command interface {
	Name() string
	Run() error
	Error() error
}

type command[T any] struct {
	name    string
	program func(conf T) error
	conf    *T
	config  Config

	// this allows NewCommand to bubble up the error.
	err error
}

func (c *command[T]) Name() string {
	return c.name
}

func (c *command[T]) Error() error {
	return c.err
}

func (c *command[T]) Run() error {
	err := c.config.Parse()
	if err != nil {
		return err
	}

	return c.program(*c.conf)
}

type configer interface {
	Config() Config
}

func (c *command[T]) Config() Config {
	return c.config
}

func NewCommand[T any](classic bool, name string, program func(conf T) error, plugins ...plugins.Plugin) Command {
	conf := new(T)

	var config Config
	var err error

	config, err = newConfig(conf, plugins)

	return &command[T]{
		name:    name,
		program: program,

		conf:   conf,
		config: config,

		err: err,
	}
}

func ClassicCommand[T any](name string, program func(conf T) error, files Files, userPlugins ...plugins.Plugin) Command {
	conf := new(T)

	var config Config
	var err error

	config, err = classicConfig(conf, files, userPlugins, true)

	return &command[T]{
		name:    name,
		program: program,

		conf:   conf,
		config: config,

		err: err,
	}
}

type program struct {
	commands map[string]Command
}

func (p *program) usage() {

	fmt.Printf("Supported Commands: ")

	names := make([]string, 0, len(p.commands))

	for name := range p.commands {
		names = append(names, name)
	}

	slices.Sort(names)

	for _, name := range names {
		if name == "" {
			fmt.Printf("[default] ")
		} else {
			fmt.Printf("%s ", name)
		}
	}

	fmt.Printf("\n")

	for _, name := range names {
		p.commands[name].(configer).Config().(*config).usage(&name)
	}

}

func (p *program) Run() error {

	arg1 := ""

	if len(os.Args) > 1 {
		arg1 = os.Args[1]
	}

	if arg1 != "" {
		if arg1 == "--help" || arg1 == "-help" {
			p.usage()
			return nil
		}

		// allow flags for default command.
		if arg1[0] == '-' {
			arg1 = ""
		}
	}

	command, ok := p.commands[arg1]

	if !ok {
		fmt.Printf("command provided but not defined: %s\n\n", arg1)
		os.Exit(1)
	}

	return command.Run()
}

func Commands(commands ...Command) error {

	commandsMap := make(map[string]Command, len(commands))

	for _, command := range commands {

		// if there is any error already
		// return it.
		if command.Error() != nil {
			return command.Error()
		}

		commandsMap[command.Name()] = command
	}
	p := program{commands: commandsMap}

	err := p.Run()

	if errors.Is(err, ErrUsage) {
		p.usage()
		os.Exit(0)
	}

	return err
}
