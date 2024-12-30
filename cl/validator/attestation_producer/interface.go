package attestation_producer

import (
	"github.com/ledgerwatch/erigon/v2/cl/cltypes/solid"
	"github.com/ledgerwatch/erigon/v2/cl/phase1/core/state"
)

type AttestationDataProducer interface {
	ProduceAndCacheAttestationData(baseState *state.CachingBeaconState, slot uint64, committeeIndex uint64) (solid.AttestationData, error)
}
