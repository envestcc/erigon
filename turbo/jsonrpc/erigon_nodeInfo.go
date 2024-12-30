package jsonrpc

import (
	"context"

	"github.com/ledgerwatch/erigon/v2/p2p"
)

const (
	// allNodesInfo used in NodeInfo request to receive meta data from all running sentries.
	allNodesInfo = 0
)

func (api *ErigonImpl) NodeInfo(ctx context.Context) ([]p2p.NodeInfo, error) {
	return api.ethBackend.NodeInfo(ctx, allNodesInfo)
}
