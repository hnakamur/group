package group

import (
	"os/user"
	"testing"
)

func TestLookup(t *testing.T) {
	name := "daemon"
	got, err := Lookup(name)
	if err != nil {
		t.Fatalf("Lookup: %v", err)
	}
	if got.Name != name {
		t.Errorf("got Name=%q; want %q", got.Name, name)
	}
}

func TestLookupId(t *testing.T) {
	u, err := user.Current()
	if err != nil {
		t.Fatalf("LookupId: %v", err)
	}
	gid := u.Gid
	got, err := LookupId(gid)
	if err != nil {
		t.Fatalf("LookupId: %v", err)
	}
	if got.Gid != gid {
		t.Errorf("got Gid=%q; want %q", got.Gid, gid)
	}
}
