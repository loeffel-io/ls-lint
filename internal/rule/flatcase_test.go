package rule

import (
	"testing"
)

func TestFlatcase(t *testing.T) {
	rule := new(Flatcase).Init()

	validTests := map[string]string{
		"expected use case":                  "abc0123",
		"lowercase letters":                  "abcdefghijklmnopqrstuvwxyz",
		"lowercase letters with digits":      "0123456789",
		"one lowercase letter":               "a",
		"one digit":                          "1",
		"lower case letters with diacritics": "åäöéèêçñ",
		"lower case cyrillic letters":        "абвгдежзийклмнопрстуфхцчшщъыьэюя",
		"lower case greek letters":           "αβγδεζηθικλμνξοπρστυφχψω", // Greek alpha lower case letter
		"lower case armenian letters":        "աաբգդեէզէթժիյկլհմնոօփքռստուֆք",
	}

	for name, value := range validTests {
		valid, err := rule.Validate(value, "", true)
		if err != nil {
			t.Errorf("Unexpected error when validating %s (%q): %v", name, value, err)
		}

		if !valid {
			t.Errorf("Expected %s (%q) to be valid", name, value)
		}
	}

	invalidCharacters := map[string]string{
		"uppercase":                 "A",
		"uppercase with diacritics": "Å", // Å is Angstrom, not a regular A

		// These languages distinct uppercase and lowercase letters
		// so any uppercase letter is invalid
		"cyrillic uppercase": "Ж",
		"Greek uppercase":    "Ω",
		"Armenian uppercase": "Ֆ",

		"underscore": "_",
		"dash":       "-",
		"dot":        ".",
		"space":      " ",
		"emoji":      "😀",

		// The following languages don't have distinct uppercase and lowercase letters
		// so all of them are invalid
		"Chinese letter":           "汉",
		"Japanese letter":          "漢",
		"Arabic letter":            "ع",
		"Hebrew letter":            "א",
		"Thai letter":              "ภ",
		"Japanese Hiragana letter": "あ",
		"Japanese Katakana letter": "ア",
	}

	validCharacter := map[string]string{
		"lowercase letters":                 "a",
		"lowercase letters with diacritics": "å",
		"digit":                             "1",
		"lower case greek letter":           "α", // Greek alpha lower case letter
	}

	invalidTests := map[string]string{
		"empty string": "",
	}
	// let's populate the tests with all combinations of invalid characters
	for name, character := range invalidCharacters {
		invalidTests[name] = character // add without any variation

		// check that adding another invalid characters don't allow them to be valid
		for variationName, variationCharacter := range invalidCharacters {
			invalidTests[name+" followed by "+variationName] = character + variationCharacter
			invalidTests[variationName+" followed by "+name] = variationCharacter + character
		}

		// check that adding valid characters don't allow them to be valid
		for variationName, variationCharacter := range validCharacter {
			invalidTests[name+" followed by "+variationName] = character + variationCharacter
			invalidTests[variationName+" followed by "+name] = variationCharacter + character
		}
	}

	for name, value := range invalidTests {
		valid, err := rule.Validate(value, "", true)
		if err != nil {
			t.Errorf("Unexpected error when validating %s (%q): %v", name, value, err)
			continue
		}

		if valid {
			t.Errorf("Expected %s (%q) to be invalid", name, value)
		}
	}
}
