package dpos

import (
	"errors"
	"sync"

	"github.com/samoslab/samos/src/cipher"
	"github.com/samoslab/samos/src/coin"
)

const (
	blockInterval = int64(10)
	epochInterval = int64(86400)
)

var (
	ErrInvalidTimestamp      = errors.New("invalid timestamp")
	ErrWaitForPrevBlock      = errors.New("wait for last block arrived")
	ErrMintFutureBlock       = errors.New("mint the future block")
	ErrBlockAlreadyCreated   = errors.New("block already created in the slot")
	ErrInvalidBlockValidator = errors.New("invalid block validator")
	ErrInvalidMintBlockTime  = errors.New("invalid time to mint the block")
)

// Dpos consensus alg
type Dpos struct {
	signer      cipher.PubKey
	mu          sync.RWMutex
	stop        chan bool
	dposContext *DposContext
	lastSlot    uint32
}

// NewDpos create dpos instance
func NewDpos(signer cipher.PubKey) *Dpos {
	return &Dpos{
		dposContext: NewDposContext(),
		signer:      signer,
		lastSlot:    0,
	}
}

// SetTrustNode set trust nodes for block creator
func (d *Dpos) SetTrustNode(trusts []cipher.PubKey) error {
	return d.dposContext.SetValidators(trusts)
}

func (d *Dpos) checkDeadline(lastBlock *coin.SignedBlock, now int64) error {
	prevSlot := PrevSlot(now)
	nextSlot := NextSlot(now)
	if int64(lastBlock.Time()) >= nextSlot {
		return ErrMintFutureBlock
	}
	if int64(lastBlock.Time()) >= prevSlot && int64(lastBlock.Time()) < nextSlot {
		return ErrBlockAlreadyCreated
	}
	if int64(lastBlock.Time()) < prevSlot {
		return nil
	}
	return ErrWaitForPrevBlock
}

// CheckValidator check current node create block or not, every node has same chance to
// make block, it is destined by time
func (d *Dpos) CheckValidator(lastBlock *coin.SignedBlock, now int64) error {
	if err := d.checkDeadline(lastBlock, now); err != nil {
		return err
	}
	epochContext := NewEpochFromDposContext(*d.dposContext, now)
	validator, err := epochContext.LookupValidator(PrevSlot(now))
	if err != nil {
		return err
	}
	if (validator == cipher.PubKey{} || validator != d.signer) {
		return ErrInvalidBlockValidator
	}
	return nil
}

// PrevSlot previos slot for create block
func PrevSlot(now int64) int64 {
	return int64((now-1)/blockInterval) * blockInterval
}

// NextSlot next slot for create block
func NextSlot(now int64) int64 {
	return int64((now+blockInterval-1)/blockInterval) * blockInterval
}

// GetValidator returns validator in the timestamp
func (d *Dpos) GetValidator(timestamp int64) (cipher.PubKey, error) {
	epochContext := NewEpochFromDposContext(*d.dposContext, timestamp)
	return epochContext.LookupValidator(PrevSlot(timestamp))
}
