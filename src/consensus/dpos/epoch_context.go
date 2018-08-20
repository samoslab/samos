package dpos

import (
	"errors"

	"github.com/samoslab/samos/src/cipher"
)

type EpochContext struct {
	DposContext DposContext
	TimeStamp   int64
}

func NewEpochFromDposContext(dc DposContext, ts int64) *EpochContext {
	return &EpochContext{
		DposContext: dc,
		TimeStamp:   ts,
	}
}

func (ec *EpochContext) LookupValidator(now int64) (validator cipher.PubKey, err error) {
	validator = cipher.PubKey{}
	offset := now % epochInterval
	if offset%blockInterval != 0 {
		return cipher.PubKey{}, ErrInvalidMintBlockTime
	}
	offset /= blockInterval

	validators, err := ec.DposContext.GetValidators()
	if err != nil {
		return cipher.PubKey{}, err
	}
	validatorSize := len(validators)
	if validatorSize == 0 {
		return cipher.PubKey{}, errors.New("failed to lookup validator")
	}
	offset %= int64(validatorSize)
	return validators[offset], nil
}
