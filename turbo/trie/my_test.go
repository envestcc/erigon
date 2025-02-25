package trie_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/holiman/uint256"
	"github.com/ledgerwatch/erigon-lib/chain"
	libcommon "github.com/ledgerwatch/erigon-lib/common"
	"github.com/ledgerwatch/erigon-lib/etl"
	"github.com/ledgerwatch/erigon-lib/kv"
	"github.com/ledgerwatch/erigon-lib/kv/mdbx"
	"github.com/ledgerwatch/erigon-lib/kv/temporal/historyv2"
	"github.com/ledgerwatch/erigon/core/state"
	"github.com/ledgerwatch/erigon/eth/stagedsync"
	"github.com/ledgerwatch/erigon/turbo/trie"
	"github.com/ledgerwatch/log/v3"
	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"
	"github.com/stretchr/testify/require"
)

const (
	dbDir     = "/tmp/mdbx"
	logPrefix = "mytrie_test"
)

var (
	contract  = libcommon.HexToAddress("0x71dd1027069078091B3ca48093B00E4735B20624")
	contract2 = libcommon.HexToAddress("0x71dd1027069078091B3ca48093B00E4735B20621")
	key       = libcommon.BigToHash(big.NewInt(1))
)

func TestMyTrieWrite(t *testing.T) {
	require.NoError(t, os.RemoveAll(dbDir))
	lg := log.New()
	lg.SetHandler(log.StdoutHandler)
	// log.Root().SetHandler(log.StdoutHandler)

	rw, err := mdbx.NewMDBX(lg).Path(dbDir).Open(context.Background())
	require.NoError(t, err)
	defer rw.Close()

	// history, err := temporal.New(rw, nil)

	// block 1
	lg.Info("block 1")
	func() {
		tx, err := rw.BeginRw(context.Background())
		require.NoError(t, err)
		defer func() { require.NoError(t, tx.Commit()) }()

		// _, tx := memdb.NewTestTx(t)
		r, tsw := state.NewDbStateReader(tx), state.NewDbStateWriter(tx, 1)
		// r, tsw := state.NewPlainStateReader(tx), state.NewPlainStateWriter(tx, tx, 1)
		intraBlockState := state.New(r)
		oldCode := []byte{0x01, 0x02, 0x03, 0x04}
		// Start the 1st transaction
		intraBlockState.CreateAccount(contract, true)
		intraBlockState.SetCode(contract, oldCode)
		intraBlockState.AddBalance(contract, uint256.NewInt(1000000000))
		intraBlockState.SetState(contract, &key, *uint256.NewInt(100))

		fmt.Println("finalizing 1st tx")
		if err := intraBlockState.FinalizeTx(&chain.Rules{}, tsw); err != nil {
			t.Errorf("error finalising 1st tx: %v", err)
		}
		fmt.Println("committing 1st tx")
		if err := intraBlockState.CommitBlock(&chain.Rules{}, tsw); err != nil {
			t.Errorf("error committing 1st tx: %v", err)
		}
		intraBlockState.Print(chain.Rules{})
		// root, err := trie.CalcRoot("", tx)
		// require.NoError(t, err)
		// t.Log("root", root)

		// acc, err := r.ReadAccountData(contract)
		// require.NoError(t, err)
		// t.Logf("acc: %+v", acc)
	}()

	// read account
	readFn := func() {
		tx, err := rw.BeginRo(context.Background())
		require.NoError(t, err)
		defer tx.Rollback()

		accTrieCollector := etl.NewCollector(logPrefix, dbDir, etl.NewSortableBuffer(etl.BufferOptimalSize), lg)
		defer accTrieCollector.Close()
		accTrieCollectorFunc := stagedsync.DebugAccountTrieCollector(accTrieCollector)
		stTrieCollector := etl.NewCollector(logPrefix, dbDir, etl.NewSortableBuffer(etl.BufferOptimalSize), lg)
		defer stTrieCollector.Close()
		stTrieCollectorFunc := stagedsync.DebugStorageTrieCollector(stTrieCollector)

		rl := trie.NewRetainList(0)
		r := state.NewDbStateReader(tx)
		// r := state.NewPlainStateReader(tx)
		a, err := r.ReadAccountData(contract)
		require.NoError(t, err)
		t.Logf("acc: %+v", a)
		pr, err := trie.NewProofRetainer(contract, a, nil, rl)
		require.NoError(t, err)
		loader := trie.NewFlatDBTrieLoader(logPrefix, rl, accTrieCollectorFunc, stTrieCollectorFunc, false)
		loader.SetProofRetainer(pr)
		root, err := loader.CalcTrieRoot(tx, nil)
		require.NoError(t, err)
		t.Log("trie root", root)
		proof, err := pr.ProofResult()
		require.NoError(t, err)
		t.Logf("account root: %+v", proof.StorageHash)
	}
	readFn()

	// block 2
	lg.Info("block 2")
	func() {
		tx, err := rw.BeginRw(context.Background())
		require.NoError(t, err)
		defer func() { require.NoError(t, tx.Commit()) }()

		r, tsw := state.NewDbStateReader(tx), state.NewDbStateWriter(tx, 2)
		// r, tsw := state.NewPlainStateReader(tx), state.NewPlainStateWriter(tx, tx, 2)
		intraBlockState := state.New(r)
		// Start the 2nd transaction
		intraBlockState.AddBalance(contract, uint256.NewInt(999))
		intraBlockState.SetState(contract, &key, *uint256.NewInt(300))
		intraBlockState.SetNonce(contract, 2)

		if err := intraBlockState.FinalizeTx(&chain.Rules{}, tsw); err != nil {
			t.Errorf("error finalising 2st tx: %v", err)
		}
		if err := intraBlockState.CommitBlock(&chain.Rules{}, tsw); err != nil {
			t.Errorf("error committing 1st tx: %v", err)
		}
	}()
	readFn()

	log.Info("walk db")
	func() {
		tx, err := rw.BeginRo(context.Background())
		require.NoError(t, err)
		defer tx.Rollback()

		tables, err := tx.ListBuckets()
		require.NoError(t, err)
		for _, table := range tables {
			err = tx.ForEach(table, nil, func(k, v []byte) error {
				t.Logf("table: %s, key: %s, value: %s", table, hex.EncodeToString(k), hex.EncodeToString(v))
				return nil
			})
			require.NoError(t, err)
		}
	}()

	return
	// block 3
	lg.Info("block 3")
	func() {
		tx, err := rw.BeginRw(context.Background())
		require.NoError(t, err)
		defer func() { require.NoError(t, tx.Commit()) }()

		// stageTrieCfg := stagedsync.StageTrieCfg(rw, false, false, false, dbDir, nil, nil, false, nil)
		// hash, err := stagedsync.RegenerateIntermediateHashes(logPrefix, tx, stageTrieCfg, libcommon.Hash{}, context.Background(), lg)
		// t.Log("hash", hash)

		r, tsw := state.NewDbStateReader(tx), state.NewDbStateWriter(tx, 3)
		// r, tsw := state.NewPlainStateReader(tx), state.NewPlainStateWriter(tx, tx, 3)
		intraBlockState := state.New(r)
		oldCode := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
		// Start the 2nd transaction
		intraBlockState.CreateAccount(contract2, true)
		intraBlockState.AddBalance(contract2, uint256.NewInt(1111000000))
		intraBlockState.SetCode(contract2, oldCode)
		intraBlockState.SetState(contract2, &key, *uint256.NewInt(200))
		intraBlockState.AddBalance(contract, uint256.NewInt(999))
		intraBlockState.SetState(contract, &key, *uint256.NewInt(300))

		if err := intraBlockState.FinalizeTx(&chain.Rules{}, tsw); err != nil {
			t.Errorf("error finalising 2st tx: %v", err)
		}
		if err := intraBlockState.CommitBlock(&chain.Rules{}, tsw); err != nil {
			t.Errorf("error committing 1st tx: %v", err)
		}
		intraBlockState.Print(chain.Rules{})
		root, err := trie.CalcRoot("test", tx)
		require.NoError(t, err)
		t.Log("root", root)

		acc, err := r.ReadAccountData(contract2)
		require.NoError(t, err)
		t.Logf("acc: %+v", acc)
	}()

}

