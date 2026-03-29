package laslig

import (
	"bytes"
	"testing"
)

// TestDefaultLayoutBuilders verifies the default layout and its builder-style overrides.
func TestDefaultLayoutBuilders(t *testing.T) {
	layout := DefaultLayout().
		WithLeadingGap(0).
		WithBlockGap(3).
		WithSectionGap(4).
		WithSectionIndent(6).
		WithListMarker(ListMarkerBullet)

	if layout.leadingGap != 0 {
		t.Fatalf("layout.leadingGap = %d, want 0", layout.leadingGap)
	}
	if layout.blockGap != 3 {
		t.Fatalf("layout.blockGap = %d, want 3", layout.blockGap)
	}
	if layout.sectionGap != 4 {
		t.Fatalf("layout.sectionGap = %d, want 4", layout.sectionGap)
	}
	if layout.sectionIndent != 6 {
		t.Fatalf("layout.sectionIndent = %d, want 6", layout.sectionIndent)
	}
	if layout.listMarker != ListMarkerBullet {
		t.Fatalf("layout.listMarker = %q, want %q", layout.listMarker, ListMarkerBullet)
	}
}

// TestResolveLayout verifies default and custom layout resolution from policy values.
func TestResolveLayout(t *testing.T) {
	defaults := resolveLayout(Policy{})
	if defaults.leadingGap != 1 || defaults.blockGap != 1 || defaults.sectionGap != 2 || defaults.sectionIndent != 2 {
		t.Fatalf("resolveLayout(Policy{}) = %+v, want default layout values", defaults)
	}

	layout := Layout{}.WithLeadingGap(-3).WithListMarker("")
	resolved := resolveLayout(Policy{Layout: &layout})
	if resolved.leadingGap != 0 {
		t.Fatalf("resolved.leadingGap = %d, want 0", resolved.leadingGap)
	}
	if resolved.listMarker != ListMarkerDash {
		t.Fatalf("resolved.listMarker = %q, want %q", resolved.listMarker, ListMarkerDash)
	}
}

// TestGlamourStyleValidation verifies supported styles validate and invalid
// values fall back to the default preset.
func TestGlamourStyleValidation(t *testing.T) {
	if got := DefaultGlamourStyle(); got != GlamourStyleDracula {
		t.Fatalf("DefaultGlamourStyle() = %q, want %q", got, GlamourStyleDracula)
	}
	if !GlamourStyleDracula.Valid() {
		t.Fatal("GlamourStyleDracula.Valid() = false, want true")
	}
	if GlamourStyle("bogus").Valid() {
		t.Fatal("GlamourStyle(\"bogus\").Valid() = true, want false")
	}
	if got := resolveGlamourStyle(Policy{GlamourStyle: GlamourStyle("bogus")}); got != GlamourStyleDracula {
		t.Fatalf("resolveGlamourStyle(invalid) = %q, want %q", got, GlamourStyleDracula)
	}
	if got := resolveGlamourStyle(Policy{GlamourStyle: GlamourStyleTokyoNight}); got != GlamourStyleTokyoNight {
		t.Fatalf("resolveGlamourStyle(valid) = %q, want %q", got, GlamourStyleTokyoNight)
	}
}

// TestCustomListMarker verifies list-marker customization affects rendered output.
func TestCustomListMarker(t *testing.T) {
	var buf bytes.Buffer
	layout := DefaultLayout().WithLeadingGap(0).WithListMarker(ListMarkerNumber)
	printer := New(&buf, Policy{
		Format: FormatPlain,
		Style:  StyleNever,
		Layout: &layout,
	})

	if err := printer.List(List{
		Title: "Items",
		Items: []ListItem{
			{Title: "first"},
			{Title: "second"},
		},
	}); err != nil {
		t.Fatalf("List() error = %v", err)
	}

	want := "Items\n1. first\n2. second\n"
	if got := buf.String(); got != want {
		t.Fatalf("List() output = %q, want %q", got, want)
	}
}
