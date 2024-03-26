package transform

import (
	"strings"
	"unicode"
)

func UpperCamelToSnake(s string) string {
	var result strings.Builder
	result.WriteRune(unicode.ToLower(rune(s[0])))

	for _, char := range s[1:] {
		if unicode.IsUpper(char) {
			result.WriteRune('_')
			result.WriteRune(unicode.ToLower(char))
		} else {
			result.WriteRune(char)
		}
	}

	return result.String()
}

func SingularToPlural(s string) string {
	if s == "" {
		return s
	}

	// 例外的な変換ルール
	irregularForms := map[string]string{
		"person":    "people",
		"child":     "children",
		"ox":        "oxen",
		"man":       "men",
		"woman":     "women",
		"tooth":     "teeth",
		"foot":      "feet",
		"goose":     "geese",
		"cactus":    "cacti",
		"fungus":    "fungi",
		"focus":     "foci",
		"datum":     "data",
		"medium":    "media",
		"analysis":  "analyses",
		"basis":     "bases",
		"diagnosis": "diagnoses",
		"ellipsis":  "ellipses",
	}
	if val, ok := irregularForms[s]; ok {
		return val
	}

	// 通常の変換ルール
	if strings.HasSuffix(s, "y") && len(s) > 1 && !strings.ContainsAny(string(s[len(s)-2]), "aeiouy") {
		return s[:len(s)-1] + "ies"
	} else if strings.HasSuffix(s, "s") {
		return s + "es"
	}

	return s + "s"
}
