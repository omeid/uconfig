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
const commandFieldName = "[command]"

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
		fs:          fs,
		args:        args,
		requiredSet: map[string]bool{},
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

	requiredSet map[string]bool
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

		required := strings.Contains(opts, "required")

		if strings.Contains(opts, "command") {
			v.command = f
			name := commandFieldName
			f.Meta()[tag] = name
			if required {
				v.requiredSet[name] = false
			}
		} else {
			usage, _ := f.Tag("usage")
			v.fs.Var(&fieldFlag{f}, name, usage)
			if required {
				v.requiredSet[name] = false
			}

			name := "-" + name
			f.Meta()[tag] = name
		}
	}

	return nil
}

func extraCommand(args []string) (string, []string, bool) {
	lastIndex := len(args) - 1

	if lastIndex < 0 {
		return "", args, false
	}

	command := args[lastIndex]

	if command != "" && command[0] == '-' {
		return "", args, false
	}

	return command, args[:lastIndex], command != ""
}

func (v *visitor) Parse() error {

	args := v.args

	if v.command != nil {
		var command string

		var set bool
		command, args, set = extraCommand(args)

		if set {
			err := v.command.Set(command)
			if err != nil {
				return err
			}
			// we have visissted the command.
			v.requiredSet[commandFieldName] = true
		}

	}

	err := v.fs.Parse(args)

	if errors.Is(err, flag.ErrHelp) {
		return plugins.ErrUsage
	}

	if err != nil {
		return err
	}

	v.fs.Visit(func(f *flag.Flag) {
		v.requiredSet[f.Name] = true
	})

	for field, set := range v.requiredSet {
		if !set {
			err = errors.Join(err, errors.New("Missing required flag: "+field))
		}
	}

	if err != nil {
		return err
	}

	return nil
}
