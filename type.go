package bert

import (
	"io"
	"io/ioutil"
)

const (
	VersionTag            = 131
	DistributionHeaderTag = 68
	CompressedTag         = 80
	SmallIntTag           = 97
	IntTag                = 98
	SmallBignumTag        = 110
	LargeBignumTag        = 111
	FloatTag              = 99
	NewFloatTag           = 70
	AtomCacheRefTag       = 82
	AtomTag               = 100
	SmallAtomTag          = 115
	AtomUtf8Tag           = 118
	SmallAtomUtf8Tag      = 119
	SmallTupleTag         = 104
	LargeTupleTag         = 105
	NilTag                = 106
	StringTag             = 107
	ListTag               = 108
	BinTag                = 109
	MapTag                = 116
	PidTag                = 103
	PortTag               = 102
	FunTag                = 117
	ReferenceTag          = 101
	NewReferenceTag       = 114
	NewFunTag             = 112
	ExportTag             = 113
)

type Atom string

const (
	BertAtom  = Atom("bert")
	NilAtom   = Atom("nil")
	TrueAtom  = Atom("true")
	FalseAtom = Atom("false")
)

const (
	MinorVersion0 = 0
	MinorVersion1 = 1
)

type Term interface{}

type Request struct {
	Kind      Atom
	Module    Atom
	Function  Atom
	Arguments []Term
}

const (
	// NewCacheEntry is used to determine if an entry is a new cache entry
	NewCacheEntry byte = 8

	// SegmentIndex can be used to extract the segment index
	SegmentIndex byte = 7

	// LongAtoms is used to determine if 2 byte atoms are used
	LongAtoms byte = 1
)

// As of erts version 5.7.2 the old atom cache protocol was dropped and a new one was introduced.
// This atom cache protocol introduced the distribution header. Nodes with erts versions earlier than
// 5.7.2 can still communicate with new nodes, but no distribution header and no atom cache will be used.
//
// The distribution header currently only contains an atom cache reference section, but could in the future
// contain more information. The distribution header precedes one or more Erlang terms on the external format.
// For more information see the documentation of the protocol between connected nodes in the distribution
// protocol documentation.
//
// ATOM_CACHE_REF entries with corresponding AtomCacheReferenceIndex in terms encoded on the external format
// following a distribution header refers to the atom cache references made in the distribution header.
// The range is 0 <= AtomCacheReferenceIndex < 255, i.e., at most 255 different atom cache references
// from the following terms can be made.
type DistributionHeader struct {

	// bucket holds all the available atoms that have been set
	// bucket[AtomCacheReferenceIndex][SegmentIndex][InternalSegmentIndex]
	bucket [255][8][256]*Atom

	// cache holds the current lookup into the bucket for the specific atom cache reference
	// cache[AtomCacheReferenceIndex]
	cache [255]*cacheIndex

	// Flags for an even AtomCacheReferenceIndex are located in the least significant half byte and flags for an
	// odd AtomCacheReferenceIndex are located in the most significant half byte.
	//
	// The flag field of an atom cache reference has the following format:
	// 1 bit	        3 bits
	// NewCacheEntryFlag	SegmentIndex
	//
	// The most significant bit is the NewCacheEntryFlag. If set, the corresponding cache reference is new.
	// The three least significant bits are the SegmentIndex of the corresponding atom cache entry.
	// An atom cache consists of 8 segments each of size 256, i.e., an atom cache can contain 2048 entries.
	flags []byte
}

// cacheIndex holds the current cache position for an atom
type cacheIndex struct {

	// SegmentIndex of the current atom cache
	segmentIndex uint8

	// InternalIndex of the current atom cache
	internalIndex uint8
}

// GetAtom will return the atom that exists for the atomCacheReferenceIndex
func (dh DistributionHeader) GetAtom(atomCacheReferenceIndex uint8) (*Atom, error) {

	// Get the atom cache index position
	atomCache := dh.cache[atomCacheReferenceIndex]

	// Look up the segment and internal index from the latest flags
	if atomCache != nil {
		atom := dh.bucket[atomCacheReferenceIndex][atomCache.segmentIndex][atomCache.internalIndex]
		if atom == nil {
			return nil, ErrMissingAtom
		}
		return atom, nil
	}
	return nil, ErrMissingAtom
}

// UpdateFlags will update the flags for the cache
func (dh DistributionHeader) Update(r io.Reader) error {
	noAtomRefs, err := read1(r)
	if err != nil {
		return err
	}

	// If NumberOfAtomCacheRefs is 0, Flags and AtomCacheRefs are omitted
	if noAtomRefs != 0 {

		// Flags consists of NumberOfAtomCacheRefs/2+1 bytes
		flags, err := ioutil.ReadAll(io.LimitReader(r, int64((noAtomRefs/2)+1)))
		if err != nil {
			return err
		}

		// Are these long atoms? Check the last half byte least significant bit
		atomLen := 1
		if dh.flags[len(dh.flags)-1]&LongAtoms == LongAtoms {
			atomLen = 2
		}

		// The flag information is stored within the
		for i, even := 0, true; i < noAtomRefs; i, even = i+1, !even {

			// Get the cache item
			cacheItem := dh.cache[i]
			if cacheItem == nil {
				cacheItem = &cacheIndex{}
				dh.cache[i] = cacheItem
			}

			// We need the cache entry and segment index
			newCacheEntry := (flags[i/2] & NewCacheEntry) == NewCacheEntry
			cacheItem.segmentIndex = flags[i/2] & SegmentIndex
			if !even {
				newCacheEntry = (flags[i/2] >> 4 & NewCacheEntry) == NewCacheEntry
				cacheItem.segmentIndex = flags[i/2] >> 4 & SegmentIndex
			}

			// We have the information to extract this atom
			internalSegmentIndex, err := read1(r)
			if err != nil {
				return err
			}
			cacheItem.internalIndex = uint8(internalSegmentIndex)

			// Extract the atom info for this is a new entry
			if newCacheEntry {

				// The length of the atom (can be 1 or 2 bytes)
				alen := 0
				if atomLen == 2 {
					alen, err = read2(r)
				} else {
					alen, err = read1(r)
				}
				if err != nil {
					return err
				}

				// Get the atom
				atomBytes, err := ioutil.ReadAll(io.LimitReader(r, int64(alen)))
				if err != nil {
					return err
				}

				// Store the atom in the bucket using the index position
				atom := Atom(atomBytes)
				dh.bucket[i][cacheItem.segmentIndex][cacheItem.internalIndex] = &atom
			}
		}
	}
	return nil
}

