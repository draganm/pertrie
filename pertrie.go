package pertrie

type Trie interface {
	Get(key []byte) (TrieOrData, error)
	Put(key, data []byte) error
	Delete(key []byte) error
	NewTrie(key []byte) (Trie, error)
	Iterator() (Iterator, error)
}

type Iterator interface {
	HasNext() bool
	GetKey() []byte
	GetValue() TrieOrData
	Next()
	Prev()
	Seek(key []byte)
}

type TrieOrData interface {
	IsTrie() bool
	GetData() []byte
	GetTrie()
}
