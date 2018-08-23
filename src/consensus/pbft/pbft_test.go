package pbft

import (
	"errors"
	"testing"

	"github.com/samoslab/samos/src/cipher"
	"github.com/stretchr/testify/assert"
)

func TestPbft(t *testing.T) {
	pbft := NewPBFT()
	hashstr := "abcd1234"
	hash := cipher.SumSHA256([]byte(hashstr))
	pubkey := cipher.MustPubKeyFromHex("03cec5e9f78524a4283868b79cf3a2b406bcd7956cd9b4be325e070a1cb1881563")
	pubkey1 := cipher.MustPubKeyFromHex("02d15bf28c4ed2c39b35b2be2f8bcde1318e2b3b65fe2a676db39b520bee9bfe86")
	pubkey2 := cipher.MustPubKeyFromHex("02e99a1338841e8b1f192337d2c6157045faa0cfe3b8a02210283aed7f5ad6880d")
	err := pbft.AddValidator(hash, pubkey)
	assert.NoError(t, err)
	err = pbft.AddValidator(hash, pubkey1)
	assert.NoError(t, err)
	num, err := pbft.ValidatorNumber(hash)
	assert.NoError(t, err)
	assert.Equal(t, 2, num)
	err = pbft.AddValidator(hash, pubkey1)
	assert.Equal(t, errors.New("the pubkey already exists"), err)
	err = pbft.AddValidator(hash, pubkey2)
	num, err = pbft.ValidatorNumber(hash)
	assert.NoError(t, err)
	assert.Equal(t, 3, num)
	assert.NoError(t, err)
}
