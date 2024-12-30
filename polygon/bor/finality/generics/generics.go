package generics

import (
	libcommon "github.com/ledgerwatch/erigon-lib/common"
	"github.com/ledgerwatch/erigon/v2/core/types"
)

func Empty[T any]() (t T) {
	return
}

type Response struct {
	Headers []*types.Header
	Hashes  []libcommon.Hash
}
