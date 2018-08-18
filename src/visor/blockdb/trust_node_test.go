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
