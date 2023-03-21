package trie_test

import (
	"testing"

	"capnproto.org/go/capnp/v3"
	"github.com/draganm/pertrie/trie"
	"github.com/stretchr/testify/require"
)

func TestNode(t *testing.T) {
	require := require.New(t)

	var d []byte
	{
		arena := capnp.SingleSegment(nil)

		msg, seg, err := capnp.NewMessage(arena)
		require.NoError(err)

		node, err := trie.NewRootNode(seg)
		require.NoError(err)

		err = node.SetPrefix([]byte("foo"))
		require.NoError(err)

		ncl, err := trie.NewNode_Child_List(seg, 16)
		require.NoError(err)
		node.SetChildren(ncl)

		for i := 0; i < 16; i++ {
			nc, err := trie.NewNode_Child(seg)
			require.NoError(err)
			val := nc.Value()
			val.SetBlockRef(uint32(333 + i))
			ncl.Set(i, nc)
		}

		nv := node.Value()
		nv.SetBlockRef(666)

		d, err = msg.Marshal()
		require.NoError(err)
	}
	require.Equal(576, len(d))

	msg, err := capnp.Unmarshal(d)
	require.NoError(err)

	node, err := trie.ReadRootNode(msg)
	require.NoError(err)
	pref, err := node.Prefix()
	require.NoError(err)
	require.Equal(string(pref), "foo")

}
