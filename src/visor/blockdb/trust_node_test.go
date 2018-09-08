package blockdb

import (
	"sort"
	"testing"

	"github.com/samoslab/samos/src/cipher"
	"github.com/samoslab/samos/src/testutil"
	"github.com/stretchr/testify/assert"
)

func TestAddNode(t *testing.T) {
	db, close := testutil.PrepareDB(t)
	defer close()
	trustNode, err := NewTrustNode(db)
	assert.Nil(t, err)

	testAddress := []string{"EX8omhDyjKtc8zHGp1KZwn7usCndaoJxSe", "2mBbNkm1pxv1Q8pYiA4n8v9zoUEWBEyKZ24"}
	sort.Strings(testAddress)
	addrList := []cipher.Address{}
	for _, addr := range testAddress {
		addrList = append(addrList, cipher.MustDecodeBase58Address(addr))
	}

	err = trustNode.AddNode(addrList)
	assert.Nil(t, err)

	trustAddrs := trustNode.GetNodes()
	assert.Equal(t, addrList, trustAddrs)

}

func TestAddPubkey(t *testing.T) {
	db, close := testutil.PrepareDB(t)
	defer close()
	trustNode, err := NewTrustNode(db)
	assert.Nil(t, err)

	testPks := []string{"03cec5e9f78524a4283868b79cf3a2b406bcd7956cd9b4be325e070a1cb1881563", "02d15bf28c4ed2c39b35b2be2f8bcde1318e2b3b65fe2a676db39b520bee9bfe86"}
	sort.Strings(testPks)
	pks := []cipher.PubKey{}
	for _, pk := range testPks {
		pks = append(pks, cipher.MustPubKeyFromHex(pk))
	}

	err = trustNode.AddNodePubkey(pks)
	assert.Nil(t, err)

	trustPks := trustNode.GetPubkeys()
	assert.Equal(t, pks, trustPks)
}
