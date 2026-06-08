package uconfig

import (
	"fmt"
	"io"
	"os"
	"path"
	"reflect"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/omeid/uconfig/flat"
	"github.com/omeid/uconfig/plugins"
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
	_, _ = fmt.Fprintf(w, "Usage:\n\t%s [flags] [command]\n\n", path.Base(os.Args[0]))
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
		if f.Folded() {
			continue
		}

		values := make([]string, len(headers))
		name, _ := f.Name("")
		values[0] = name
		for i, header := range headers[1:] {
			value := f.Meta()[header]
			values[i+1] = value
		}

		_, _ = fmt.Fprintln(w, strings.Join(values, "\t"))

	}

	_ = w.Flush()

	var foldedFields flat.Fields
	for _, f := range c.fields {
		if f.Folded() {
			foldedFields = append(foldedFields, f)
		}
	}

	var usageHeaders []string
	groups := make(map[string][]string)

	for _, p := range c.plugins {
		if u, ok := p.(plugins.Usage); ok {
			header, text := u.Usage(".")
			if header != "" && text != "" {
				if len(groups[header]) == 0 {
					usageHeaders = append(usageHeaders, header)
				}
				groups[header] = append(groups[header], text)
			}
		}
	}

	if len(groups["Files"]) > 0 {
		printFileGroup(w, "Files", groups["Files"])
	}

	if len(foldedFields) > 0 {
		_, _ = fmt.Fprintf(w, "\n\nFOLDED\tTYPE\tUSAGE\n")
		_, _ = fmt.Fprintf(w, "------\t----\t-----\n")

		for i, cf := range foldedFields {
			if i > 0 {
				_, _ = fmt.Fprintf(w, "\t\t\n")
			}
			name, _ := cf.Name("")

			// We use reflect to inspect the element type
			t := reflect.TypeOf(cf.Interface())
			elem := t
			if t.Kind() == reflect.Slice || t.Kind() == reflect.Map {
				elem = t.Elem()
			}

			typeName := t.String()
			pkgPath := elem.PkgPath()
			if pkgPath != "" {
				typeName = strings.ReplaceAll(typeName, pkgPath+".", "")
			}

			_, _ = fmt.Fprintf(w, "%s\t%s\t\n", name, typeName)

			if elem.Kind() == reflect.Struct {
				ptr := reflect.New(elem).Interface()
				subFields, err := flat.View(ptr)
				if err == nil {
					setUsageMeta(subFields)
					for _, sf := range subFields {
						sn, _ := sf.Name("")
						st := reflect.TypeOf(sf.Interface())

						usageMeta := ""
						if u, ok := sf.Meta()["usage"]; ok {
							usageMeta = u
						}

						_, _ = fmt.Fprintf(w, "    .%s\t%s\t%s\n", sn, st.String(), usageMeta)
					}
				}
			}

			// Print any fileset snippets associated with this folded field.
			var filesetTexts []string
			for _, p := range c.plugins {
				if u, ok := p.(plugins.Usage); ok {
					header, text := u.Usage(name)
					if header == "Fileset" && text != "" {
						filesetTexts = append(filesetTexts, text)
					}
				}
			}
			if len(filesetTexts) > 0 {
				_, _ = fmt.Fprintf(w, "\t\t\n")
				printFilesetGroup(w, "Fileset", filesetTexts)
			}
		}
		_ = w.Flush()
	}

	var printedOtherSnippets bool
	for _, header := range usageHeaders {
		if header == "Files" || strings.HasSuffix(header, " Files") {
			continue
		}
		if !printedOtherSnippets {
			_, _ = fmt.Fprintf(w, "\n")
			printedOtherSnippets = true
		} else {
			_, _ = fmt.Fprintf(w, "\n")
		}
		_, _ = fmt.Fprintf(w, "%s:\n", header)
		for _, text := range groups[header] {
			_, _ = fmt.Fprint(w, text)
		}
	}

	err := w.Flush()
	if err != nil {
		// we are asked for usage which means it is interactive use
		// and so panicking is acceptable.
		panic(err)
	}
}

// printFileGroup prints a file group with the header inline:
//
//	Header: path1
//	        path2
//	        path3
func printFileGroup(w io.Writer, header string, texts []string) {
	prefix := header + ": "
	padding := strings.Repeat(" ", len(prefix))

	// Collect all lines from all texts.
	var lines []string
	for _, text := range texts {
		for _, line := range strings.Split(strings.TrimRight(text, "\n"), "\n") {
			line = strings.TrimSpace(line)
			if line != "" {
				lines = append(lines, line)
			}
		}
	}

	for i, line := range lines {
		if i == 0 {
			_, _ = fmt.Fprintf(w, "\n%s%s\n", prefix, line)
		} else {
			_, _ = fmt.Fprintf(w, "%s%s\n", padding, line)
		}
	}
}

// printFilesetGroup prints a fileset group inline:
//
//	    Fileset:            absolute: /etc/app/*.yaml
//	                        relative: apps.d/*.json
func printFilesetGroup(w io.Writer, header string, texts []string) {
	prefix := "    " + header + ":\t"
	padding := "\t"

	// Collect all lines from all texts.
	var lines []string
	for _, text := range texts {
		for _, line := range strings.Split(strings.TrimRight(text, "\n"), "\n") {
			if line != "" {
				lines = append(lines, line)
			}
		}
	}

	for i, line := range lines {
		if i == 0 {
			_, _ = fmt.Fprintf(w, "%s%s\n", prefix, line)
		} else {
			_, _ = fmt.Fprintf(w, "%s%s\n", padding, line)
		}
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

	tags = append(tags, "fields")

	for key := range tagMap {
		tags = append(tags, key)
	}

	weights := map[string]int{
		"fields": 1,
		"usage":  99,
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
