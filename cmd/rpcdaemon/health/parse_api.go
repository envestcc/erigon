package health

import (
	"github.com/ledgerwatch/erigon/v2/rpc"
)

func parseAPI(api []rpc.API) (netAPI NetAPI, ethAPI EthAPI) {
	for _, rpc := range api {
		if rpc.Service == nil {
			continue
		}

		if netCandidate, ok := rpc.Service.(NetAPI); ok {
			netAPI = netCandidate
		}

		if ethCandidate, ok := rpc.Service.(EthAPI); ok {
			ethAPI = ethCandidate
		}
	}
	return netAPI, ethAPI
}
