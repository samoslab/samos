package dpos

import (
	"errors"
	"fmt"
	"sync"

	"github.com/samoslab/samos/src/cipher"
	"github.com/samoslab/samos/src/coin"
)

const (
	blockInterval = int64(10)
	epochInterval = int64(86400)
)

var (
	// errUnknownBlock is returned when the list of signers is requested for a block
	// that is not part of the local blockchain.
	errUnknownBlock = errors.New("unknown block")

	ErrInvalidTimestamp           = errors.New("invalid timestamp")
	ErrWaitForPrevBlock           = errors.New("wait for last block arrived")
	ErrMintFutureBlock            = errors.New("mint the future block")
	ErrBlockAlreadyCreated        = errors.New("block already created in the slot")
	ErrMismatchSignerAndValidator = errors.New("mismatch block signer and validator")
	ErrInvalidBlockValidator      = errors.New("invalid block validator")
	ErrInvalidMintBlockTime       = errors.New("invalid time to mint the block")
	ErrNilBlockHeader             = errors.New("nil block header returned")
)

type Dpos struct {
	signer      cipher.PubKey
	mu          sync.RWMutex
	stop        chan bool
	dposContext *DposContext
	lastSlot    uint32
}

func NewDpos(signer cipher.PubKey) *Dpos {
	return &Dpos{
		dposContext: NewDposContext(),
		signer:      signer,
		lastSlot:    0,
	}
}

func (d *Dpos) SetTrustNode(trusts []cipher.PubKey) error {
	return d.dposContext.SetValidators(trusts)
}

func (d *Dpos) checkDeadline(lastBlock *coin.SignedBlock, now int64) error {
	prevSlot := PrevSlot(now)
	nextSlot := NextSlot(now)
	fmt.Printf("prev %d, next %d, now %d, last %d\n", prevSlot, nextSlot, now, lastBlock.Time())
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

func (d *Dpos) CheckValidator(lastBlock *coin.SignedBlock, now int64) error {
	if err := d.checkDeadline(lastBlock, now); err != nil {
		return err
	}
	epochContext := NewEpochFromDposContext(*d.dposContext, now)
	validator, err := epochContext.LookupValidator(PrevSlot(now))
	if err != nil {
		return err
	}
	fmt.Printf("create block validator %s\n", validator.Hex())
	if (validator == cipher.PubKey{} || validator != d.signer) {
		return ErrInvalidBlockValidator
	}
	fmt.Printf("in turn the validator %s create block\n", validator.Hex())
	return nil
}

func PrevSlot(now int64) int64 {
	return int64((now-1)/blockInterval) * blockInterval
}

func NextSlot(now int64) int64 {
	return int64((now+blockInterval-1)/blockInterval) * blockInterval
}
