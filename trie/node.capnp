using Go = import "/go.capnp";
@0xa21abfe7a0a832e0;
$Go.package("trie");
$Go.import("trie/node");

struct Node {
    prefix @0 :Prefix;
    value :union {
        nil @1 :Void;
        content @2 :TrieRootOrValue;
        blockRef @3 :UInt32;
    }

    struct Prefix {
        data @0 :Data;
        lastByteIsHalfByte @1 :Bool;
    }


    children @4 :List(Child);

    struct Child {
        value :union {
            nil @0 :Void;
            embedded @1 :Node;
            blockRef @2 :UInt32;
        }
    }
}


struct TrieRoot {
    count @0 :UInt64;
    seq @1 :UInt64;
    node @2 :Node;
}

struct ValueSegment {
    data @0 :Data;
    nextSegment @1 :UInt32;
}

struct Value {
        size @0 :UInt64;
        firstSegment @1 :ValueSegment;
}

struct TrieRootOrValue {
    value :union {
        trieRoot @0 :TrieRoot;
        value @1 :Value;
    }  
}

