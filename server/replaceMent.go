package main

import (
	"fmt"
	"strings"
)

func replacement(message string) string {
	message = strings.Replace(message, "#", "From:", -1)
	message = strings.Replace(message, "^", " To:", -1)
	message = strings.Replace(message, "$", " Amount:", -1)
	message = strings.Replace(message, "/", " Bank:", -1)
	message = strings.Replace(message, "-", " Time:", -1)
	//message = strings.Replace(message, "|", "", 1)

	parts := strings.Split(message, "|")

	// Initialize a new string to hold the result
	var result strings.Builder

	// Track packet count and iterate over each part
	packetCount := 1
	for _, part := range parts {
		if part != "" {
			result.WriteString(fmt.Sprintf("packet:%d %s ", packetCount, part))
			packetCount++
		}
	}

	// Trim any trailing spaces and print the result
	finalString := strings.TrimSpace(result.String())
	return finalString
}
