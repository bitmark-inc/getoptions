// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package getoptions_test

import (
	"fmt"
	"github.com/bitmark-inc/getoptions"
)

func ExampleGetopt() {

	// define options
	flags := []getoptions.Option{
		{Long: "help", HasArg: getoptions.NO_ARGUMENT, Short: 'h'},
		{Long: "output", HasArg: getoptions.REQUIRED_ARGUMENT, Short: 'o'},
		{Long: "verbose", HasArg: getoptions.NO_ARGUMENT, Short: 'v'},
	}

	// simulated command-line arguments
	args := []string{"--help", "--output=data1", "zero", "-odata2", "-vvv", "one", "two"}

	// parse options
	options, arguments, err := getoptions.Getopt(args, flags)

	// display results
	if nil != err {
		fmt.Printf("parse error: %v\n", err)
	} else {
		for _, op := range []string{"help", "output", "verbose"} {
			fmt.Printf("option[%s]: %#v\n", op, options[op])
		}
		fmt.Printf("arguments: %#v\n", arguments)
	}
	// Output:
	// option[help]: []string{""}
	// option[output]: []string{"data1", "data2"}
	// option[verbose]: []string{"", "", ""}
	// arguments: []string{"zero", "one", "two"}
}
