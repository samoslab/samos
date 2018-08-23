package dpos

import (
	"errors"

	"github.com/samoslab/samos/src/cipher"
)

type DposContext struct {
	candidate []cipher.PubKey
}

func NewDposContext() *DposContext {
	return &DposContext{}
}

func (dc *DposContext) GetValidators() ([]cipher.PubKey, error) {
	if len(dc.candidate) == 0 {
		return dc.candidate, errors.New("zero validator")
	}
	return dc.candidate, nil
}

func (dc *DposContext) SetValidators(validators []cipher.PubKey) error {
	dc.candidate = validators
	return nil
}
