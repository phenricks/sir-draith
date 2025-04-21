package events

import "strings"

// stringContainsAny verifica se uma string contém qualquer uma das substrings fornecidas
func stringContainsAny(s string, substrings []string) bool {
	for _, sub := range substrings {
		if stringContains(s, sub) {
			return true
		}
	}
	return false
}

// stringContains verifica se uma string contém uma substring (case insensitive)
func stringContains(s, substr string) bool {
	s, substr = strings.ToLower(s), strings.ToLower(substr)
	return strings.Contains(s, substr)
}
