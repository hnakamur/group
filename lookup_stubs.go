// +build !cgo windows

package group

import (
	"fmt"
	"runtime"
)

func init() {
	implemented = false
}

func lookup(groupname string) (*Group, error) {
	return nil, fmt.Errorf("group: Lookup not implemented on %s/%s", runtime.GOOS, runtime.GOARCH)
}

func lookupId(gid string) (*Group, error) {
	return nil, fmt.Errorf("group: LookupId not implemented on %s/%s", runtime.GOOS, runtime.GOARCH)
}
