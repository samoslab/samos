package dpos

import (
	"errors"

	"github.com/samoslab/samos/src/cipher"
)

// EpochContext epoch for time polling
type EpochContext struct {
	DposContext DposContext
	TimeStamp   int64
}

// NewEpochFromDposContext new instance
func NewEpochFromDposContext(dc DposContext, ts int64) *EpochContext {
	return &EpochContext{
		DposContext: dc,
		TimeStamp:   ts,
	}
}

func calSolt(now int64) (int64, error) {
	offset := now % epochInterval
	if offset%blockInterval != 0 {
		return 0, ErrInvalidMintBlockTime
	}
	offset /= blockInterval
	return offset, nil
}

// LookupValidator lookup a valid validator according to time
func (ec *EpochContext) LookupValidator(now int64) (validator cipher.PubKey, err error) {
	validator = cipher.PubKey{}
	offset, err := calSolt(now)
	if err != nil {
		return cipher.PubKey{}, err
	}

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
