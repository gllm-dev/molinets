package textnorm

import (
	"strings"
	"unicode"
)

var folds = map[rune]string{
	'ä': "a", 'à': "a", 'á': "a", 'â': "a",
	'ö': "o", 'ò': "o", 'ó': "o", 'ô': "o",
	'ü': "u", 'ù': "u", 'ú': "u",
	'è': "e", 'é': "e", 'ê': "e", 'ë': "e",
	'ì': "i", 'í': "i", 'ï': "i",
	'ç': "c", 'ñ': "n",
	'+': "plus",
}

func Key(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		if f, ok := folds[r]; ok {
			b.WriteString(f)
			continue
		}
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return strings.ReplaceAll(b.String(), "oe", "o")
}
