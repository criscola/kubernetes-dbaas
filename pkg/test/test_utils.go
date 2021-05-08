// Package test provides common utilities to be used in tests throughout the project
package test

import (
	"fmt"
	"strings"
)

type TestType string

const (
	Unit        TestType = "u"
	Integration TestType = "i"
	E2e         TestType = "e2e"
)

// FormatTestDesc should be used throughout the project to format test descriptions for the Ginkgo testing
// library. This is needed because it provides a higher-level abstraction for the format of test descriptions in
// order to allow --skip or --focus test flags to work uniformly.
// Example output: [e2e] RenderOperation [extra1] [extra2]
// TODO: The primitivity of Ginkgo's test runner is a long-standing issue of its 1.x version, see https://github.com/onsi/ginkgo/issues/664
// and https://github.com/onsi/ginkgo/issues/144. This function will likely become obsolete from the 2.x version of
// Ginkgo.
func FormatTestDesc(label TestType, description string, extras ...string) string {
	var extraStrings string
	if len(extras) > 0 {
		for _, s := range extras {
			extraStrings += "[" + s + "] "
		}
	}
	extraStrings = strings.TrimRight(extraStrings, " ")
	return fmt.Sprintf("[%s] %s %s", label, description, extraStrings)
}
