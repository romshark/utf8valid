package utf8_test

import (
	_ "embed"
	"flag"
	"testing"

	std "unicode/utf8"

	optimized "github.com/romshark/utf8valid"
)

//go:embed testdata/long_ascii.txt
var longASCII string

//go:embed testdata/ukranian_poetry.txt
var allUTF8UkranianPoetry string

//go:embed testdata/wikipedia_diacritic.html
var wikipediaDiacriticHTML string

//go:embed testdata/wikipedia_japan.html
var wikipediaJapanHTML string

var tests = []struct {
	name, input string
}{
	{"empty", ""},
	{"single_byte", "x"},
	{"single_utf8_rune", "ジ"},
	{"short_ascii", "Lorem Ipsum"},
	{"short_utf8", "ジабв🤯гдеёжзジ"},
	{"all_utf8_ukranian_poetry", allUTF8UkranianPoetry},
	{"long_ascii", longASCII},
	{"wikipedia_diacritic_html", wikipediaDiacriticHTML},
	{"wikipedia_japan_html", wikipediaJapanHTML},
	{"invalid_surrogate_max", "\xed\xbf\xbf\x80"},
}

func FuzzValid(f *testing.F) {
	for _, td := range tests {
		f.Add(td.input)
	}
	f.Fuzz(func(t *testing.T, s string) {
		b := []byte(s)
		opt, std := optimized.Valid(b), std.Valid(b)
		if opt != std {
			t.Errorf("expected: %t; received: %t", std, opt)
		}
	})
}

func FuzzValidString(f *testing.F) {
	for _, td := range tests {
		f.Add(td.input)
	}
	f.Fuzz(func(t *testing.T, s string) {
		opt, std := optimized.ValidString(s), std.ValidString(s)
		if opt != std {
			t.Errorf("expected: %t; received: %t", std, opt)
		}
	})
}

func TestBenchValid(t *testing.T) {
	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			bytesInput := []byte(td.input)
			if optimized.Valid(bytesInput) != std.Valid(bytesInput) {
				t.Errorf("unexpected mismatch for bytes input: %q", td.name)
			}
			if optimized.ValidString(td.input) != std.ValidString(td.input) {
				t.Errorf("unexpected mismatch for string input: %q", td.name)
			}
		})
	}
}

var fBenchFunc = flag.String(
	"benchfunc", "std",
	`benchmark function (either "std" or "optimized")`,
)

func getFunc[T any](tb testing.TB, std, optimized T) (f T) {
	switch *fBenchFunc {
	case "std":
		return std
	case "optimized":
		return optimized
	}
	tb.Fatalf("unsupported benchmark function: %q", *fBenchFunc)
	return
}

func BenchmarkFnValidString(b *testing.B) {
	f := getFunc(b, std.ValidString, optimized.ValidString)
	for _, td := range tests {
		b.Run(td.name, func(b *testing.B) {
			for range b.N {
				f(td.input)
			}
		})
	}
}

func BenchmarkFnValid(b *testing.B) {
	f := getFunc(b, std.Valid, optimized.Valid)
	for _, td := range tests {
		input := []byte(td.input)
		b.Run(td.name, func(b *testing.B) {
			for range b.N {
				f(input)
			}
		})
	}
}
