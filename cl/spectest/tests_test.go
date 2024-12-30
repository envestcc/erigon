package spectest

import (
	"os"
	"testing"

	"github.com/ledgerwatch/erigon/v2/spectest"

	"github.com/ledgerwatch/erigon/v2/cl/transition"

	"github.com/ledgerwatch/erigon/v2/cl/spectest/consensus_tests"
)

func Test(t *testing.T) {
	spectest.RunCases(t, consensus_tests.TestFormats, transition.ValidatingMachine, os.DirFS("./tests"))
}
