package dpos

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/samoslab/samos/src/cipher"
	"github.com/samoslab/samos/src/coin"
)

const (
	extraVanity        = 32   // Fixed number of extra-data prefix bytes reserved for signer vanity
	extraSeal          = 65   // Fixed number of extra-data suffix bytes reserved for signer seal
	inmemorySignatures = 4096 // Number of recent block signatures to keep in memory

	blockInterval    = int64(10)
	epochInterval    = int64(86400)
	maxValidatorSize = 21
	safeSize         = maxValidatorSize*2/3 + 1
	consensusSize    = maxValidatorSize*2/3 + 1
)

var (
	big0  = big.NewInt(0)
	big8  = big.NewInt(8)
	big32 = big.NewInt(32)

	frontierBlockReward  *big.Int = big.NewInt(5e+18) // Block reward in wei for successfully mining a block
	byzantiumBlockReward *big.Int = big.NewInt(3e+18) // Block reward in wei for successfully mining a block upward from Byzantium

	timeOfFirstBlock = int64(0)

	confirmedBlockHead = []byte("confirmed-block-head")
)

var (
	// errUnknownBlock is returned when the list of signers is requested for a block
	// that is not part of the local blockchain.
	errUnknownBlock = errors.New("unknown block")
	// errMissingVanity is returned if a block's extra-data section is shorter than
	// 32 bytes, which is required to store the signer vanity.
	errMissingVanity = errors.New("extra-data 32 byte vanity prefix missing")
	// errMissingSignature is returned if a block's extra-data section doesn't seem
	// to contain a 65 byte secp256k1 signature.
	errMissingSignature = errors.New("extra-data 65 byte suffix signature missing")
	// errInvalidMixDigest is returned if a block's mix digest is non-zero.
	errInvalidMixDigest = errors.New("non-zero mix digest")
	// errInvalidUncleHash is returned if a block contains an non-empty uncle list.
	errInvalidUncleHash  = errors.New("non empty uncle hash")
	errInvalidDifficulty = errors.New("invalid difficulty")

	// ErrInvalidTimestamp is returned if the timestamp of a block is lower than
	// the previous block's timestamp + the minimum block period.
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
	signer      cipher.Address
	mu          sync.RWMutex
	stop        chan bool
	dposContext *DposContext
}

func NewDpos() *Dpos {
	return &Dpos{
		dposContext: NewDposContext(),
	}
}

func (d *Dpos) SetTrustNode(trusts []cipher.Address) error {
	return d.dposContext.SetValidators(trusts)
}

func (d *Dpos) SetSigner(signer cipher.Address) {
	d.signer = signer
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
	// last block was arrived, or time's up
	if int64(lastBlock.Time()) == prevSlot || nextSlot-now <= 1 {
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
	fmt.Printf("validator %s\n", validator.String())
	if (validator == cipher.Address{} || bytes.Compare(validator.Bytes(), d.signer.Bytes()) != 0) {
		return ErrInvalidBlockValidator
	}
	return nil
}

func PrevSlot(now int64) int64 {
	return int64((now-1)/blockInterval) * blockInterval
}

func NextSlot(now int64) int64 {
	return int64((now+blockInterval-1)/blockInterval) * blockInterval
}
