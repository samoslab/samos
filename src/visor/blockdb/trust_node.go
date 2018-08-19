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
func (tn *TrustNode) AddNode(addresses []cipher.Address) error {
	return tn.db.Update(func(tx *bolt.Tx) error {
		return tn.AddNodeWithTx(tx, addresses)
	})
}

// AddNodeWithTx adds block with *bolt.Tx
func (tn *TrustNode) AddNodeWithTx(tx *bolt.Tx, addresses []cipher.Address) error {
	bkt := tx.Bucket(tn.node.Name)
	if bkt == nil {
		return fmt.Errorf("bucket %s doesn't eist", tn.node.Name)
	}

	trustAddrs := []string{}
	for _, addr := range addresses {
		trustAddrs = append(trustAddrs, addr.String())
	}

	sort.Strings(trustAddrs)

	return bkt.Put([]byte("addresses"), []byte(strings.Join(trustAddrs, ",")))
}

// GetNodes get all trust nodes
func (tn *TrustNode) GetNodes() []cipher.Address {
	return tn.getNode()
}

func (tn *TrustNode) getNode() []cipher.Address {
	addresses := []cipher.Address{}
	bin := tn.node.Get([]byte("addresses"))
	if bin == nil {
		return addresses
	}
	for _, addr := range strings.Split(string(bin), ",") {
		addresses = append(addresses, cipher.MustDecodeBase58Address(addr))
	}
	return addresses
}

// AddNodePubkey write the node into blocks trust_node
func (tn *TrustNode) AddNodePubkey(pubkeys []cipher.PubKey) error {
	return tn.db.Update(func(tx *bolt.Tx) error {
		return tn.AddNodePubkeyWithTx(tx, pubkeys)
	})
}

// AddNodePubkeyWithTx adds block with *bolt.Tx
func (tn *TrustNode) AddNodePubkeyWithTx(tx *bolt.Tx, pubkeys []cipher.PubKey) error {
	bkt := tx.Bucket(tn.node.Name)
	if bkt == nil {
		return fmt.Errorf("bucket %s doesn't eist", tn.node.Name)
	}

	trustPks := []string{}
	for _, pk := range pubkeys {
		trustPks = append(trustPks, pk.Hex())
	}

	sort.Strings(trustPks)

	return bkt.Put([]byte("pubkey"), []byte(strings.Join(trustPks, ",")))
}

// GetPubkeys get all trust nodes
func (tn *TrustNode) GetPubkeys() []cipher.PubKey {
	return tn.getPubkeys()
}

func (tn *TrustNode) getPubkeys() []cipher.PubKey {
	pubkeys := []cipher.PubKey{}
	bin := tn.node.Get([]byte("pubkey"))
	if bin == nil {
		return pubkeys
	}
	for _, pk := range strings.Split(string(bin), ",") {
		pubkeys = append(pubkeys, cipher.MustPubKeyFromHex(pk))
	}
	return pubkeys
}
