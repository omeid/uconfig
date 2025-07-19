package uconfig

import (
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/omeid/uconfig/flat"
	"github.com/omeid/uconfig/plugins"
	"github.com/omeid/uconfig/plugins/file"
	"github.com/omeid/uconfig/plugins/flag"
)

const usageTag = "usage"

func init() {
	plugins.RegisterTag(usageTag)
}

// UsageOutput is the io.Writer used by Usage message printer.
var UsageOutput io.Writer = os.Stdout

// Usage prints out the current config fields, flags, env vars
// and any other source and setting.
func (c *config[C]) Usage() {
	setUsageMeta(c.fields)
	headers := getHeaders(c.fields)

	w := tabwriter.NewWriter(UsageOutput, 0, 0, 4, ' ', 0)
	_, _ = fmt.Fprintf(w, "Usage:\n\t%s [flags] [command]\n", path.Base(os.Args[0]))
	_, _ = fmt.Fprintf(w, "\nConfigurations:\n")
	_, _ = fmt.Fprintln(w, strings.ToUpper(strings.Join(headers, "\t")))

	dashes := make([]string, len(headers))
	for i, f := range headers {
		n := max(len(f), 5)
		dashes[i] = strings.Repeat("-", n)
	}
	_, _ = fmt.Fprintln(w, strings.Join(dashes, "\t"))

	sort.SliceStable(c.fields, func(i, j int) bool {
		return flag.IsCommand(c.fields[j]) // move command to last.
	})

	for _, f := range c.fields {

		values := make([]string, len(headers))
		name, _ := f.Name("")
		values[0] = name
		for i, header := range headers[1:] {
			value := f.Meta()[header]
			values[i+1] = value
		}

		_, _ = fmt.Fprintln(w, strings.Join(values, "\t"))

	}

	files := []string{}

	for _, p := range c.plugins {
		if p, ok := p.(file.Plugin); ok {
			files = append(files, p.FilePath())
		}
	}

	if len(files) > 0 {
		_, _ = fmt.Fprintf(w, "\nConfiguration Files:\n")
		for _, fp := range files {
			_, _ = fmt.Fprintf(w, "\t%s\n", fp)
		}

	}

	err := w.Flush()
	if err != nil {
		// we are asked for usage which means it is interactive use
		// and so panicing is acceptable.
		panic(err)
	}
}

func setUsageMeta(fs flat.Fields) {
	for _, f := range fs {
		usage, ok := f.Tag(usageTag)
		if !ok {
			continue
		}

		f.Meta()[usageTag] = usage

	}
}

func getHeaders(fs flat.Fields) []string {
	tagMap := map[string]struct{}{}

	for _, f := range fs {
		for key := range f.Meta() {
			tagMap[key] = struct{}{}
		}
	}

	tags := make([]string, 0, len(tagMap)+2)

	tags = append(tags, "field")

	for key := range tagMap {
		tags = append(tags, key)
	}

	weights := map[string]int{
		"field": 1,
		"usage": 99,
		"flag":  3,
		"env":   4,
	}

	weight := func(tags []string, i int) int {
		key := tags[i]
		w, ok := weights[key]
		if !ok {
			return 98
		}
		return w
	}

	sort.SliceStable(tags, func(i, j int) bool {
		iw := weight(tags, i)
		jw := weight(tags, j)

		if iw == jw {
			return tags[i] < tags[j]
		}

		return iw < jw
	})

	return tags
}
