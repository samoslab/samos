package pbft

import (
	"errors"
	"sync"

	"github.com/samoslab/samos/src/cipher"
	"github.com/samoslab/samos/src/coin"
)

// PBFT pending block data
type PBFT struct {
	Status        int
	PendingBlocks map[cipher.SHA256]coin.SignedBlock
	PreparedInfos map[cipher.SHA256][]cipher.PubKey
	mutex         sync.Mutex
}

// NewPBFT new pbft
func NewPBFT() *PBFT {
	return &PBFT{
		PendingBlocks: make(map[cipher.SHA256]coin.SignedBlock, 1),
		Status:        0,
		PreparedInfos: make(map[cipher.SHA256][]cipher.PubKey, 1),
	}
}

// GetSignedBlock get SignedBlock for the hash
func (p *PBFT) GetSignedBlock(hash cipher.SHA256) (coin.SignedBlock, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if sb, ok := p.PendingBlocks[hash]; ok {
		return sb, nil
	}
	return coin.SignedBlock{}, errors.New("block not exists")
}

// DeleteHash delete hash from pending block map
func (p *PBFT) DeleteHash(hash cipher.SHA256) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.PendingBlocks[hash]; !ok {
		return errors.New("block hash not exists")
	}

	delete(p.PendingBlocks, hash)
	delete(p.PreparedInfos, hash)

	return nil
}

// WaitingConfirmedBlockHash block hash that waiting other validator message
func (p *PBFT) WaitingConfirmedBlockHash() []cipher.SHA256 {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	watingHash := []cipher.SHA256{}
	for hash := range p.PendingBlocks {
		watingHash = append(watingHash, hash)
	}
	return watingHash
}

// AddSignedBlock add a signed block
func (p *PBFT) AddSignedBlock(sb coin.SignedBlock) error {
	p.mutex.Lock()
	bh := sb.Block.HashHeader()
	if _, ok := p.PendingBlocks[bh]; ok {
		p.mutex.Unlock()
		return errors.New("the block has added")
	}
	p.PendingBlocks[bh] = sb

	pubkeyRec, err := cipher.PubKeyFromSig(sb.Sig, bh) //recovered pubkey
	if err != nil {
		p.mutex.Unlock()
		return errors.New("Invalid sig: PubKey recovery failed")
	}
	p.mutex.Unlock()
	return p.AddValidator(bh, pubkeyRec)
}

// GetBlockValidators returns all pubkeys for block hash
func (p *PBFT) GetBlockValidators(hash cipher.SHA256) ([]cipher.PubKey, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	validators, ok := p.PreparedInfos[hash]
	if !ok {
		return []cipher.PubKey{}, errors.New("not exists")
	}
	return validators, nil
}

// CheckPubkeyExists check pubkey exists for the block hash
func (p *PBFT) CheckPubkeyExists(hash cipher.SHA256, pubkey cipher.PubKey) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	validators, ok := p.PreparedInfos[hash]
	if !ok {
		return errors.New("not exists")
	}
	for _, pk := range validators {
		if pk == pubkey {
			return nil
		}
	}
	return errors.New("not exists")
}

// AddValidator add a validator , if validator number exceed threshold, then make block into chain
func (p *PBFT) AddValidator(hash cipher.SHA256, pubkey cipher.PubKey) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	validators, ok := p.PreparedInfos[hash]
	if !ok {
		validators = []cipher.PubKey{}
	}
	for _, pk := range validators {
		if pk == pubkey {
			return errors.New("the pubkey already exists")
		}
	}
	validators = append(validators, pubkey)
	p.PreparedInfos[hash] = validators
	return nil
}

// ValidatorNumber the nunber of validator for the block hash
func (p *PBFT) ValidatorNumber(hash cipher.SHA256) (int, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	validators, ok := p.PreparedInfos[hash]
	if !ok {
		return 0, errors.New("the hash not exists")
	}
	return len(validators), nil
}