func TestMyTrieRead(t *testing.T) {
	rw, err := mdbx.NewMDBX(log.New()).Path(dbDir).Open(context.Background())
	require.NoError(t, err)
	defer rw.Close()

	tx, err := rw.BeginRo(context.Background())
	require.NoError(t, err)
	defer tx.Rollback()

	// r := state.NewPlainStateReader(tx)
	r := state.NewDbStateReader(tx)
	a, err := r.ReadAccountData(contract)
	require.NoError(t, err)
	t.Logf("acc: %+v", a)
	a2, err := r.ReadAccountData(contract2)
	require.NoError(t, err)
	t.Logf("acc2: %+v", a2)

	rl := trie.NewRetainList(0)
	// retainKeys := []libcommon.Hash{key}
	// for _, retainKey := range retainKeys {
	// 	rl.AddKeyWithMarker(retainKey.Bytes(), true)
	// }
	lg := log.New()
	lg.SetHandler(log.StdoutHandler)
	accTrieCollector := etl.NewCollector(logPrefix, dbDir, etl.NewSortableBuffer(etl.BufferOptimalSize), lg)
	defer accTrieCollector.Close()
	accTrieCollectorFunc := stagedsync.DebugAccountTrieCollector(accTrieCollector)

	stTrieCollector := etl.NewCollector(logPrefix, dbDir, etl.NewSortableBuffer(etl.BufferOptimalSize), lg)
	defer stTrieCollector.Close()
	stTrieCollectorFunc := stagedsync.DebugStorageTrieCollector(stTrieCollector)

	loader := trie.NewFlatDBTrieLoader(logPrefix, rl, accTrieCollectorFunc, stTrieCollectorFunc, true)
	root, err := loader.CalcTrieRoot(tx, nil)
	require.NoError(t, err)
	t.Log("root", root)

	rl = trie.NewRetainList(0)
	pr, err := trie.NewProofRetainer(contract, a, nil, rl)
	require.NoError(t, err)
	loader.SetProofRetainer(pr)
	root, err = loader.CalcTrieRoot(tx, nil)
	require.NoError(t, err)
	t.Log("contract root", root)
	proof, err := pr.ProofResult()
	require.NoError(t, err)
	t.Logf("contract proof: %+v", proof)

	rl = trie.NewRetainList(0)
	pr, err = trie.NewProofRetainer(contract2, a2, nil, rl)
	require.NoError(t, err)
	loader.SetProofRetainer(pr)
	root, err = loader.CalcTrieRoot(tx, nil)
	require.NoError(t, err)
	t.Log("contract2 root", root)
	proof, err = pr.ProofResult()
	require.NoError(t, err)
	t.Logf("contract2 proof: %+v", proof)

	trie1, trie2, root := naiveTriesAndHashFromDB(t, rw)
	t.Log("native.root", root)
	t.Log("native.trie1.root", hex.EncodeToString(trie1.Root()))
	t.Log("native.trie2.root", hex.EncodeToString(trie2.Root()))
}

