package flagger

import (
	"fmt"
	"os"
)

var UsageDefault func()
var UsageFooter, UsageHeader string

func UsageWithHeaderAndFooter() {
	if len(UsageHeader) > 0 {
		fmt.Fprintln(os.Stderr, UsageHeader)
	}
	if UsageDefault != nil {
		UsageDefault()
	}
	if len(UsageFooter) > 0 {
		fmt.Fprintln(os.Stderr, UsageFooter)
	}
}
