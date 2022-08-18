package types

type Permission int

const (
	// Displays as: pv
	//
	// Other modules cannot access the section or member.
	PermissionPrivate Permission = iota

	// Displays as: ro
	//
	// Other modules can access the section or member, but can only read its
	// value. It is effectively immutable.
	// 
	// Data sections, member variables, etc: The value can be read by other
	// modules but not altered by them
	//
	// Functions: The function can be called by other modules.
	//
	// Methods: The method can be called by other modules, but cannot be
	// overriden by a type defined in another module inheriting from this
	// method's reciever.
	PermissionReadOnly

	// Displays as: rw
	//
	// Other modules cannot only access the section or member's value, but
	// can alter it. It is effectively mutable. 
	// 
	// Data sections, member variables, etc: The value can be read and
	// altered at will by other modules.
	//
	// Functions: This permission cannot be applied to non-method functions.
	//
	// Methods: The method can be called by other modules, and overridden by
	// types defined in other modules inheriting from the method's reciever.
	PermissionReadWrite
)

// PermissionFrom creates a new permission value from the specified text. If the
// input text was not valid, the function returns false for worked. Otherwise,
// it returns true.
func PermissionFrom (data string) (permission Permission, worked bool) {
	worked = true
	switch data {
	case "pv": permission = PermissionPrivate
	case "ro": permission = PermissionReadOnly
	case "rw": permission = PermissionReadWrite
	default:   worked = false
	}
	return
}

// ToString converts the permission value into a string.
func (permission Permission) ToString () (output string) {
	switch permission {
	case PermissionPrivate:   output = "pv"
	case PermissionReadOnly:  output = "ro"
	case PermissionReadWrite: output = "rw"
	}
	return
}
