package textnorm

import "testing"

func TestKey(t *testing.T) {
	cases := map[string]string{
		"Mahlkönig X54":      "mahlkonigx54",
		"mahlkoenig x54":     "mahlkonigx54",
		"Virtuoso+":          "virtuosoplus",
		" Chestnut C3S-Pro ": "chestnutc3spro",
		"Cafetera Italiana":  "cafeteraitaliana",
		"":                   "",
		"---":                "",
	}
	for in, want := range cases {
		if got := Key(in); got != want {
			t.Errorf("Key(%q) = %q, want %q", in, got, want)
		}
	}
}
