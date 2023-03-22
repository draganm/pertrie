package pertrie

import (
	"fmt"
	"os"
	"syscall"

	"capnproto.org/go/capnp/v3"
	"github.com/draganm/pertrie/meta"
	"github.com/draganm/pertrie/trie"
	"golang.org/x/sys/unix"
)

type DB struct {
	openFile *os.File
	mapped   []byte
}

func (d *DB) Close() error {
	err := unix.Munmap(d.mapped)
	if err != nil {
		return fmt.Errorf("failed to munmap: %w", err)
	}
	return d.openFile.Close()
}

var pageSize = os.Getpagesize()

func (d *DB) init() error {

	err := d.initRootDatabase()
	if err != nil {
		return err
	}

	err = d.initRootTrie()
	if err != nil {
		return err
	}

	return d.openFile.Sync()

}

func (d *DB) initRootDatabase() error {

	metaBlock := make([]byte, pageSize)

	arena := capnp.SingleSegment(nil)

	msg, seg, err := capnp.NewMessage(arena)
	if err != nil {
		return fmt.Errorf("could not initialize meta block arena: %w", err)
	}

	db, err := meta.NewRootDatabase(seg)
	if err != nil {
		return fmt.Errorf("could not create root database object: %w", err)
	}

	metaMessage, err := msg.Marshal()
	if err != nil {
		return fmt.Errorf("could not marshal root meta: %w", err)
	}

	copy(metaBlock, metaMessage)

	db.SetNextUnallocatedBlock(2)
	db.SetVersion(1)
	_, err = d.openFile.Write(metaBlock)
	if err != nil {
		return fmt.Errorf("could not write meta block: %w", err)
	}

	return nil

}

func (d *DB) initRootTrie() error {

	rootTrieBlock := make([]byte, pageSize)

	arena := capnp.SingleSegment(nil)

	msg, seg, err := capnp.NewMessage(arena)
	if err != nil {
		return fmt.Errorf("could not initialize meta block arena")
	}

	root, err := trie.NewTrieRoot(seg)
	if err != nil {
		return fmt.Errorf("could not create root database object: %w", err)
	}

	rootNode, err := root.NewNode()
	if err != nil {
		return fmt.Errorf("could not create root node: %w", err)
	}

	_, err = rootNode.NewChildren(16)
	if err != nil {
		return fmt.Errorf("could not create root children: %w", err)
	}

	rootTrieMessage, err := msg.Marshal()
	if err != nil {
		return fmt.Errorf("could not marshal root trie: %w", err)
	}

	copy(rootTrieBlock, rootTrieMessage)
	_, err = d.openFile.Write(rootTrieBlock)
	if err != nil {
		return fmt.Errorf("could not write root node")
	}

	return nil

}

func (d *DB) mmap() error {
	info, err := d.openFile.Stat()
	if err != nil {
		return fmt.Errorf("could not stat db file: %w", err)
	}

	b, err := unix.Mmap(int(d.openFile.Fd()), 0, int(info.Size()), syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return fmt.Errorf("mmap failed: %w", err)
	}

	// Advise the kernel that the mmap is accessed randomly.
	err = unix.Madvise(b, syscall.MADV_RANDOM)
	if err != nil && err != syscall.ENOSYS {
		// Ignore not implemented error in kernel because it still works.
		return fmt.Errorf("madvise: %s", err)
	}

	d.mapped = b
	return nil
}

func Open(name string, mode os.FileMode) (*DB, error) {

	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, mode)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}

	st, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("could not stat db file: %w", err)
	}

	db := &DB{
		openFile: f,
	}

	if st.Size() == 0 {
		err = db.init()
		if err != nil {
			return nil, fmt.Errorf("could not init db: %w", err)
		}
	}

	err = db.mmap()
	if err != nil {
		return nil, fmt.Errorf("unable to mmap file: %w", err)
	}

	return db, nil

}
