package utils

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TopicToMethodName(topic string) string {
	// Find the last dot in the topic
	lastDotIndex := strings.LastIndex(topic, ".")

	// If there is no dot, use the entire topic
	var lastPart string
  
	if lastDotIndex == -1 {
		lastPart = topic
	} else {
		// Take the last part after the last dot
		lastPart = topic[lastDotIndex+1:]
	}

	// Split the last part by hyphens
	parts := strings.Split(lastPart, "-")

	// Capitalize each part and concatenate them
	var capitalizedParts []string
	for _, part := range parts {
		capitalizedPart := cases.Title(language.Und).String(part)
		capitalizedParts = append(capitalizedParts, capitalizedPart)
	}

	// Join the capitalized parts and append "Handler"
	return strings.Join(capitalizedParts, "") + "Handler"
}
