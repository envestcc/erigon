package handler

import (
	"github.com/ledgerwatch/erigon/v2/cl/beacon/beaconhttp"
)

func newBeaconResponse(data any) *beaconhttp.BeaconResponse {
	return beaconhttp.NewBeaconResponse(data)
}