const dbpath = "/tmp/mdbx2"

func TestTrieStorage(t *testing.T) {
	// create database and do much changes, to see how many k v pairs are stored
	// test data: 10w account, 1000 storage per account

	// 1 - create db
	ctx := context.Background()
	require.NoError(t, os.RemoveAll(dbpath))
	lg := log.New()
	rw, err := mdbx.NewMDBX(lg).Path(dbpath).Open(ctx)
	require.NoError(t, err)
	defer rw.Close()
	t.Logf("create db at %s", dbpath)

	// prepare data
	codesTempl := make([]byte, 128)
	for i := 0; i < len(codesTempl); i++ {
		codesTempl[i] = byte(i)
	}

	// 2 - execute txs
	t.Log("execute txs")
	func() {
		tx, err := rw.BeginRw(ctx)
		require.NoError(t, err)
		defer func() { require.NoError(t, tx.Commit()) }()

		readers := []state.StateReader{state.NewDbStateReader(tx), state.NewPlainStateReader(tx)}
		writers := []state.StateWriter{state.NewDbStateWriter(tx, 1), state.NewPlainStateWriter(tx, tx, 1)}
		for i, r := range readers {
			tsw := writers[i]
			intraBlockState := state.New(r)

			accCnt := 10000
			bar := progressbar.New(accCnt)
			bar.RenderBlank()
			for i := 0; i < accCnt; i++ {
				bar.Add(1)
				contract := libcommon.BigToAddress(big.NewInt(int64(i)))
				intraBlockState.CreateAccount(contract, true)
				intraBlockState.SetBalance(contract, uint256.NewInt(uint64(i*i)))
				intraBlockState.SetNonce(contract, uint64(i))
				intraBlockState.SetCode(contract, codesTempl)
				for j := 0; j < 200; j++ {
					v := j
					key := libcommon.BigToHash(big.NewInt(int64(v)))
					intraBlockState.SetState(contract, &key, *uint256.NewInt(uint64(v)))
				}
			}

			t.Logf("finalize tx")
			err = intraBlockState.FinalizeTx(&chain.Rules{}, tsw)
			require.NoError(t, err)
			t.Logf("commit block")
			err = intraBlockState.CommitBlock(&chain.Rules{}, tsw)
			require.NoError(t, err)
			t.Logf("commit")
		}
	}()

	// 3 - walk db
	err = walkdb(rw)
	require.NoError(t, err)
}

