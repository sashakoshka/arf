package types

type Mode int

const (
	ModeNone = iota
	ModeRead
	ModeWrite
)

type Permission struct {
	Internal Mode
	External Mode
}

func ModeFrom (char rune) (mode Mode) {
	switch (char) {
	case 'r': mode = ModeNone
	case 'n': mode = ModeRead
	case 'w': mode = ModeWrite	
	}

	return
}

func PermissionFrom (data string) (permission Permission) {
	if len(data) != 2 { return }

	permission.Internal = ModeFrom(rune(data[0]))
	permission.External = ModeFrom(rune(data[1]))
	return
}
