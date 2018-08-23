package pbft

import (
	"errors"

	"github.com/samoslab/samos/src/cipher"
	"github.com/samoslab/samos/src/coin"
)

type PBFT struct {
	Status        int
	PendingBlocks map[cipher.SHA256]coin.SignedBlock
	PreparedInfos map[cipher.SHA256][]cipher.PubKey
}

func NewPBFT() *PBFT {
	return &PBFT{
		PendingBlocks: make(map[cipher.SHA256]coin.SignedBlock, 1),
		Status:        0,
		PreparedInfos: make(map[cipher.SHA256][]cipher.PubKey, 1),
	}
}

func (p *PBFT) GetSignedBlock(hash cipher.SHA256) (coin.SignedBlock, error) {
	if sb, ok := p.PendingBlocks[hash]; ok {
		return sb, nil
	}
	return coin.SignedBlock{}, errors.New("block not exists")
}

func (p *PBFT) AddSignedBlock(sb coin.SignedBlock) error {
	bh := sb.Block.HashHeader()
	if _, ok := p.PendingBlocks[bh]; ok {
		return errors.New("the block has added")
	}

	pubkeyRec, err := cipher.PubKeyFromSig(sb.Sig, bh) //recovered pubkey
	if err != nil {
		return errors.New("Invalid sig: PubKey recovery failed")
	}

	return p.AddValidator(bh, pubkeyRec)
}

func (p *PBFT) AddValidator(hash cipher.SHA256, pubkey cipher.PubKey) error {
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

func (p *PBFT) ValidatorNumber(hash cipher.SHA256) (int, error) {
	validators, ok := p.PreparedInfos[hash]
	if !ok {
		return 0, errors.New("the hash not exists")
	}
	return len(validators), nil
}
