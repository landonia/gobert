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
}

func TestEncodeWithCompression(t *testing.T) {
	// TODO
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
	assertEncodeUsingMinor(t, actual, expected, false, MinorVersion1)
}

func assertEncodeUsingMinor(t *testing.T, actual interface{}, expected []byte, compress bool, minorVersion int) {
	val, err := EncodeAndCompressUsingMinorVersion(actual, compress, minorVersion)
	if err != nil {
		t.Errorf("Encode(%v) returned error '%v'", actual, err)
	} else if !reflect.DeepEqual(val, expected) {
		t.Errorf("Decode(%v) = %v, expected %v", actual, val, expected)
	}
}
