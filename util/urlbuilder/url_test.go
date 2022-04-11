package urlbuilder

import "testing"

func TestJoin(t *testing.T) {
	u := Join("synod.dev", "users").AddQuery("version", "1.1.0")
	addr := u.Build()
	expected := "http://synod.dev/users?version=1.1.0"

	if addr != expected {
		t.Errorf("expected %s, got %s\n", expected, addr)
	}
}

func TestNew(t *testing.T) {
	u := New().Schema("https").Host("synod.dev").Path("users/1").AddQuery("page", "1")
	addr := u.Build()
	expected := "https://synod.dev/users/1?page=1"

	if addr != expected {
		t.Errorf("expected %s, got %s\n", expected, addr)
	}
}

func TestOf(t *testing.T) {
	baseURL := "synod.dev/users?perSize=1"
	expected := "synod.dev/users?perSize=1&parsed=1"

	u := Of(baseURL).AddQuery("parsed", "1")
	addr := u.Build()

	if addr != expected {
		t.Errorf("expected %s, got %s\n", expected, addr)
	}
}
