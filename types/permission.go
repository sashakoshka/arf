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
	case 'n': mode = ModeNone
	case 'r': mode = ModeRead
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

func (mode Mode) ToString () (output string) {
	switch mode {
	case ModeNone:  output = "n"
	case ModeRead:  output = "r"
	case ModeWrite: output = "w"
	}

	return
}

func (permission Permission) ToString () (output string) {
	output += permission.Internal.ToString()
	output += permission.External.ToString()
	return
}
