package blocks

import (
	"context"
	"fmt"

	"github.com/ledgerwatch/erigon/v2/cmd/devnet/devnet"
	"github.com/ledgerwatch/erigon/v2/cmd/devnet/requests"
	"github.com/ledgerwatch/erigon/v2/turbo/jsonrpc"
)

var CompletionChecker = BlockHandlerFunc(
	func(ctx context.Context, node devnet.Node, block *requests.Block, transaction *jsonrpc.RPCTransaction) error {
		traceResults, err := node.TraceTransaction(transaction.Hash)

		if err != nil {
			return fmt.Errorf("Failed to trace transaction: %s: %w", transaction.Hash, err)
		}

		for _, traceResult := range traceResults {
			if traceResult.TransactionHash == transaction.Hash {
				if len(traceResult.Error) != 0 {
					return fmt.Errorf("Transaction error: %s", traceResult.Error)
				}

				break
			}
		}

		return nil
	})