func TestWalkDB(t *testing.T) {
	rw, err := mdbx.NewMDBX(log.New()).Path("/Users/chenchen/iotex-var-standalone/data/historyindex").Open(context.Background())
	require.NoError(t, err)
	defer rw.Close()

	err = walkdb(rw)
	require.NoError(t, err)
}

func walkdb(rw kv.RoDB) error {
	tx, err := rw.BeginRo(context.Background())
	if err != nil {
		return errors.Wrap(err, "begin ro")
	}
	defer tx.Rollback()

	dbsize, err := tx.DBSize()
	if err != nil {
		return errors.Wrap(err, "db size")
	}
	fmt.Printf("db size: %d\n", dbsize)
	tables, err := tx.ListBuckets()
	if err != nil {
		return errors.Wrap(err, "list tables")
	}
	// fmt.Printf("tables: %v\n", tables)
	total := uint64(0)
	for _, table := range tables {
		tsize, err := tx.BucketSize(table)
		if err != nil {
			return errors.Wrapf(err, "table size: %s", table)
		}
		if tsize == 0 {
			continue
		}
		total += tsize
		keynum := uint64(0)
		err = tx.ForEach(table, nil, func(k, v []byte) error {
			if keynum < 300 {
				fmt.Printf("table: %s, key: %x, value: %x\n", table, k, v)
			}
			keynum++
			return nil
		})
		if err != nil {
			return errors.Wrapf(err, "for each table: %s", table)
		}
		fmt.Printf("table: %s, size: %d, keynum: %d\n", table, tsize, keynum)
	}
	fmt.Printf("total: %d\n", total)
	return nil
}

const historydbpath = "/tmp/mdbx-history"

