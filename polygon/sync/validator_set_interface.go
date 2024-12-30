package sync

import (
	libcommon "github.com/ledgerwatch/erigon-lib/common"
	"github.com/ledgerwatch/erigon/v2/polygon/bor"
)

// valset.ValidatorSet abstraction for unit tests
type validatorSetInterface interface {
	bor.ValidateHeaderTimeSignerSuccessionNumber
	IncrementProposerPriority(times int)
	Difficulty(signer libcommon.Address) (uint64, error)
}
