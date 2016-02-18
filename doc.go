// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// Parses options of the forms:
//   --help           - set empty string
//   -h               - set empty string
//   -v               - set empty string
//   -v value         - set value
//   --output=value   - set value
//   --output value   - set value
//   --o value        - set value
//   --ovalue         - set value
//   --               - stop option parsing
//
// Notes:
//   --help value     - value is an argument, it is not assigned to "help"
//   --output -v      - is an error as value is missing for required argument
//   --output=-v      - is allowed, the value is "-v"
//
package getoptions
