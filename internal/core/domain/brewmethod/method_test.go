package brewmethod

import "testing"

func TestMatching(t *testing.T) {
	ms := []Method{
		{ID: "espresso", Name: "Espresso", MinUM: 180, MaxUM: 380},
		{ID: "v60", Name: "V60", MinUM: 400, MaxUM: 700},
	}
	if got := Matching(ms, 550); len(got) != 1 || got[0] != "V60" {
		t.Fatalf("550 µm = %v", got)
	}
	if got := Matching(ms, 390); got == nil || len(got) != 0 {
		t.Fatalf("390 µm should be an empty, non-nil list, got %#v", got)
	}
	if got := Matching(ms, 380); len(got) != 1 {
		t.Fatalf("range bounds are inclusive, got %v", got)
	}
}
