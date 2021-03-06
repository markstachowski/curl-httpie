package args

import (
	"sort"
	"strings"
)

type Opts []string

func (opts Opts) index(opt string) int {
	off := 1
	if len(opt) > 1 {
		off = 2
	}
	for i, o := range opts {
		if len(o) >= 2 && o[0] == '-' {
			if o[off:] == opt {
				return i
			}
		}
	}
	return -1
}

// Has returns true if opt flag is in opts.
func (opts Opts) Has(opt string) bool {
	return opts.index(opt) != -1
}

func (opts Opts) Val(opt string) string {
	if idx := opts.index(opt); idx != -1 && idx+1 < len(opts) {
		return opts[idx+1]
	}
	return ""
}

// Parse converts an HTTPie like argv into a list of curl options.
func Parse(argv []string) (opts Opts) {
	args := []string{}
	sort.Strings(curlLongValues)
	more := true
	for i := 1; i < len(argv); i++ {
		arg := argv[i]
		if !more || len(arg) < 2 || arg[0] != '-' {
			args = append(args, arg)
			continue
		}
		if arg == "--" {
			// Enf of opts marker
			more = false
			continue
		}
		opts = append(opts, arg)
		if arg[1] == '-' {
			if longHasValue(arg[2:]) && i+1 < len(argv) {
				opts = append(opts, argv[i+1])
				i++
			}
			continue
		}
		for j := 1; j < len(arg); j++ {
			if strings.IndexByte(curlShortValues, arg[j]) != -1 {
				if j == len(arg)-1 && i+1 < len(argv) {
					opts = append(opts, argv[i+1])
				}
				break
			}
		}
	}
	if len(args) > 0 {
		opts = append(opts, parseFancyArgs(args)...)
	}
	return
}

func longHasValue(arg string) bool {
	i := sort.SearchStrings(curlLongValues, arg)
	return i < len(curlLongValues) && curlLongValues[i] == arg
}
