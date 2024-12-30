package requests

import (
	"context"

	"github.com/ledgerwatch/erigon/v2/p2p"
)

func (reqGen *requestGenerator) AdminNodeInfo() (p2p.NodeInfo, error) {
	var result p2p.NodeInfo

	if err := reqGen.rpcCall(context.Background(), &result, Methods.AdminNodeInfo); err != nil {
		return p2p.NodeInfo{}, err
	}

	return result, nil
}
