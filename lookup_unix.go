// +build darwin freebsd linux netbsd openbsd
// +build cgo

package group

import (
	"fmt"
	"runtime"
	"strconv"
	"syscall"
	"unsafe"
)

/*
#include <unistd.h>
#include <sys/types.h>
#include <grp.h>
#include <stdlib.h>

static int mygetgrgid_r(int gid, struct group *grp,
	char *buf, size_t buflen, struct group **result) {
 return getgrgid_r(gid, grp, buf, buflen, result);
}

static char **next_mem(char **mem) {
	return ++mem;
}

static int member_count(char **mem) {
	int c = 0;
	for (; *mem; ++mem) {
		++c;
	}
	return c;
}
*/
import "C"

func lookup(groupname string) (*Group, error) {
	return lookupUnix(-1, groupname, true)
}

func lookupId(gid string) (*Group, error) {
	i, e := strconv.Atoi(gid)
	if e != nil {
		return nil, e
	}
	return lookupUnix(i, "", false)
}

func lookupUnix(gid int, groupname string, lookupByName bool) (*Group, error) {
	var grp C.struct_group
	var result *C.struct_group

	var bufSize C.long
	if runtime.GOOS == "freebsd" {
		// FreeBSD doesn't have _SC_GETGR_R_SIZE_MAX
		// and just returns -1.  So just use the same
		// size that Linux returns
		// TODO: Confirm this!
		bufSize = 1024
	} else {
		bufSize = C.sysconf(C._SC_GETGR_R_SIZE_MAX)
		if bufSize <= 0 || bufSize > 1<<20 {
			return nil, fmt.Errorf("group: unreasonable _SC_GETGR_R_SIZE_MAX of %d", bufSize)
		}
	}
	buf := C.malloc(C.size_t(bufSize))
	defer C.free(buf)
	var rv C.int
	if lookupByName {
		nameC := C.CString(groupname)
		defer C.free(unsafe.Pointer(nameC))
		rv = C.getgrnam_r(nameC,
			&grp,
			(*C.char)(buf),
			C.size_t(bufSize),
			&result)
		if rv != 0 {
			return nil, fmt.Errorf("group: lookup groupname %s: %s", groupname, syscall.Errno(rv))
		}
		if result == nil {
			return nil, UnknownGroupError(groupname)
		}
	} else {
		// mygetgrgid_r is a wrapper around getgrgid_r to
		// avoid using gid_t because C.gid_t(gid) for
		// unknown reasons doesn't work on linux.
		rv = C.mygetgrgid_r(C.int(gid),
			&grp,
			(*C.char)(buf),
			C.size_t(bufSize),
			&result)
		if rv != 0 {
			return nil, fmt.Errorf("group: lookup groupid %d: %s", gid, syscall.Errno(rv))
		}
		if result == nil {
			return nil, UnknownGroupIdError(gid)
		}
	}

	var members []string
	if grp.gr_mem != nil {
		members = make([]string, C.member_count(grp.gr_mem))
		i := 0
		for mem := grp.gr_mem; *mem != nil; mem = C.next_mem(mem) {
			members[i] = C.GoString(*mem)
			i += 1
		}
	} else {
		members = make([]string, 0)
	}
	g := &Group{
		Gid:     strconv.Itoa(int(grp.gr_gid)),
		Name:    C.GoString(grp.gr_name),
		Members: members,
	}
	return g, nil
}
