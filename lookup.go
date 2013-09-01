
package group

// Lookup looks up a group by groupname. If the group cannot be found, the
// returned error is of type UnknownGroupError.
func Lookup(groupname string) (*Group, error) {
	return lookup(groupname)
}

// LookupId looks up a group by groupid. If the group cannot be found, the
// returned error is of type UnknownGroupIdError.
func LookupId(groupid string) (*Group, error) {
	return lookupId(groupid)
}
