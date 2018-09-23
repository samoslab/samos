package pbft

import (
	"errors"
	"testing"

	"github.com/samoslab/samos/src/cipher"
	"github.com/samoslab/samos/src/coin"
	"github.com/stretchr/testify/assert"
)

func _feeCalc(t *coin.Transaction) (uint64, error) {
	return 0, nil
}
func makeNewBlock(uxHash cipher.SHA256) (*coin.Block, error) {
	body := coin.BlockBody{
		Transactions: coin.Transactions{coin.Transaction{}},
	}

	prev := coin.Block{
		Body: body,
		Head: coin.BlockHeader{
			Version:  0x02,
			Time:     100,
			BkSeq:    0,
			Fee:      10,
			PrevHash: cipher.SHA256{},
			BodyHash: body.Hash(),
		}}
	return coin.NewBlock(prev, 100+20, uxHash, coin.Transactions{coin.Transaction{}}, _feeCalc)
}
func makeNewBlock1(uxHash cipher.SHA256) (*coin.Block, error) {
	body := coin.BlockBody{
		Transactions: coin.Transactions{coin.Transaction{}},
	}

	prev := coin.Block{
		Body: body,
		Head: coin.BlockHeader{
			Version:  0x02,
			Time:     1000,
			BkSeq:    1,
			Fee:      100,
			PrevHash: uxHash,
			BodyHash: body.Hash(),
		}}
	return coin.NewBlock(prev, 1000+20, uxHash, coin.Transactions{coin.Transaction{}}, _feeCalc)
}

func TestPbft(t *testing.T) {
	pbft := NewPBFT()
	hashstr := "abcd1234"
	uxhash := cipher.SumSHA256([]byte(hashstr))
	pubkey := cipher.MustPubKeyFromHex("025d8360bc9439aa94044df96605f7693f50bb35386b37ae5003787d840a98bf43")
	seckey := cipher.MustSecKeyFromHex("4f36d5784d96a5b0e29d6876dd4eba422a2d92a29e81c67487fff7c403fa105b")
	pubkey1 := cipher.MustPubKeyFromHex("02d15bf28c4ed2c39b35b2be2f8bcde1318e2b3b65fe2a676db39b520bee9bfe86")
	pubkey2 := cipher.MustPubKeyFromHex("02e99a1338841e8b1f192337d2c6157045faa0cfe3b8a02210283aed7f5ad6880d")

	block, err := makeNewBlock(uxhash)
	assert.NoError(t, err)
	sig := cipher.SignHash(block.HashHeader(), seckey)
	sb := coin.SignedBlock{
		Block: *block,
		Sig:   sig,
	}

	err = pbft.AddSignedBlock(sb)
	assert.NoError(t, err)
	err = pbft.AddSignedBlock(sb)
	assert.Equal(t, err, errors.New("the block has added"))
	hash := sb.HashHeader()
	err = pbft.AddValidator(hash, pubkey)
	assert.Error(t, err)
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

	temphashstr := "abcd12345678"
	temphash := cipher.SumSHA256([]byte(temphashstr))
	err = pbft.AddValidator(temphash, pubkey2)
	assert.Equal(t, err, errors.New("this block hash not added into prepared infos"))

	block1, err := makeNewBlock1(temphash)
	assert.NoError(t, err)
	seckey1 := cipher.MustSecKeyFromHex("ebfb674103d75439360694cba709f522ab461d75914c7868204101927bf08b36")
	sig1 := cipher.SignHash(block.HashHeader(), seckey1)
	sb1 := coin.SignedBlock{
		Block: *block1,
		Sig:   sig1,
	}
	err = pbft.AddSignedBlock(sb1)
	assert.Equal(t, err, errors.New("has unconfirmed block, this block cannot added"))

	err = pbft.CheckPubkeyExists(hash, pubkey1)
	assert.NoError(t, err)

	_, err = pbft.GetSignedBlock(hash)
	assert.Nil(t, err)
	err = pbft.DeleteHash(hash)
	assert.Nil(t, err)
}
