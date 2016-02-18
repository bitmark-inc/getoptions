// Copyright (c) 2014-2015 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package getoptions

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type optionType int

// argument requirements
const (
	NO_ARGUMENT = optionType(iota)
	REQUIRED_ARGUMENT = optionType(iota)
	OPTIONAL_ARGUMENT = optionType(iota)
)

// option structure to represent a single option definition
type Option struct {
	Long   string     // long option e.g.: "verbose"
	HasArg optionType // one of: NO_ARGUMENT, REQUIRED_ARGUMENT or OPTIONAL_ARGUMENT
	Short  rune       // short option e.g.: 'v'
}

// returned options, the Long name is used in this map
// repeated option values are returned as a string slice
type OptionsMap map[string][]string


// error strings
const (
	errorUnknown = "option: %q is unknown"
	errorMissing = "option: %q is missing its argument"
	errorNoArg   = "option: %q does not take value: %q"
)

// parse options from OS command-line
//
// Return values:
//   program_name     - string
//   options          - map["option"]=[]string{"value1","value2"}
//                      use len(option["verbose"]) to detect a sequence like: -v -v -v
//                      the actual value will be options["verbose"] = []string{"", "", ""}
//   arguments        - []string  (all items not starting with "-" that do not belong to option and everything after --)
//   err              - nil if parsing was sucessful
func GetOS(flags []Option) (program string, options OptionsMap, arguments []string, err error) {
	options, arguments, err = Getopt(os.Args[1:], flags)
	program = filepath.Base(os.Args[0])
	return
}

// parse options from an arbitrary array of strings
//
// Note that the input string slices does not contain any program name
func Getopt(inputs []string, flags []Option) (options OptionsMap, arguments []string, err error) {

	options = make(OptionsMap)
	arguments = make([]string, 0, 10)
	hasArg := make(map[string]optionType)
	alias := make(map[rune]string)

	for _, f := range flags {
		hasArg[f.Long] = f.HasArg
		alias[f.Short] = f.Long
	}

	n := 0

	// for the cases:  -t arg  --test arg
	wantArgument := NO_ARGUMENT
	wantedBy := ""

	// parse each option
loop:
	for i, item := range inputs {

		// previous option needs an argument
		if NO_ARGUMENT != wantArgument {
			name := wantedBy
			value := item
			doContinue := true
			if len(item) > 1 && '-' == item[0] {
				if REQUIRED_ARGUMENT == wantArgument {
					err = fmt.Errorf(errorMissing, name)
					return
				}
				value = ""
				doContinue = false
			}
			options[name] = append(options[name], value)
			wantArgument = NO_ARGUMENT
			if doContinue {
				continue loop
			}
		}

		// empty string and any single character is an argument
		// even a '-' (often used to represint stdin/stdout
		// any value that does not start with a '-'in an argument
		if len(item) <= 1 || '-' != item[0] {
			arguments = append(arguments, item)
			continue loop
		}

		// check for end of options
		if "--" == item {
			n = i + 1
			break loop
		}

		// check for long option
		// --file name  --file=name
		if "--" == item[:2] {
			name := item[2:]
			value := ""
			s := strings.SplitN(name, "=", 2)

			if 2 == len(s) {
				name = s[0]
				value = s[1]
				h, ok := hasArg[name]
				if !ok {
					err = fmt.Errorf(errorUnknown, name)
					return
				}
				if NO_ARGUMENT == h {
					err = fmt.Errorf(errorNoArg, name, value)
					return
				}
				options[name] = append(options[name], value)
				continue loop
			}
			h, ok := hasArg[name]
			if !ok {
				err = fmt.Errorf(errorUnknown, name)
				return
			}
			if NO_ARGUMENT == h {
				options[name] = append(options[name], "")
			} else {
				wantArgument = h
				wantedBy = name
			}
			continue loop
		}

		// also checks for merged short options like:
		//   -abcfdata
		// if -f has optional/required setting then this is equivalent to:
		//   -a -b -c -f data
		// note that the option with an argument uses the rest of the string
		if '-' == item[0] {
			shortOptions := item[1:]
			end := len(shortOptions) - 1
		short:
			for j, c := range shortOptions {
				name, ok := alias[c]
				if !ok {
					err = fmt.Errorf(errorUnknown, name)
					return
				}
				h, ok := hasArg[name]
				if !ok {
					err = fmt.Errorf(errorUnknown, name)
					return
				}
				// if argument possible the consume the remaining string
				// or if lat character the consume next item
				if NO_ARGUMENT != h {
					if j == end {
						wantArgument = h
						wantedBy = name
						continue loop
					}
					options[name] = append(options[name], shortOptions[j+1:])
					break short
				}
				options[name] = append(options[name], "")
			}
		}
	}

	// check if missing argument
	if NO_ARGUMENT != wantArgument {
		err := fmt.Errorf(errorMissing, wantedBy)
		return nil, nil, err
	}

	// remaining items are noramal arguments
	if 0 != n {
		arguments = append(arguments, inputs[n:]...)
	}

	return
}
