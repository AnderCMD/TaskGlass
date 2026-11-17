package state

import (
	"testing"

	"github.com/AnderCMD/TaskGlass/internal/config"
)

func TestResolve(t *testing.T) {
	base := config.Default()
	base.Maximized.Enabled = true
	base.StartMenu.Enabled = true

	noMax := base
	noMax.Maximized.Enabled = false
	noStart := base
	noStart.StartMenu.Enabled = false

	tests := []struct {
		name string
		cfg  config.Config
		ctx  Context
		want Kind
	}{
		{"idle", base, Context{}, KindDefault},
		{"maximized", base, Context{Maximized: true}, KindMaximized},
		{"start", base, Context{StartMenu: true}, KindStartMenu},
		{"start beats maximized", base, Context{StartMenu: true, Maximized: true}, KindStartMenu},
		{"maximized rule disabled", noMax, Context{Maximized: true}, KindDefault},
		{"start rule disabled falls to maximized", noStart, Context{StartMenu: true, Maximized: true}, KindMaximized},
		{"start rule disabled alone", noStart, Context{StartMenu: true}, KindDefault},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got := Resolve(tt.cfg, tt.ctx)
			if got != tt.want {
				t.Fatalf("got %s want %s", got, tt.want)
			}
		})
	}
}

func TestResolveReturnsMatchingAppearance(t *testing.T) {
	c := config.Default()
	c.Maximized.Enabled = true
	ap, _ := Resolve(c, Context{Maximized: true})
	if ap != c.Maximized.Appearance {
		t.Fatalf("got %+v", ap)
	}
}