// TestHistoryTrieStorage is a test to see how it works with history trie storage
func TestHistoryTrieStorage(t *testing.T) {
	// 1 - create db
	ctx := context.Background()
	require.NoError(t, os.RemoveAll(historydbpath))
	lg := log.New()
	rw, err := mdbx.NewMDBX(lg).Path(historydbpath).Open(ctx)
	require.NoError(t, err)
	defer rw.Close()
	t.Logf("create history db at %s", historydbpath)

	// prepare data
	codesTempl := make([]byte, 128)
	for i := 0; i < len(codesTempl); i++ {
		codesTempl[i] = byte(i)
	}

	// 2 - execute txs
	t.Log("execute txs")
	func() {
		tx, err := rw.BeginRw(ctx)
		require.NoError(t, err)
		defer func() {
			t.Log("commit")
			require.NoError(t, tx.Commit())
		}()

		readers := []state.StateReader{state.NewPlainStateReader(tx)}
		writers := []state.WriterWithChangeSets{state.NewPlainStateWriter(tx, tx, 1)}
		for i, r := range readers {
			tsw := writers[i]
			intraBlockState := state.New(r)

			accCnt := 1
			for i := 0; i < accCnt; i++ {
				contract := libcommon.BigToAddress(big.NewInt(int64(i)))
				intraBlockState.CreateAccount(contract, true)
				intraBlockState.SetBalance(contract, uint256.NewInt(uint64(i*i)))
				intraBlockState.SetNonce(contract, uint64(i))
				intraBlockState.SetCode(contract, codesTempl)
				for j := 0; j < 1; j++ {
					v := j
					key := libcommon.BigToHash(big.NewInt(int64(v)))
					intraBlockState.SetState(contract, &key, *uint256.NewInt(uint64(v)))
				}
			}

			t.Logf("finalize tx")
			err = intraBlockState.FinalizeTx(&chain.Rules{}, tsw)
			require.NoError(t, err)
			t.Logf("commit block")
			err = intraBlockState.CommitBlock(&chain.Rules{}, tsw)
			require.NoError(t, err)
			t.Logf("write changesets")
			err = tsw.WriteChangeSets()
			require.NoError(t, err)
			t.Logf("write history")
			err = tsw.WriteHistory()
			require.NoError(t, err)
		}
	}()

	// 3 - walk db
	err = walkdbkv(rw)
	require.NoError(t, err)
}

func walkdbkv(rw kv.RoDB) error {
	tx, err := rw.BeginRo(context.Background())
	if err != nil {
		return errors.Wrap(err, "begin ro")
	}
	defer tx.Rollback()

	// dbsize, err := tx.DBSize()
	// if err != nil {
	// 	return errors.Wrap(err, "db size")
	// }
	// fmt.Printf("db size: %d\n", dbsize)
	tables, err := tx.ListBuckets()
	if err != nil {
		return errors.Wrap(err, "list tables")
	}
	// fmt.Printf("tables: %v\n", tables)
	total := uint64(0)
	for _, table := range tables {
		tsize, err := tx.BucketSize(table)
		if err != nil {
			return errors.Wrapf(err, "table size: %s", table)
		}
		if tsize == 0 {
			continue
		}
		total += tsize
		keynum := uint64(0)
		err = tx.ForEach(table, nil, func(k, v []byte) error {
			fmt.Printf("table: %s, key: %x, value: %x\n", table, k, v)
			keynum++
			return nil
		})
		if err != nil {
			return errors.Wrapf(err, "for each table: %s", table)
		}
		// fmt.Printf("table: %s, size: %d, keynum: %d\n", table, tsize, keynum)
	}
	// fmt.Printf("total: %d\n", total)
	return nil
}

func TestWalkHistory(t *testing.T) {
	rw, err := mdbx.NewMDBX(log.New()).Path("/Users/chenchen/iotex-var-standalone/data/historyindex").Open(context.Background())
	require.NoError(t, err)
	defer rw.Close()

	err = walkdbkv(rw)
	require.NoError(t, err)
}

