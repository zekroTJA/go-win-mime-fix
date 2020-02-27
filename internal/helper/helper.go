// package helper contains general
// helper functionalities
package helper

import (
	"fmt"
	"strings"

	"github.com/bclicn/color"
)

// isValidContentType returns true if the passed
// content type string ct starts with the expected
// "application/javascript" content type.
func IsValidContentType(ct string) bool {
	return strings.HasPrefix(ct, "application/javascript")
}

// PrintStatusLine prints a formatted check status
// line with a given description, the collected
// value of the check and the validity state of
// the check.
func PrintStatusLine(desc, val string, valid bool) {
	validS := color.BGreen("as expected")
	if !valid {
		validS = color.BRed("not as expected")
	}

	fmt.Printf("%s → '%s' ► %s\n",
		desc, color.Yellow(val), validS)
}
