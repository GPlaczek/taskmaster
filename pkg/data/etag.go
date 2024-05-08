package data

type ETag interface {
	ETagUpdate() error
	ETagGet() []byte
	ETagCompare([]byte) bool
}