func TestMyHistoryTrieWrite(t *testing.T) {
	require.NoError(t, os.RemoveAll(historydbpath))
	lg := log.New()
	lg.SetHandler(log.StdoutHandler)
	log.Root().SetHandler(log.StdoutHandler)

	rw, err := mdbx.NewMDBX(lg).Path(historydbpath).WithTableCfg(func(defaultBuckets kv.TableCfg) kv.TableCfg {
		defaultBuckets["erigonsystem"] = kv.TableCfgItem{}
		return defaultBuckets
	} ).Open(context.Background())
	require.NoError(t, err)
	defer rw.Close()

	// block 1
	lg.Info("block 1")
	func() {
		tx, err := rw.BeginRw(context.Background())
		require.NoError(t, err)
		defer func() { require.NoError(t, tx.Commit()) }()

		r, tsw := state.NewPlainStateReader(tx), state.NewPlainStateWriter(tx, tx, 1)
		intraBlockState := state.New(r)
		oldCode := []byte{0x01, 0x02, 0x03, 0x04}
		// Start the 1st transaction
		intraBlockState.CreateAccount(contract, true)
		intraBlockState.SetCode(contract, oldCode)
		intraBlockState.AddBalance(contract, uint256.NewInt(1000000000))
		intraBlockState.SetState(contract, &key, *uint256.NewInt(100))

		fmt.Println("finalizing 1st tx")
		if err := intraBlockState.FinalizeTx(&chain.Rules{}, tsw); err != nil {
			t.Errorf("error finalising 1st tx: %v", err)
		}

		// start the 2nd transaction
		intraBlockState.SetState(contract, &key, *uint256.NewInt(109))
		fmt.Println("finalizing 2st tx")
		if err := intraBlockState.FinalizeTx(&chain.Rules{}, tsw); err != nil {
			t.Errorf("error finalising 2st tx: %v", err)
		}

		// start the 3rd transaction
		intraBlockState.SetState(contract, &key, *uint256.NewInt(100))
		fmt.Println("finalizing 3st tx")
		if err := intraBlockState.FinalizeTx(&chain.Rules{}, tsw); err != nil {
			t.Errorf("error finalising 3st tx: %v", err)
		}

		fmt.Println("committing 1st block")
		if err := intraBlockState.CommitBlock(&chain.Rules{}, tsw); err != nil {
			t.Errorf("error committing 1st block: %v", err)
		}
		intraBlockState.Print(chain.Rules{})

		t.Logf("write changesets")
		err = tsw.WriteChangeSets()
		require.NoError(t, err)
		t.Logf("write history")
		err = tsw.WriteHistory()
		require.NoError(t, err)
	}()

	// read account
	readFn := func() {
		tx, err := rw.BeginRo(context.Background())
		require.NoError(t, err)
		defer tx.Rollback()

		// accTrieCollector := etl.NewCollector(logPrefix, dbDir, etl.NewSortableBuffer(etl.BufferOptimalSize), lg)
		// defer accTrieCollector.Close()
		// accTrieCollectorFunc := stagedsync.DebugAccountTrieCollector(accTrieCollector)
		// stTrieCollector := etl.NewCollector(logPrefix, dbDir, etl.NewSortableBuffer(etl.BufferOptimalSize), lg)
		// defer stTrieCollector.Close()
		// stTrieCollectorFunc := stagedsync.DebugStorageTrieCollector(stTrieCollector)

		// rl := trie.NewRetainList(0)
		r := state.NewPlainState(tx, 2, nil)
		a, err := r.ReadAccountData(contract)
		require.NoError(t, err)
		t.Logf("acc: %+v", a)
		// pr, err := trie.NewProofRetainer(contract, a, nil, rl)
		// require.NoError(t, err)
		// loader := trie.NewFlatDBTrieLoader(logPrefix, rl, accTrieCollectorFunc, stTrieCollectorFunc, false)
		// loader.SetProofRetainer(pr)
		// root, err := loader.CalcTrieRoot(tx, nil)
		// require.NoError(t, err)
		// t.Log("trie root", root)
		// proof, err := pr.ProofResult()
		// require.NoError(t, err)
		// t.Logf("account root: %+v", proof.StorageHash)
	}
	readFn()

	// block 2
	lg.Info("block 2")
	func() {
		tx, err := rw.BeginRw(context.Background())
		require.NoError(t, err)
		defer func() { require.NoError(t, tx.Commit()) }()

		r, tsw := state.NewPlainStateReader(tx), state.NewPlainStateWriter(tx, tx, 2)
		intraBlockState := state.New(r)
		// Start the 1st transaction
		intraBlockState.AddBalance(contract, uint256.NewInt(1000000000))
		intraBlockState.SetState(contract, &key, *uint256.NewInt(200))

		fmt.Println("finalizing 1st tx")
		if err := intraBlockState.FinalizeTx(&chain.Rules{}, tsw); err != nil {
			t.Errorf("error finalising 1st tx: %v", err)
		}
		fmt.Println("committing 1st tx")
		if err := intraBlockState.CommitBlock(&chain.Rules{}, tsw); err != nil {
			t.Errorf("error committing 1st tx: %v", err)
		}
		intraBlockState.Print(chain.Rules{})

		t.Logf("write changesets")
		err = tsw.WriteChangeSets()
		require.NoError(t, err)
		t.Logf("write history")
		err = tsw.WriteHistory()
		require.NoError(t, err)
	}()

	func() {
		tx, err := rw.BeginRo(context.Background())
		require.NoError(t, err)
		defer tx.Rollback()

		for i := 0; i <= 3; i++ {
			r := state.NewPlainState(tx, uint64(i), nil)
			a, err := r.ReadAccountData(contract)
			require.NoError(t, err)
			v, err := r.ReadAccountStorage(contract, 1, &key)
			require.NoError(t, err)
			
			ir := state.New(r)
			val := uint256.NewInt(0)
			ir.GetState(contract, &key, val)

			t.Logf("acc at %d: %+v, value: %x, value2: %x", i, a, v, val)
		}
	}()
}

