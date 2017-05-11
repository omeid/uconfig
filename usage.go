package uconfig

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/omeid/uconfig/flat"
)

func (c *config) Usage() {

	headers := getHeaders(c.fields)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.Debug)
	fmt.Fprintf(w, "\nSupported Fields:\n")
	fmt.Fprintln(w, strings.ToUpper(strings.Join(headers, "\t")))

	for _, f := range c.fields {

		values := make([]string, len(headers))
		values[0] = f.Name()
		for i, header := range headers[1:] {
			value, _ := f.Meta()[header]
			values[i+1] = value
		}

		fmt.Fprintln(w, strings.Join(values, "\t"))

	}

	w.Flush()
}

type null struct{}

func getHeaders(fs flat.Fields) []string {
	tagMap := map[string]null{}

	for _, f := range fs {
		for key := range f.Meta() {
			tagMap[key] = struct{}{}
		}
	}

	tags := make([]string, 0, len(tagMap)+1)

	tags = append(tags, "field")

	for key := range tagMap {
		tags = append(tags, key)
	}

	weights := map[string]int{
		"field": 1,
		"flag":  2,
		"usage": 4,
		"env":   3,
	}

	weight := func(tags []string, i int) int {
		key := tags[i]
		w, ok := weights[key]
		if !ok {
			return 99
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
