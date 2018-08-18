package blockdb

import (
	"fmt"
	"sort"
	"strings"

	"github.com/boltdb/bolt"

	"github.com/samoslab/samos/src/cipher"
	"github.com/samoslab/samos/src/visor/bucket"
)

// TrustNode use the trustnode store all trust node info
type TrustNode struct {
	db   *bolt.DB
	node *bucket.Bucket
}

// NewBlockTree create buckets in blockdb if does not exist.
func NewTrustNode(db *bolt.DB) (*TrustNode, error) {
	node, err := bucket.New([]byte("trust_node"), db)
	if err != nil {
		return nil, err
	}

	return &TrustNode{
		node: node,
		db:   db,
	}, nil
}

// AddNode write the node into blocks trust_node
func (bt *TrustNode) AddNode(addresses []cipher.Address) error {
	return bt.db.Update(func(tx *bolt.Tx) error {
		return bt.AddNodeWithTx(tx, addresses)
	})
}

// AddNodeWithTx adds block with *bolt.Tx
func (bt *TrustNode) AddNodeWithTx(tx *bolt.Tx, addresses []cipher.Address) error {
	bkt := tx.Bucket(bt.node.Name)
	if bkt == nil {
		return fmt.Errorf("bucket %s doesn't eist", bt.node.Name)
	}

	trustAddrs := []string{}
	for _, addr := range addresses {
		trustAddrs = append(trustAddrs, addr.String())
	}

	sort.Strings(trustAddrs)

	return bkt.Put([]byte("addresses"), []byte(strings.Join(trustAddrs, ",")))
}

// GetNodes get all trust nodes
func (bt *TrustNode) GetNodes() []cipher.Address {
	return bt.getBlock()
}

func (bt *TrustNode) getBlock() []cipher.Address {
	addresses := []cipher.Address{}
	bin := bt.node.Get([]byte("addresses"))
	if bin == nil {
		return addresses
	}
	for _, addr := range strings.Split(string(bin), ",") {
		addresses = append(addresses, cipher.MustDecodeBase58Address(addr))
	}
	return addresses
}