// Reference wraps the REFERENCE_EXT tag type (101)
//
// Encode a reference object (an object generated with make_ref/0).
//
// The Node term is an encoded atom, i.e. ATOM_EXT, SMALL_ATOM_EXT or ATOM_CACHE_REF.
//
// The ID field contains a big-endian unsigned integer, but should be regarded as uninterpreted data
// since this field is node specific. Creation is a byte containing a node serial number that makes it
// possible to separate old (crashed) nodes from a new one.
//
// In ID, only 18 bits are significant; the rest should be 0. In Creation, only 2 bits are significant;
// the rest should be 0. See NEW_REFERENCE_EXT.
type Reference struct {
	Node     Atom
	ID       uint32
	Creation uint8
}

// NewReference wraps the NEW_REFERENCE_EXT tag type (114)
//
// Node and Creation are as in REFERENCE_EXT.
//
// ID contains a sequence of big-endian unsigned integers (4 bytes each, so N' is a multiple of 4),
// but should be regarded as uninterpreted data.
//
// N' = 4 * Len.
//
// In the first word (four bytes) of ID, only 18 bits are significant, the rest should be 0.
//
// In Creation, only 2 bits are significant, the rest should be 0.
//
// NEW_REFERENCE_EXT was introduced with distribution version 4. In version 4, N' should be at most 12
type NewReference struct {
	Node     Atom
	Creation uint8
	ID       []uint32
}

// Port wraps the PORT_EXT tag type (102)
//
// Encode a port object (obtained form open_port/2). The ID is a node specific identifier for a local port.
// Port operations are not allowed across node boundaries. The Creation works just like in REFERENCE_EXT.
type Port struct {
	Node     Atom
	ID       uint32
	Creation uint8
}

// Pid wraps the PID_EXT tag type (103)
//
// Encode a process identifier object (obtained from spawn/3 or friends). The ID and Creation fields
// works just like in REFERENCE_EXT, while the Serial field is used to improve safety. In ID,
// only 15 bits are significant; the rest should be 0.
type Pid struct {
	Node     Atom
	ID       uint32
	Serial   uint32
	Creation uint8
}

// Func wraps the FUN_EXT tag type (117)
//
// Pid
// is a process identifier as in PID_EXT. It represents the process in which the fun was created.
//
// Module
// is an encoded as an atom, using ATOM_EXT, SMALL_ATOM_EXT or ATOM_CACHE_REF. This is the module that the fun is implemented in.
//
// Index
// is an integer encoded using SMALL_INTEGER_EXT or INTEGER_EXT. It is typically a small index into the module's fun table.
//
// Uniq
// is an integer encoded using SMALL_INTEGER_EXT or INTEGER_EXT. Uniq is the hash value of the parse for the fun.
//
// Free vars
// is NumFree number of terms, each one encoded according to its type.
type Func struct {
	Pid      Pid
	Module   Atom
	Index    uint32
	Uniq     uint32
	FreeVars []Term
}

// NewFunc wraps the NEW_FUN_EXT tag type (112)
// This is the new encoding of internal funs: fun F/A and fun(Arg1,..) -> ... end.
//
// Size
// is the total number of bytes, including the Size field.
//
// Arity
// is the arity of the function implementing the fun.
//
// Uniq
// is the 16 bytes MD5 of the significant parts of the Beam file.
//
// Index
// is an index number. Each fun within a module has an unique index. Index is stored in big-endian byte order.
//
// NumFree
// is the number of free variables.
//
// Module
// is an encoded as an atom, using ATOM_EXT, SMALL_ATOM_EXT or ATOM_CACHE_REF. This is the module that the fun is implemented in.
//
// OldIndex
// is an integer encoded using SMALL_INTEGER_EXT or INTEGER_EXT. It is typically a small index into the module's fun table.
//
// OldUniq
// is an integer encoded using SMALL_INTEGER_EXT or INTEGER_EXT. Uniq is the hash value of the parse tree for the fun.
//
// Pid
// is a process identifier as in PID_EXT. It represents the process in which the fun was created.
//
// Free vars
// is NumFree number of terms, each one encoded according to its type.
type NewFunc struct {
	Arity     uint8
	Uniq      []byte
	Index     uint32
	Module    Atom
	OldIndex  uint32
	OldUnique uint32
	Pid       Pid
	FreeVars  []Term
}

// Export wraps the EXPORT_EXT tag type (113)
// This term is the encoding for external funs: fun M:F/A.
//
// Module and Function are atoms (encoded using ATOM_EXT, SMALL_ATOM_EXT or ATOM_CACHE_REF).
//
// Arity is an integer encoded using SMALL_INTEGER_EXT.
type Export struct {
	Module   Atom
	Function Atom
	Arity    uint8
}
