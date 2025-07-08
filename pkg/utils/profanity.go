package utils

import "strings"

func ReplaceProfanity(text string) string {
	profaneWords := map[string]bool{"kerfuffle": true, "sharbert": true, "fornax": true}
	words := strings.Fields(text)

	for i, word := range words {
		if profaneWords[strings.ToLower(word)] {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}