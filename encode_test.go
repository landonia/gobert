package bert

import (
	"bytes"
	"math/big"
	"reflect"
	"testing"
)

func TestEncodeNoCompression(t *testing.T) {

	// Small Integer
	assertEncode(t, 1, []byte{131, 97, 1})
	assertEncode(t, 42, []byte{131, 97, 42})

	// Integer
	assertEncode(t, 257, []byte{131, 98, 0, 0, 1, 1})
	assertEncode(t, 1025, []byte{131, 98, 0, 0, 4, 1})
	assertEncode(t, -1, []byte{131, 98, 255, 255, 255, 255})
	assertEncode(t, -8, []byte{131, 98, 255, 255, 255, 248})
	assertEncode(t, 5000, []byte{131, 98, 0, 0, 19, 136})
	assertEncode(t, -5000, []byte{131, 98, 255, 255, 236, 120})

	// Small Bignum
	assertEncode(t, big.NewInt(987654321), []byte{131, 110, 4, 0, 177, 104, 222, 58})
	assertEncode(t, big.NewInt(-987654321), []byte{131, 110, 4, 1, 177, 104, 222, 58})

	// Large Bignum
	bigNum := big.Int{}
	bigNum.SetBytes([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	assertEncode(t, bigNum, []byte{131, 111, 0, 0, 1, 44, 0, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1})

	// Put the number as negative
	bigNum.Neg(&bigNum)
	assertEncode(t, bigNum, []byte{131, 111, 0, 0, 1, 44, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1})

	// Old Float
	assertEncodeUsingMinor(t, 0.5, []byte{131, 99, 53, 46, 48, 48, 48, 48, 48, 48,
		48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 101,
		45, 48, 49, 0, 0, 0, 0, 0,
	}, false, MinorVersion0)
	assertEncodeUsingMinor(t, 3.14159, []byte{131, 99, 51, 46, 49, 52, 49, 53, 57,
		48, 49, 49, 56, 52, 48, 56, 50, 48, 51, 49, 50, 53, 48, 48,
		101, 43, 48, 48, 0, 0, 0, 0, 0,
	}, false, MinorVersion0)
	assertEncodeUsingMinor(t, -3.14159, []byte{131, 99, 45, 51, 46, 49, 52, 49, 53,
		57, 48, 49, 49, 56, 52, 48, 56, 50, 48, 51, 49, 50, 53, 48,
		48, 101, 43, 48, 48, 0, 0, 0, 0,
	}, false, MinorVersion0)

	// NewFloat
	assertEncode(t, 0.5,
		[]byte{131, 70, 63, 224, 0, 0, 0, 0, 0, 0})
	assertEncode(t, 3.14159,
		[]byte{131, 70, 64, 9, 33, 249, 240, 27, 134, 110})
	assertEncode(t, -3.14159,
		[]byte{131, 70, 192, 9, 33, 249, 240, 27, 134, 110})

	// Atom
	assertEncode(t, Atom("foo"),
		[]byte{131, 100, 0, 3, 102, 111, 111})

	// Atom UTF8
	assertEncode(t, Atom("foo, 世界"),
		[]byte{131, 118, 0, 11, 102, 111, 111, 44, 32, 228,
			184, 150, 231, 149, 140})

	// Small Tuple
	assertEncode(t, []Term{Atom("foo")},
		[]byte{131, 104, 1, 100, 0, 3, 102, 111, 111})
	assertEncode(t, []Term{Atom("foo"), Atom("bar")},
		[]byte{131, 104, 2,
			100, 0, 3, 102, 111, 111,
			100, 0, 3, 98, 97, 114,
		})
	assertEncode(t, []Term{Atom("coord"), 23, 42},
		[]byte{131, 104, 3,
			100, 0, 5, 99, 111, 111, 114, 100,
			97, 23,
			97, 42,
		})

	// Small Tuple with UTF8
	assertEncode(t, []Term{Atom("foo, 世界")},
		[]byte{131, 104, 1, 118, 0, 11, 102, 111, 111, 44, 32, 228,
			184, 150, 231, 149, 140})

	// Nil
	assertEncode(t, nil, []byte{131, 106})

	// String
	assertEncode(t, "foo", []byte{131, 107, 0, 3, 102, 111, 111})

	// List
	assertEncode(t, [1]Term{1},
		[]byte{131, 108, 0, 0, 0, 1, 97, 1, 106})
	assertEncode(t, [3]Term{1, 2, 3},
		[]byte{131, 108, 0, 0, 0, 3,
			97, 1, 97, 2, 97, 3,
			106,
		})
	assertEncode(t, [2]Term{Atom("a"), Atom("b")},
		[]byte{131, 108, 0, 0, 0, 2,
			100, 0, 1, 97, 100, 0, 1, 98,
			106,
		})

	// Map
	assertEncode(t, map[Term]Term{"key1": "value1", "key2": "value2", "key3": "value3"},
		[]byte{131, 116, 0, 0, 0, 3, 107, 0, 4, 107, 101, 121, 49, 107, 0, 6, 118,
			97, 108, 117, 101, 49, 107, 0, 4, 107, 101, 121, 50, 107, 0, 6, 118,
			97, 108, 117, 101, 50, 107, 0, 4, 107, 101, 121, 51, 107, 0, 6, 118,
			97, 108, 117, 101, 51,
		})

	// Pid
	pid := Pid{}
	pid.Node = Atom("Node")
	pid.ID = 123456789
	pid.Serial = 987654321
	pid.Creation = 128
	assertEncode(t, pid,
		[]byte{131, 103, 100, 0, 4, 78, 111,
			100, 101, 7, 91, 205, 21, 58,
			222, 104, 177, 128,
		})

	// Port
	port := Port{}
	port.Node = Atom("Node")
	port.ID = 123456789
	port.Creation = 128
	assertEncode(t, port,
		[]byte{131, 102, 100, 0, 4, 78, 111, 100,
			101, 7, 91, 205, 21, 128,
		})

	// Reference
	reference := Reference{}
	reference.Node = Atom("Node")
	reference.ID = 123456789
	reference.Creation = 128
	assertEncode(t, reference,
		[]byte{131, 101, 100, 0, 4, 78, 111, 100,
			101, 7, 91, 205, 21, 128,
		})

	// New Reference
	newReference := NewReference{}
	newReference.Creation = 128
	newReference.ID = []uint32{123, 234, 345, 456, 567, 678, 789, 890}
	newReference.Node = Atom("node")
	assertEncode(t, newReference,
		[]byte{131, 114, 0, 8, 100, 0, 4, 110, 111, 100,
			101, 128, 0, 0, 0, 123, 0, 0, 0, 234, 0,
			0, 1, 89, 0, 0, 1, 200, 0, 0, 2, 55, 0, 0,
			2, 166, 0, 0, 3, 21, 0, 0, 3, 122,
		})

	// Function
	function := Func{}
	function.Pid = Pid{Atom("Node"), 123456789, 987654321, 128}
	function.Module = Atom("module")
	function.Index = 1234567
	function.FreeVars = []Term{1, 2.3, "string"}
	function.Uniq = 987654321
	assertEncode(t, function,
		[]byte{131, 117, 0, 0, 0, 3, 103, 100, 0, 4, 78,
			111, 100, 101, 7, 91, 205, 21, 58, 222,
			104, 177, 128, 100, 0, 6, 109, 111, 100,
			117, 108, 101, 98, 0, 18, 214, 135, 98, 58,
			222, 104, 177, 97, 1, 70, 64, 2, 102, 102,
			102, 102, 102, 102, 107, 0, 6, 115, 116,
			114, 105, 110, 103,
		})

	// New Function
	newFunction := NewFunc{}
	newFunction.Uniq = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	newFunction.Arity = 2
	newFunction.FreeVars = []Term{1, 2.3, "string"}
	newFunction.Index = 1234567
	newFunction.Module = Atom("module")
	newFunction.OldIndex = 1234567
	newFunction.OldUnique = 987654321
	newFunction.Pid = Pid{Atom("Node"), 123456789, 987654321, 128}
	assertEncode(t, newFunction,
		[]byte{131, 112, 0, 0, 0, 85, 2, 1, 2, 3, 4, 5, 6,
			7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 0, 18,
			214, 135, 0, 0, 0, 3, 100, 0, 6, 109, 111,
			100, 117, 108, 101, 98, 0, 18, 214, 135, 98,
			58, 222, 104, 177, 103, 100, 0, 4, 78, 111,
			100, 101, 7, 91, 205, 21, 58, 222, 104, 177,
			128, 97, 1, 70, 64, 2, 102, 102, 102, 102, 102,
			102, 107, 0, 6, 115, 116, 114, 105, 110, 103,
		})

	// Export
	export := Export{}
	export.Module = Atom("module")
	export.Arity = 2
	export.Function = Atom("function")
	assertEncode(t, export,
		[]byte{131, 113, 100, 0, 6, 109, 111, 100, 117, 108,
			101, 100, 0, 8, 102, 117, 110, 99, 116, 105,
			111, 110, 97, 2,
		})
}

func TestEncodeWithCompression(t *testing.T) {

	// Test a compressed list
	assertEncodeAndCompress(t, []Term{Atom("foo, 世界"), Atom("foo, 世界"),
		Atom("foo, 世界"), Atom("foo, 世界"),
		Atom("foo, 世界"), Atom("foo, 世界")},
		[]byte{131, 80, 0, 0, 0, 86, 120, 156, 202, 96, 43,
			99, 224, 78, 203, 207, 215, 81, 120, 178, 99,
			218, 243, 169, 61, 212, 225, 1, 2, 0, 0, 255,
			255, 48, 96, 38, 49,
		}, true)
}

func TestMarshal(t *testing.T) {
	var buf bytes.Buffer
	Marshal(&buf, 42)
	assertEqual(t, []byte{131, 97, 42}, buf.Bytes())
}

func TestMarshalResponse(t *testing.T) {
	var buf bytes.Buffer
	MarshalResponse(&buf, []Term{Atom("reply"), 42})
	assertEqual(t, []byte{0, 0, 0, 13,
		131, 104, 2,
		100, 0, 5, 114, 101, 112, 108,
		121, 97, 42,
	},
		buf.Bytes())
}

func assertEncode(t *testing.T, actual interface{}, expected []byte) {
	assertEncodeAndCompress(t, actual, expected, false)
}

func assertEncodeAndCompress(t *testing.T, actual interface{}, expected []byte, compress bool) {
	assertEncodeUsingMinor(t, actual, expected, compress, MinorVersion1)
}

func assertEncodeUsingMinor(t *testing.T, actual interface{}, expected []byte, compress bool, minorVersion int) {
	val, err := EncodeAndCompressUsingMinorVersion(actual, compress, minorVersion)
	if err != nil {
		t.Errorf("Encode(%v) returned error '%v'", actual, err)
	} else if !reflect.DeepEqual(val, expected) {
		t.Errorf("Decode(%v) = %v, expected %v", actual, val, expected)
	}
}
