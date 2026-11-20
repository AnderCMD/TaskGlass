package i18n

import "testing"

func TestEveryKeyIsTranslated(t *testing.T) {
	for _, k := range Keys(English) {
		if _, ok := es[k]; !ok {
			t.Errorf("key %q has no Spanish translation", k)
		}
	}
	for _, k := range Keys(Spanish) {
		if _, ok := en[k]; !ok {
			t.Errorf("key %q exists only in Spanish", k)
		}
	}
}

func TestResolve(t *testing.T) {
	cases := []struct {
		l    Lang
		loc  string
		want Lang
	}{
		{Auto, "es-MX", Spanish}, {Auto, "ES-es", Spanish}, {Auto, "en-US", English},
		{Auto, "", English}, {Spanish, "en-US", Spanish}, {English, "es-MX", English},
	}
	for _, c := range cases {
		if got := Resolve(c.l, c.loc); got != c.want {
			t.Errorf("Resolve(%q,%q)=%q want %q", c.l, c.loc, got, c.want)
		}
	}
}

func TestFallback(t *testing.T) {
	if T(Spanish, "no.such.key") != "no.such.key" {
		t.Error("missing key should return the key")
	}
}