type testKVTx struct {
	kv.RwTx
}

func (t *testKVTx) Tx() kv.Tx {
	return t.RwTx
}

func TestMyHistoryDBTrieWrite(t *testing.T) {
	require.NoError(t, os.RemoveAll(historydbpath))
	lg := log.New()
	lg.SetHandler(log.StdoutHandler)
	// log.Root().SetHandler(log.StdoutHandler)

	rw, err := mdbx.NewMDBX(lg).Path(historydbpath).Open(context.Background())
	require.NoError(t, err)

	// agg, err := libstate.NewAggregator(context.Background(), "/tmp/snaphistory", "/tmp/tmp", config3.HistoryV3AggregationStep, rw, log.New())
	// require.NoError(t, err)
	// err = agg.OpenFolder()
	// require.NoError(t, err)
	// rw, err = temporal.New(rw, agg)
	// require.NoError(t, err)
	// defer rw.Close()

	// block 1
	lg.Info("block 1")
	func() {
		tx, err := rw.BeginRw(context.Background())
		require.NoError(t, err)
		defer func() { require.NoError(t, tx.Commit()) }()

		r, tsw := state.NewDbStateReader(tx), state.NewDbStateWriter(&testKVTx{tx}, 1)
		intraBlockState := state.New(r)
		oldCode := []byte{0x01, 0x02, 0x03, 0x04}
		// Start the 1st transaction
		intraBlockState.CreateAccount(contract, true)
		intraBlockState.SetCode(contract, oldCode)
		intraBlockState.AddBalance(contract, uint256.NewInt(1000000000))
		intraBlockState.SetState(contract, &key, *uint256.NewInt(100))

		fmt.Println("finalizing 1st tx")
		if err := intraBlockState.FinalizeTx(&chain.Rules{}, tsw); err != nil {
			t.Errorf("error finalising 1st tx: %v", err)
		}
		fmt.Println("committing 1st tx")
		if err := intraBlockState.CommitBlock(&chain.Rules{}, tsw); err != nil {
			t.Errorf("error committing 1st tx: %v", err)
		}
		intraBlockState.Print(chain.Rules{})

		t.Logf("write changesets")
		err = tsw.WriteChangeSets()
		require.NoError(t, err)
		t.Logf("write history")
		err = tsw.WriteHistory()
		require.NoError(t, err)
	}()

	// read account
	readFn := func() {
		tx, err := rw.BeginRo(context.Background())
		require.NoError(t, err)
		defer tx.Rollback()

		accTrieCollector := etl.NewCollector(logPrefix, dbDir, etl.NewSortableBuffer(etl.BufferOptimalSize), lg)
		defer accTrieCollector.Close()
		accTrieCollectorFunc := stagedsync.DebugAccountTrieCollector(accTrieCollector)
		stTrieCollector := etl.NewCollector(logPrefix, dbDir, etl.NewSortableBuffer(etl.BufferOptimalSize), lg)
		defer stTrieCollector.Close()
		stTrieCollectorFunc := stagedsync.DebugStorageTrieCollector(stTrieCollector)

		rl := trie.NewRetainList(0)
		r := state.NewPlainState(tx, 2, nil)
		c, err := tx.CursorDupSort(kv.E2AccountsHistory)
		require.NoError(t, err)
		historyv2.FindAccount(c, 2, contract[:])
		a, err := r.ReadAccountData(contract)
		require.NoError(t, err)
		t.Logf("acc: %+v", a)
		pr, err := trie.NewProofRetainer(contract, a, nil, rl)
		require.NoError(t, err)
		loader := trie.NewFlatDBTrieLoader(logPrefix, rl, accTrieCollectorFunc, stTrieCollectorFunc, false)
		loader.SetProofRetainer(pr)
		root, err := loader.CalcTrieRoot(tx, nil)
		require.NoError(t, err)
		t.Log("trie root", root)
		// proof, err := pr.ProofResult()
		// require.NoError(t, err)
		// t.Logf("account root: %+v", proof.StorageHash)
	}
	readFn()

	// block 2
	lg.Info("block 2")
	func() {
		tx, err := rw.BeginRw(context.Background())
		require.NoError(t, err)
		defer func() { require.NoError(t, tx.Commit()) }()

		r, tsw := state.NewDbStateReader(tx), state.NewDbStateWriter(&testKVTx{tx}, 2)
		intraBlockState := state.New(r)
		// Start the 1st transaction
		intraBlockState.AddBalance(contract, uint256.NewInt(1000000000))
		intraBlockState.SetState(contract, &key, *uint256.NewInt(200))

		fmt.Println("finalizing 1st tx")
		if err := intraBlockState.FinalizeTx(&chain.Rules{}, tsw); err != nil {
			t.Errorf("error finalising 1st tx: %v", err)
		}
		fmt.Println("committing 1st tx")
		if err := intraBlockState.CommitBlock(&chain.Rules{}, tsw); err != nil {
			t.Errorf("error committing 1st tx: %v", err)
		}
		intraBlockState.Print(chain.Rules{})

		t.Logf("write changesets")
		err = tsw.WriteChangeSets()
		require.NoError(t, err)
		t.Logf("write history")
		err = tsw.WriteHistory()
		require.NoError(t, err)
	}()

	func() {
		tx, err := rw.BeginRo(context.Background())
		require.NoError(t, err)
		defer tx.Rollback()

		for i := 0; i <= 10; i++ {
			r := state.NewHistoryReaderV3()
			r.SetTx(tx)
			r.SetTxNum(uint64(i))
			// r := state.NewPlainState(tx, uint64(i), nil)
			a, err := r.ReadAccountData(contract)
			require.NoError(t, err)
			v, err := r.ReadAccountStorage(contract, 1, &key)
			require.NoError(t, err)
			t.Logf("acc at %d: %+v, value: %x", i, a, v)
		}
	}()
}
