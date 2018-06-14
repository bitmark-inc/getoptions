// Copyright (c) 2014-2017 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package getoptions_test

import (
	"reflect"
	"testing"

	"github.com/bitmark-inc/getoptions"
)

type testItem struct {
	in []string
	op getoptions.OptionsMap
	ar []string
}

func TestGetOptions(t *testing.T) {

	flags := []getoptions.Option{
		{Long: "verbose", HasArg: getoptions.NO_ARGUMENT, Short: 'v'},
		{Long: "hello", HasArg: getoptions.REQUIRED_ARGUMENT, Short: 'H'},
		{Long: "say", HasArg: getoptions.REQUIRED_ARGUMENT, Short: 's'},
		{Long: "xyz", HasArg: getoptions.NO_ARGUMENT, Short: 'x'},
		{Long: "test", HasArg: getoptions.OPTIONAL_ARGUMENT, Short: 't'},
	}

	tests := []testItem{
		{
			in: []string{"-v", "-x", "-v", "--hello=yes", "--test", "data", "argon", "999", "--verbose"},
			op: getoptions.OptionsMap{"test": []string{"data"}, "verbose": []string{"", "", ""}, "xyz": []string{""}, "hello": []string{"yes"}},
			ar: []string{"argon", "999"},
		},
		{
			in: []string{"-vv", "--hello", "yes", "-xtdata", "argon", "999", "--verbose"},
			op: getoptions.OptionsMap{"test": []string{"data"}, "verbose": []string{"", "", ""}, "xyz": []string{""}, "hello": []string{"yes"}},
			ar: []string{"argon", "999"},
		},
		{
			in: []string{"-v", "-x", "--hello=yes", "--", "multi-word", "999", "-verbose"},
			op: getoptions.OptionsMap{"verbose": []string{""}, "xyz": []string{""}, "hello": []string{"yes"}},
			ar: []string{"multi-word", "999", "-verbose"},
		},
		{
			in: []string{"--say=hello", "--say=there", "--say=world", "--", "hello", "earth"},
			op: getoptions.OptionsMap{"say": []string{"hello", "there", "world"}},
			ar: []string{"hello", "earth"},
		},
		{
			in: []string{"--say", "hello", "--say", "there", "--say", "world", "--", "hello", "earth"},
			op: getoptions.OptionsMap{"say": []string{"hello", "there", "world"}},
			ar: []string{"hello", "earth"},
		},
		{
			in: []string{"-s", "hello", "-sthere", "--say", "world", "--", "hello", "earth"},
			op: getoptions.OptionsMap{"say": []string{"hello", "there", "world"}},
			ar: []string{"hello", "earth"},
		},
		{
			in: []string{"-t", "--test", "-t", "one", "--test", "two", "alpha", "--test=three", "beta", "-tfour", "gamma"},
			op: getoptions.OptionsMap{"test": []string{"", "", "one", "two", "three", "four"}},
			ar: []string{"alpha", "beta", "gamma"},
		},
	}

	for i, s := range tests {
		options, arguments, err := getoptions.Getopt(s.in, flags)
		if nil != err {
			t.Errorf("%d: error: %v", i, err)
		}
		if !reflect.DeepEqual(options, s.op) {
			t.Errorf("%d: options: %#v  expected: %#v", i, options, s.op)
		}
		if !reflect.DeepEqual(arguments, s.ar) {
			t.Errorf("%d: arguments: %#v expected: %#v", i, arguments, s.ar)
		}
	}

}
