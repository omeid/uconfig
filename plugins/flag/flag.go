// Package flag provides flags support for uconfig
package flag

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/omeid/uconfig/flat"
	"github.com/omeid/uconfig/plugins"
)

const tag = "flag"

func init() {
	plugins.RegisterTag(tag)
}

// ErrorHandling defines how FlagSet.Parse behaves if the parse fails.
type ErrorHandling flag.ErrorHandling

// These constants cause FlagSet.Parse to behave as described if the parse fails.
const (
	ContinueOnError = ErrorHandling(flag.ContinueOnError)
	ExitOnError     = ErrorHandling(flag.ExitOnError)
	PanicOnError    = ErrorHandling(flag.PanicOnError)
)

// New returns a new Flags
func New(name string, errorHandling ErrorHandling, args []string) plugins.Plugin {

	fs := flag.NewFlagSet(name, flag.ErrorHandling(errorHandling))
	fs.Usage = func() {}

	return &visitor{
		fs:   fs,
		args: args,
	}
}

// Standard returns a set of flags configured in the common way.
// It is same as: `New(os.Args[0], ContinueOnError, os.Args[1:])`
func Standard() plugins.Plugin {
	return New(os.Args[0], ContinueOnError, os.Args[1:])
}

var _ plugins.Visitor = (*visitor)(nil)

type visitor struct {
	fs      *flag.FlagSet
	args    []string
	command flat.Field
}

func makeFlagName(name string) string {
	name = strings.Replace(name, ".", "-", -1)
	name = strings.ToLower(name)
	return name
}

type fieldFlag struct {
	flat.Field
}

func (ff *fieldFlag) String() string {
	if ff == nil {
		return ""
	}

	return fmt.Sprintf("%s", ff.Field.Interface())
}

// Used by standard library flag package.
func (f *fieldFlag) IsBoolFlag() bool {
	return reflect.ValueOf(f.Field.Interface()).Kind() == reflect.Bool
}

func (v *visitor) Visit(fields flat.Fields) error {

	for _, f := range fields {

		name, explicit := f.Name(tag)

		if name == "-" {
			continue
		}

		if !explicit {
			name = makeFlagName(name)
		}

		opts, _ := f.Tag(tag)
		_, opts, _ = strings.Cut(opts, ",")
		if strings.Contains(opts, "command") {
			v.command = f
			f.Meta()[tag] = "[command]"
		} else {
			usage, _ := f.Tag("usage")
			f.Meta()[tag] = "-" + name
			v.fs.Var(&fieldFlag{f}, name, usage)
		}
	}

	return nil
}

func extraCommand(args []string) (string, []string) {
	if len(args) == 0 {
		return "", args
	}

	command := args[0]

	if command != "" && command[0] == '-' {
		return "", args
	}

	return args[0], args[1:]
}

func (v *visitor) Parse() error {

	args := v.args

	if v.command != nil {
		var command string
		command, args = extraCommand(args)

		err := v.command.Set(command)
		if err != nil {
			return err
		}
	}

	err := v.fs.Parse(args)

	if errors.Is(err, flag.ErrHelp) {
		return plugins.ErrUsage
	}

	return err
}
