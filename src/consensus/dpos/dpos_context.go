package dpos

import (
	"errors"

	"github.com/samoslab/samos/src/cipher"
)

// DposContext store trust pubkey
type DposContext struct {
	candidate []cipher.PubKey
}

// NewDposContext new instance
func NewDposContext() *DposContext {
	return &DposContext{}
}

// GetValidators return all validators
func (dc *DposContext) GetValidators() ([]cipher.PubKey, error) {
	if len(dc.candidate) == 0 {
		return dc.candidate, errors.New("zero validator")
	}
	return dc.candidate, nil
}

// SetValidators set valid validator
func (dc *DposContext) SetValidators(validators []cipher.PubKey) error {
	dc.candidate = validators
	return nil
}
