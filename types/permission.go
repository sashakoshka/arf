package types

type Mode int

const (
	ModeRead = iota
	ModeWrite
	ModeNone
)

type Permission struct {
	Internal Mode
	External Mode
}
