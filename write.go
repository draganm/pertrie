package pertrie

import (
	"fmt"

	"capnproto.org/go/capnp/v3"
	"github.com/draganm/pertrie/trie"
)

type wtx struct {
	db        *DB
	rootBlock uint32
}

var _ Trie = &wtx{}

func (w *wtx) Get(key []byte) (TrieOrData, error) {
	return nil, fmt.Errorf("not yet implemented")
}
func (w *wtx) Put(key, data []byte) error {
	block, err := w.db.getBlock(w.rootBlock)
	if err != nil {
		return err
	}
	msg, err := capnp.Unmarshal(block)
	if err != nil {
		return fmt.Errorf("could not unmarshal root node: %w", err)
	}
	rootNode, err := trie.ReadRootNode(msg)
	if err != nil {
		return fmt.Errorf("could not read root node: %w", err)
	}

	return rootNode.Put(key, data)

}
func (w *wtx) Delete(key []byte) error {
	return fmt.Errorf("not yet imlemented")
}
func (w *wtx) NewTrie(key []byte) (Trie, error) {
	return nil, fmt.Errorf("not yet implemented")
}
func (w *wtx) Iterator() (Iterator, error) {
	return nil, fmt.Errorf("not yet imlemented")
}

func (w *wtx) Size() uint64 {
	return 0
}

func (d *DB) Write(fn func(t Trie) error) error {
	tx := &wtx{db: d, rootBlock: 1}
	return fn(tx)
}
