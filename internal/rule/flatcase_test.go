package rule

import (
	"errors"
	"testing"
)

func TestFlatcase(t *testing.T) {
	rule := new(Flatcase).Init()

	tests := []*ruleTest{
		{name: "expected use case", value: "abc0123", expected: true},

		{name: "lowercase with space", value: "abc def", expected: false},
		{name: "lowercase with underscore", value: "abc_def", expected: false},
		{name: "lowercase with dash", value: "abc-def", expected: false},
		{name: "lowercase with dot", value: "abc.def", expected: false},
		{name: "lowercase with emoji", value: "abc😀def", expected: false},
		{name: "lowercase with uppercase", value: "abcDEF", expected: false},

		// edge cases
		{name: "empty string", value: "", expected: false},
		{name: "one digit", value: "1", expected: true},
		{name: "one lowercase", value: "a", expected: true},
		{name: "one lowercase with diacritics", value: "å", expected: true},
		{name: "one uppercase", value: "A", expected: false},
		{name: "one uppercase with diacritics", value: "Å", expected: false},
		{name: "digits only", value: "0123456789", expected: true},
		{name: "emoji", value: "😀", expected: false},

		// Lowercase from other languages
		{name: "diacritics on lowercase", value: "åäöéèêçñ", expected: true},
		{name: "Cyrillic lowercase", value: "абвгдежзийклмнопрстуфхцчшщъыьэюя", expected: true},
		{name: "Greek lowercase", value: "αβγδεζηθικλμνξοπρστυφχψω", expected: true},
		{name: "Armenian lowercase", value: "աաբգդեէզէթժիյկլհմնոօփքռստուֆք", expected: true},
		{name: "Georgian lowercase", value: "ააბგდევზთიკლმნოპჟრსტუფქღყშჩცძწჭხჯჰ", expected: true},

		// Uppercase from other languages
		{name: "diacritics on uppercase", value: "ÅÄÖÉÈÊÇÑ", expected: false},
		{name: "Cyrillic uppercase", value: "Ж", expected: false},
		{name: "Greek uppercase", value: "Ω", expected: false},
		{name: "Armenian uppercase", value: "Ֆ", expected: false},

		// The following languages don't have distinct uppercase and lowercase letters
		// so all of them are invalid
		{name: "Chinese letter", value: "汉", expected: false},
		{name: "Japanese letter", value: "漢", expected: false},
		{name: "Arabic letter", value: "ع", expected: false},
		{name: "Hebrew letter", value: "א", expected: false},
		{name: "Thai letter", value: "ภ", expected: false},
		{name: "Japanese Hiragana letter", value: "あ", expected: false},
		{name: "Japanese Katakana letter", value: "ア", expected: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := rule.Validate(test.value, "", true)

			if !errors.Is(err, test.err) {
				t.Errorf("unmatched error - %e", err)
				return
			}

			if res != test.expected {
				t.Errorf("unmatched return value - %+v", res)
				return
			}
		})
	}
}
