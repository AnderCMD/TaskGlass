//go:build windows

package app

import (
	"testing"
	"time"
)

func TestPanicGuardAllowsBurstsUpToMax(t *testing.T) {
	g := newPanicGuard(3, time.Minute)
	for i := 1; i <= 3; i++ {
		if !g.allow() {
			t.Fatalf("panic %d of 3 should be tolerated", i)
		}
	}
	if g.allow() {
		t.Fatal("the 4th panic inside the window must stop the app")
	}
}

func TestPanicGuardForgetsOldPanics(t *testing.T) {
	g := newPanicGuard(1, 20*time.Millisecond)
	if !g.allow() {
		t.Fatal("first panic should be tolerated")
	}
	time.Sleep(40 * time.Millisecond)
	if !g.allow() {
		t.Fatal("a panic outside the window must not count")
	}
}
