// General utilities including utilities for unit testing
package utils

import (
	"encoding/json"
	"regexp"
	"strings"
)

var re = regexp.MustCompile("[^a-z0-9]+")

// Create a slug (all lowercase and replacing character groups matching "[^a-z0-9]+" with single hyphen)
func Slug(s string) string {
	return strings.Trim(re.ReplaceAllString(strings.ToLower(s), "-"), "-")
}

// Create a JSON string from an interface, or an empty string if it cannot be marshalled
func JsonStringify(data interface{}) string {

	raw, err := json.Marshal(data)

	if err != nil {
		return ""
	}

	return string(raw[:])
}

// Converts stack and panic message into JSON for readability on a single log line
func JsonStack(panicMsg interface{}, rawTrace []byte) string {

	msg, ok := panicMsg.(string)

	if !ok {
		msg = "Unprintable"
	}

	trace := strings.Replace(string(rawTrace), "\t", "", -1)

	lines := strings.Split(trace, "\n")

	traceData := struct {
		Panic string
		Stack []string
	}{
		Panic: msg,
		Stack: lines,
	}

	jsonData, err := json.Marshal(traceData)

	if err != nil {
		return "Panic:" + msg + ": " + trace
	}

	return string(jsonData)
}


