package utils

import (
	"regexp"
	"strings"
)

var (
	bulletRegex = regexp.MustCompile(`(?m)^\s*[-*•]\s+`)
	numberRegex = regexp.MustCompile(`(?m)^\s*\d+\.\s+`)
)

func ParseBreakdownPoints(input string) []string {
	if input == "" {
		return []string{}
	}

	lines := strings.Split(input, "\n")
	result := make([]string, 0, len(lines))

	for _, line := range lines {
		l := strings.TrimSpace(line)
		if l == "" {
			continue
		}

		l = bulletRegex.ReplaceAllString(l, "")

		l = numberRegex.ReplaceAllString(l, "")

		l = strings.TrimSpace(l)

		if l == "" {
			continue
		}

		result = append(result, l)
	}

	return result
}
