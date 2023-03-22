package pertrie

import "fmt"

type wtx struct {
	db *DB
}

var _ Trie = &wtx{}

func (w *wtx) Get(key []byte) (TrieOrData, error) {
	return nil, fmt.Errorf("not yet implemented")
}
func (w *wtx) Put(key, data []byte) error {
	return fmt.Errorf("not yet imlemented")
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
	tx := &wtx{db: d}
	return fn(tx)
}
