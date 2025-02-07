// Copyright 2022 The Erigon Authors
// This file is part of Erigon.
//
// Erigon is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Erigon is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Erigon. If not, see <http://www.gnu.org/licenses/>.

package bptree

import (
	"fmt"
)

type Stats struct {
	ExposedCount  uint
	RehashedCount uint
	CreatedCount  uint
	UpdatedCount  uint
	DeletedCount  uint
	OpeningHashes uint
	ClosingHashes uint
}

type Tree23 struct {
	root *Node23
}

func NewEmptyTree23() *Tree23 {
	return &Tree23{}
}

func NewTree23(kvItems KeyValues) *Tree23 {
	tree := new(Tree23).Upsert(kvItems)
	tree.reset()
	return tree
}

func (t *Tree23) String() string {
	return fmt.Sprintf("root={keys=%v #children=%d} size=%d", deref(t.root.keys), t.root.childrenCount(), t.Size())
}

func (t *Tree23) Size() int {
	count := 0
	t.WalkPostOrder(func(n *Node23) interface{} { count++; return nil })
	return count
}

func (t *Tree23) RootHash() []byte {
	if t.root == nil {
		return []byte{}
	}
	return t.root.hashNode()
}

func (t *Tree23) IsValid() (bool, error) {
	if t.root == nil {
		return true, nil
	}
	// Last leaf must have sentinel next key
	if lastLeaf := t.root.lastLeaf(); lastLeaf.keyCount() > 0 && lastLeaf.nextKey() != nil {
		return false, fmt.Errorf("no sentinel next key in last leaf %d", &lastLeaf)
	}
	return t.root.isValid()
}

func (t *Tree23) Graph(filename string, debug bool) {
	graph := NewGraph(t.root)
	graph.saveDot(filename, debug)
}

func (t *Tree23) GraphAndPicture(filename string) error {
	graph := NewGraph(t.root)
	return graph.saveDotAndPicture(filename, false)
}

func (t *Tree23) GraphAndPictureDebug(filename string) error {
	graph := NewGraph(t.root)
	return graph.saveDotAndPicture(filename, true)
}

func (t *Tree23) Height() int {
	if t.root == nil {
		return 0
	}
	return t.root.height()
}

func (t *Tree23) KeysInLevelOrder() []Felt {
	if t.root == nil {
		return []Felt{}
	}
	return t.root.keysInLevelOrder()
}

func (t *Tree23) WalkPostOrder(w Walker) []interface{} {
	if t.root == nil {
		return make([]interface{}, 0)
	}
	return t.root.walkPostOrder(w)
}

func (t *Tree23) WalkKeysPostOrder() []Felt {
	keyPointers := make([]*Felt, 0)
	t.WalkPostOrder(func(n *Node23) interface{} {
		if n.isLeaf && n.keyCount() > 0 {
			keyPointers = append(keyPointers, n.keys[:len(n.keys)-1]...)
		}
		return nil
	})
	keys := deref(keyPointers)
	return keys
}

func (t *Tree23) Upsert(kvItems KeyValues) *Tree23 {
	return t.UpsertWithStats(kvItems, &Stats{})
}

func (t *Tree23) UpsertWithStats(kvItems KeyValues, stats *Stats) *Tree23 {
	promoted, _, intermediateKeys := upsert(t.root, kvItems, stats)
	ensure(len(promoted) > 0, "nodes length is zero")
	if len(promoted) == 1 {
		t.root = promoted[0]
	} else {
		t.root = promote(promoted, intermediateKeys, stats)
	}
	stats.RehashedCount, stats.ClosingHashes = t.countUpsertRehashedNodes()
	return t
}

func (t *Tree23) Delete(keyToDelete []Felt) *Tree23 {
	return t.DeleteWithStats(keyToDelete, &Stats{})
}

func (t *Tree23) DeleteWithStats(keysToDelete []Felt, stats *Stats) *Tree23 {
	newRoot, nextKey, intermediateKeys := del(t.root, keysToDelete, stats)
	t.root, _ = demote(newRoot, nextKey, intermediateKeys, stats)
	stats.RehashedCount, stats.ClosingHashes = t.countDeleteRehashedNodes()
	return t
}

func (t *Tree23) countUpsertRehashedNodes() (rehashedCount uint, closingHashes uint) {
	t.WalkPostOrder(func(n *Node23) interface{} {
		if n.exposed {
			rehashedCount++
			closingHashes += n.howManyHashes()
		}
		return nil
	})
	return rehashedCount, closingHashes
}

func (t *Tree23) countDeleteRehashedNodes() (rehashedCount uint, closingHashes uint) {
	t.WalkPostOrder(func(n *Node23) interface{} {
		if n.updated {
			rehashedCount++
			closingHashes += n.howManyHashes()
		}
		return nil
	})
	return rehashedCount, closingHashes
}

func (t *Tree23) reset() {
	if t.root == nil {
		return
	}
	t.root.reset()
}
