package data

import (
	"encoding/json"
	"crypto/sha1"
)

type ETagger interface {
	ETagUpdate() error
	ETagGet() []byte
	ETagCompare([]byte) bool
}

type ETag struct {
	eTag []byte `json:"-"`
}

func (e *ETag) ETagUpdate() error {
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	t := sha1.Sum(data)
	e.eTag = t[:]

	return nil
}

func (e *ETag) ETagCompare(tag []byte) bool {
	if e.eTag == nil {
		return true
	}

	if len(tag) != 20 {
		return false
	}

	for i, b := range e.eTag {
		if tag[i] != b {
			return false
		}
	}

	return true
}

func (e *ETag) ETagGet() []byte {
	return e.eTag
}
