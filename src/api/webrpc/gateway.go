package webrpc

import (
	"github.com/samoslab/samos/src/cipher"
	"github.com/samoslab/samos/src/coin"
	"github.com/samoslab/samos/src/daemon"
	"github.com/samoslab/samos/src/visor"
	"github.com/samoslab/samos/src/visor/historydb"
)

//go:generate goautomock -template=testify Gatewayer

// Gatewayer provides interfaces for getting samos related info.
type Gatewayer interface {
	GetLastBlocks(num uint64) (*visor.ReadableBlocks, error)
	GetBlocks(start, end uint64) (*visor.ReadableBlocks, error)
	GetBlocksInDepth(vs []uint64) (*visor.ReadableBlocks, error)
	GetUnspentOutputs(filters ...daemon.OutputsFilter) (*visor.ReadableOutputSet, error)
	GetTransaction(txid cipher.SHA256) (*visor.Transaction, error)
	InjectBroadcastTransaction(tx coin.Transaction) error
	GetAddrUxOuts(addr []cipher.Address) ([]*historydb.UxOut, error)
	GetTimeNow() uint64
}
