package machine

import (
	"github.com/ledgerwatch/erigon-lib/common"
	"github.com/ledgerwatch/erigon/v2/cl/abstract"
	"github.com/ledgerwatch/erigon/v2/cl/cltypes"
	"github.com/ledgerwatch/erigon/v2/cl/phase1/core/state"
)

func executionEnabled(s abstract.BeaconState, payload *cltypes.Eth1Block) bool {
	return (!state.IsMergeTransitionComplete(s) && payload.BlockHash != common.Hash{}) || state.IsMergeTransitionComplete(s)
}
