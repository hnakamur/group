package group

import (
	"strconv"
)

var implemented = true // set to false by lookup_stubs.go's init

// Group represents a group.
//
// On posix systems Gid contain a decimal number
// representing gid.
type Group struct {
	Gid     string // group id
	Name    string
	Members []string // usernames who are members of this group as
	// a supplementary group, not a primary group.
}

// UnknownGroupIdError is returned by LookupId when
// a group cannot be found.
type UnknownGroupIdError int

func (e UnknownGroupIdError) Error() string {
	return "group: unknown groupid " + strconv.Itoa(int(e))
}

// UnknownGroupError is returned by Lookup when
// a group cannot be found.
type UnknownGroupError string

func (e UnknownGroupError) Error() string {
	return "group: unknown group " + string(e)
}
