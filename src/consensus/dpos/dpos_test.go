package dpos

import (
	"testing"

	"github.com/samoslab/samos/src/cipher"
	"github.com/samoslab/samos/src/coin"
	"github.com/stretchr/testify/assert"
)

func TestSlot(t *testing.T) {
	testCases := []struct {
		now int64
		ps  int64
		ns  int64
	}{
		{
			now: 12345,
			ps:  12340,
			ns:  12350,
		},
		{
			now: 12340,
			ps:  12330,
			ns:  12340,
		},
		{
			now: 12349,
			ps:  12340,
			ns:  12350,
		},
		{
			now: 12341,
			ps:  12340,
			ns:  12350,
		},
	}
	for _, cs := range testCases {
		p := PrevSlot(cs.now)
		assert.Equal(t, cs.ps, p)
		n := NextSlot(cs.now)
		assert.Equal(t, cs.ns, n)
	}
}

func oneBlock(ts uint64) *coin.SignedBlock {
	block := coin.SignedBlock{}
	block.Block = coin.Block{
		Head: coin.BlockHeader{
			BodyHash: cipher.SHA256{},
			Version:  uint32(123),
			PrevHash: cipher.SHA256{},
			Time:     ts,
			BkSeq:    1 + 1,
			Fee:      1,
			UxHash:   cipher.SHA256{},
		},
	}
	return &block
}

func TestCheckDeadline(t *testing.T) {
	ts := uint64(12345678)
	block := oneBlock(ts)
	pubkey := cipher.MustPubKeyFromHex("03cec5e9f78524a4283868b79cf3a2b406bcd7956cd9b4be325e070a1cb1881563")
	dpos := NewDpos(pubkey)
	testCases := []struct {
		now int64
		err error
	}{
		{
			now: 12345677,
			err: ErrBlockAlreadyCreated,
		},
		{
			now: 12345681,
			err: nil,
		},
		{
			now: 12345670,
			err: ErrMintFutureBlock,
		},
		{
			now: 12345678,
			err: ErrBlockAlreadyCreated,
		},
	}
	for _, cs := range testCases {
		err := dpos.checkDeadline(block, cs.now)
		assert.Equal(t, cs.err, err)
	}
}

func TestCheckValidator(t *testing.T) {
	ts := uint64(12345678)
	block := oneBlock(ts)
	now := int64(12345680)
	pubkey := cipher.MustPubKeyFromHex("03cec5e9f78524a4283868b79cf3a2b406bcd7956cd9b4be325e070a1cb1881563")
	dpos := NewDpos(pubkey)
	pubkeys := []string{"03cec5e9f78524a4283868b79cf3a2b406bcd7956cd9b4be325e070a1cb1881563", "02d15bf28c4ed2c39b35b2be2f8bcde1318e2b3b65fe2a676db39b520bee9bfe86", "02e99a1338841e8b1f192337d2c6157045faa0cfe3b8a02210283aed7f5ad6880d"}
	trusts := []cipher.PubKey{}
	for _, pk := range pubkeys {
		trusts = append(trusts, cipher.MustPubKeyFromHex(pk))
	}
	dpos.SetTrustNode(trusts)
	err := dpos.CheckValidator(block, now)
	assert.Equal(t, ErrBlockAlreadyCreated, err)

	now = int64(12345681)
	err = dpos.CheckValidator(block, now)
	assert.Equal(t, err, ErrInvalidBlockValidator)

	now = int64(12345691)
	err = dpos.CheckValidator(block, now)
	assert.NoError(t, err)

	pubkeyValidator, err := dpos.GetValidator(int64(ts))
	assert.NoError(t, err)
	assert.Equal(t, pubkeyValidator, trusts[1])
}
