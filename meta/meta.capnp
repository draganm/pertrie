using Go = import "/go.capnp";
@0xc7a6b8fe322e2a04;
$Go.package("meta");
$Go.import("meta/meta");

struct FreeBlock {
    next @0 :UInt32;
}

struct Database {
    version @0 :UInt32;
    firstFreeBlock @1 :UInt32;
    nextUnallocatedBlock @2 :UInt32;
}

