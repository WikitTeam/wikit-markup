package wikitm

import "testing"

func TestRenderPlainText(t *testing.T) {
	got := Render("ProjectWikit!")
	want := "ProjectWikit!"
	if got != want {
		t.Errorf("Render(%q) = %q; want %q", "ProjectWikit!", got, want)
	}
}
