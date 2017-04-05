package bert

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
	"math/big"
)

func ExampleDecode() {
	i, err := Decode([]byte{131, 97, 42})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%#v\n", i)
	s, err := Decode([]byte{131, 107, 0, 3, 102, 111, 111})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%#v\n", s)
	a, err := Decode([]byte{131, 104, 1, 100, 0, 3, 102, 111, 111})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%#v\n", a)
	// Output:
	// 42
	// "foo"
	// []bert.Term{"foo"}
}

func TestDecodeNoCompression(t *testing.T) {
	// Small Integer
	assertDecode(t, []byte{131, 97, 1}, 1)
	assertDecode(t, []byte{131, 97, 2}, 2)
	assertDecode(t, []byte{131, 97, 3}, 3)
	assertDecode(t, []byte{131, 97, 4}, 4)
	assertDecode(t, []byte{131, 97, 42}, 42)

	// Integer
	assertDecode(t, []byte{131, 98, 0, 0, 1, 1}, 257)
	assertDecode(t, []byte{131, 98, 0, 0, 4, 1}, 1025)
	assertDecode(t, []byte{131, 98, 255, 255, 255, 255}, -1)
	assertDecode(t, []byte{131, 98, 255, 255, 255, 248}, -8)
	assertDecode(t, []byte{131, 98, 0, 0, 19, 136}, 5000)
	assertDecode(t, []byte{131, 98, 255, 255, 236, 120}, -5000)

	// Small Bignum
	assertDecode(t, []byte{131, 110, 4, 0, 177, 104, 222, 58},
		*big.NewInt(987654321))
	assertDecode(t, []byte{131, 110, 4, 1, 177, 104, 222, 58},
		*big.NewInt(-987654321))

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
	assertDecode(t, []byte{131, 111, 0, 0, 1, 44, 0, 1, 1, 1, 1, 1, 1, 1, 1,
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
		1, 1, 1, 1, 1, 1},
		bigNum)

	// Put the number as negative
	bigNum.Neg(&bigNum)
	assertDecode(t, []byte{131, 111, 0, 0, 1, 44, 1, 1, 1, 1, 1, 1, 1, 1, 1,
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
		1, 1, 1, 1, 1, 1},
		bigNum)

	// Float
	assertDecode(t, []byte{131, 99, 53, 46, 48, 48, 48, 48, 48, 48, 48,
		48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 101, 45,
		48, 49, 0, 0, 0, 0, 0,
	},
		float32(0.5))
	assertDecode(t, []byte{131, 99, 51, 46, 49, 52, 49, 53, 56,
		57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 56, 56, 50, 54, 50,
		101, 43, 48, 48, 0, 0, 0, 0, 0,
	},
		float32(3.14159))
	assertDecode(t, []byte{131, 99, 45, 51, 46, 49, 52, 49, 53,
		56, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 56, 56, 50, 54,
		50, 101, 43, 48, 48, 0, 0, 0, 0,
	},
		float32(-3.14159))

	// NewFloat
	assertDecode(t,
		[]byte{131, 70, 63, 224, 0, 0, 0, 0, 0, 0,
	},
		0.5)
	assertDecode(t,
		[]byte{131, 70, 64, 9, 33, 249, 240, 27, 134, 110,
	},
		3.14159)
	assertDecode(t,
		[]byte{131, 70, 192, 9, 33, 249, 240, 27, 134, 110,
	},
		-3.14159)

	// Atom
	assertDecode(t, []byte{131, 100, 0, 3, 102, 111, 111},
		Atom("foo"))
	assertDecode(t, []byte{131, 100, 0, 5, 104, 101, 108, 108, 111},
		Atom("hello"))

	// Atom UTF8
	assertDecode(t,
		[]byte{131, 118, 0, 11, 102, 111, 111, 44, 32, 228,
			184, 150, 231, 149, 140},
		Atom("foo, 世界"))

	// Small Tuple
	assertDecode(t, []byte{131, 104, 0}, []Term{})
	assertDecode(t, []byte{131, 104, 1,
		100, 0, 3, 102, 111, 111,
	},
		[]Term{Atom("foo")})
	assertDecode(t, []byte{131, 104, 2,
		100, 0, 3, 102, 111, 111,
		100, 0, 3, 98, 97, 114,
	},
		[]Term{Atom("foo"), Atom("bar")})
	assertDecode(t, []byte{131, 104, 3,
		100, 0, 5, 99, 111, 111, 114, 100,
		97, 23,
		97, 42,
	},
		[]Term{Atom("coord"), 23, 42})
	assertDecode(t, []byte{131, 104, 4,
		100, 0, 4, 99, 97, 108, 108,
		100, 0, 6, 112, 104, 111, 116, 111, 120,
		100, 0, 8, 105, 109, 103, 95, 115, 105, 122, 101,
		108, 0, 0, 0, 1, 97, 99,
		106,
	},
		[]Term{Atom("call"), Atom("photox"), Atom("img_size"), []Term{99}})

	// Small Tuple with UTF8
	assertDecode(t,
		[]byte{131, 104, 1, 118, 0, 11, 102, 111, 111, 44, 32, 228,
			184, 150, 231, 149, 140},
		[]Term{Atom("foo, 世界")})

	// Large Tuple

	// String
	assertDecode(t, []byte{131, 107, 0, 3, 102, 111, 111}, "foo")
	assertDecode(t, []byte{131, 107, 0, 1, 0}, "\000")
	assertDecode(t, []byte{131, 107, 0, 1, 1}, "\001")

	// List
	assertDecode(t, []byte{131, 106}, []Term{})
	assertDecode(t, []byte{131, 108, 0, 0, 0, 1, 97, 1, 106},
		[]Term{1})
	assertDecode(t, []byte{131, 108, 0, 0, 0, 1, 98, 0, 0, 1, 0, 106},
		[]Term{256})
	assertDecode(t, []byte{131, 108, 0, 0, 0, 1, 107, 0, 1, 97, 106},
		[]Term{"a"})
	assertDecode(t, []byte{131, 108, 0, 0, 0, 1, 100, 0, 1, 97, 106},
		[]Term{Atom("a")})
	assertDecode(t, []byte{131, 108, 0, 0, 0, 3,
		97, 1, 97, 2, 97, 3,
		106,
	},
		[]Term{1, 2, 3})
	assertDecode(t, []byte{131, 108, 0, 0, 0, 2,
		107, 0, 1, 97, 107, 0, 1, 98, 106,
	},
		[]Term{"a", "b"})
	assertDecode(t, []byte{131, 108, 0, 0, 0, 1, 107, 0, 2, 97, 98, 106},
		[]Term{"ab"})
	assertDecode(t, []byte{131, 108, 0, 0, 0, 2,
		100, 0, 1, 97, 100, 0, 1, 98, 106,
	},
		[]Term{Atom("a"), Atom("b")})
	assertDecode(t, []byte{131, 108, 0, 0, 0, 2, 100, 0, 1, 97, 97, 1, 106},
		[]Term{Atom("a"), 1})
	assertDecode(t, []byte{131, 108, 0, 0, 0, 2,
		107, 0, 1, 97, 107, 0, 2, 1, 2, 106,
	},
		[]Term{"a", "\001\002"})
	assertDecode(t, []byte{131, 108, 0, 0, 0, 2,
		100, 0, 1, 97, 107, 0, 2, 1, 2, 106,
	},
		[]Term{Atom("a"), "\001\002"})
	assertDecode(t, []byte{131, 108, 0, 0, 0, 2, 100, 0, 1, 97, 108, 0, 0, 0, 1, 98, 0, 0, 1, 0, 106,
		106,
	},
		[]Term{Atom("a"), []Term{256}})

	// Binary
	assertDecode(t, []uint8{131, 109, 0, 0, 0, 3, 102, 111, 111},
		bintag{102, 111, 111})
	assertDecode(t, []byte{131, 109, 0, 0, 0, 5, 104, 101, 108, 108, 111},
		bintag{104, 101, 108, 108, 111})

	// Complex
	assertDecode(t, []byte{131, 104, 2, 100, 0, 4, 98, 101, 114, 116, 100, 0, 3, 110, 105, 108}, nil)
	assertDecode(t, []byte{131, 104, 2, 100, 0, 4, 98, 101, 114, 116, 100, 0, 4, 116, 114, 117, 101}, true)
	assertDecode(t, []byte{131, 104, 2, 100, 0, 4, 98, 101, 114, 116, 100, 0, 5, 102, 97, 108, 115, 101}, false)

	assertDecode(t, []byte{131, 104, 4,
		100, 0, 4, 99, 97, 108, 108,
		100, 0, 6, 112, 104, 111, 116, 111, 120,
		100, 0, 8, 105, 109, 103, 95, 115, 105, 122, 101,
		108, 0, 0, 0, 1, 97, 99,
		106,
	},
		[]Term{Atom("call"), Atom("photox"), Atom("img_size"), []Term{99}})

	// Map
	assertDecode(t,
		[]byte{131, 116, 0, 0, 0, 3, 107, 0, 4, 107, 101, 121, 49, 107, 0, 6, 118,
			97, 108, 117, 101, 49, 107, 0, 4, 107, 101, 121, 50, 107, 0, 6, 118,
			97, 108, 117, 101, 50, 107, 0, 4, 107, 101, 121, 51, 107, 0, 6, 118,
			97, 108, 117, 101, 51,
		}, maptag{"key1":"value1", "key2":"value2", "key3":"value3"})

	// Pid
	pid := Pid{}
	pid.Node = Atom("Node")
	pid.ID = 123456789
	pid.Serial = 987654321
	pid.Creation = 128
	assertDecode(t,
		[]byte{131, 103, 100, 0, 4, 78, 111,
			100, 101, 7, 91, 205, 21, 58,
			222, 104, 177, 128,
		}, pid)

	// Port
	port := Port{}
	port.Node = Atom("Node")
	port.ID = 123456789
	port.Creation = 128
	assertDecode(t,
		[]byte{131, 102, 100, 0, 4, 78, 111, 100,
			101, 7, 91, 205, 21, 128,
		}, port)

	// Reference
	reference := Reference{}
	reference.Node = Atom("Node")
	reference.ID = 123456789
	reference.Creation = 128
	assertDecode(t,
		[]byte{131, 101, 100, 0, 4, 78, 111, 100,
			101, 7, 91, 205, 21, 128,
		}, reference)

	// New Reference
	newReference := NewReference{}
	newReference.Creation = 128
	newReference.ID = []uint32{123,234,345,456,567,678,789,890}
	newReference.Node = Atom("node")
	assertDecode(t,
		[]byte{131, 114, 0, 8, 100, 0, 4, 110, 111, 100,
			101, 128, 0, 0, 0, 123, 0, 0, 0, 234, 0,
			0, 1, 89, 0, 0, 1, 200, 0, 0, 2, 55, 0, 0,
			2, 166, 0, 0, 3, 21, 0, 0, 3, 122,
		}, newReference)

	// Function
	function := Func{}
	function.Pid = Pid{Atom("Node"),123456789,987654321,128}
	function.Module = Atom("module")
	function.Index = 1234567
	function.FreeVars = []Term{1, 2.3, "string"}
	function.Uniq = 987654321
	assertDecode(t,
		[]byte{131, 117, 0, 0, 0, 3, 103, 100, 0, 4, 78,
			111, 100, 101, 7, 91, 205, 21, 58, 222,
			104, 177, 128, 100, 0, 6, 109, 111, 100,
			117, 108, 101, 98, 0, 18, 214, 135, 98, 58,
			222, 104, 177, 97, 1, 70, 64, 2, 102, 102,
			102, 102, 102, 102, 107, 0, 6, 115, 116,
			114, 105, 110, 103,
		}, function)

	// New Function
	newFunction := NewFunc{}
	newFunction.Uniq = []byte{1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16}
	newFunction.Arity = 2
	newFunction.FreeVars = []Term{1, 2.3, "string"}
	newFunction.Index = 1234567
	newFunction.Module = Atom("module")
	newFunction.OldIndex = 1234567
	newFunction.OldUnique = 987654321
	newFunction.Pid = Pid{Atom("Node"),123456789,987654321,128}
	assertDecode(t,
		[]byte{131, 112, 0, 0, 0, 85, 2, 1, 2, 3, 4, 5, 6,
			7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 0, 18,
			214, 135, 0, 0, 0, 3, 100, 0, 6, 109, 111,
			100, 117, 108, 101, 98, 0, 18, 214, 135, 98,
			58, 222, 104, 177, 103, 100, 0, 4, 78, 111,
			100, 101, 7, 91, 205, 21, 58, 222, 104, 177,
			128, 97, 1, 70, 64, 2, 102, 102, 102, 102, 102,
			102, 107, 0, 6, 115, 116, 114, 105, 110, 103,
		}, newFunction)

	// Export
	export := Export{}
	export.Module = Atom("module")
	export.Arity = 2
	export.Function = Atom("function")
	assertDecode(t,
		[]byte{131, 113, 100, 0, 6, 109, 111, 100, 117, 108,
			101, 100, 0, 8, 102, 117, 110, 99, 116, 105,
			111, 110, 97, 2,
		}, export)
}

func TestDecodeWithCompression(t *testing.T) {

	// Test a compressed list
	assertDecode(t,
		[]byte{131, 80, 0, 0, 0, 86, 120, 156, 202, 96, 43,
			99, 224, 78, 203, 207, 215, 81, 120, 178, 99,
			218, 243, 169, 61, 212, 225, 1, 2, 0, 0, 255,
			255, 48, 96, 38, 49,
		}, []Term{Atom("foo, 世界"), Atom("foo, 世界"),
			Atom("foo, 世界"), Atom("foo, 世界"),
			Atom("foo, 世界"), Atom("foo, 世界")})
}

func assertDecode(t *testing.T, data []byte, expected interface{}) {
	val, err := Decode(data)
	if err != nil {
		t.Errorf("Decode(%v) returned error '%v'", data, err)
	} else if !reflect.DeepEqual(val, expected) {
		t.Errorf("Decode(%v) = %#v, expected %#v", data, val, expected)
	}
}

func TestUnmarshal(t *testing.T) {
	var a struct {
		First Atom
	}
	Unmarshal([]byte{131, 104, 1, 100, 0, 3, 102, 111, 111}, &a)
	assertEqual(t, Atom("foo"), a.First)

	var b struct {
		First int
	}
	Unmarshal([]byte{131, 104, 1, 97, 42}, &b)
	assertEqual(t, 42, b.First)

	var c struct {
		First  Atom
		Second Atom
	}
	Unmarshal([]byte{131, 104, 2, 100, 0, 3, 102, 111, 111, 100, 0, 3, 98, 97, 114}, &c)
	assertEqual(t, Atom("foo"), c.First)
	assertEqual(t, Atom("bar"), c.Second)

	var req Request
	Unmarshal([]byte{131, 104, 4,
		100, 0, 4, 99, 97, 108, 108,
		100, 0, 6, 112, 104, 111, 116, 111, 120,
		100, 0, 8, 105, 109, 103, 95, 115, 105, 122, 101,
		108, 0, 0, 0, 1, 97, 99,
		106,
	},
		&req)
	assertEqual(t, Atom("call"), req.Kind)
	assertEqual(t, Atom("photox"), req.Module)
	assertEqual(t, Atom("img_size"), req.Function)
	assertEqual(t, []Term{99}, req.Arguments)
}

func TestUnmarshalRequest(t *testing.T) {
	buf := bytes.NewBuffer([]byte{
		0, 0, 0, 38,
		131, 104, 4,
		100, 0, 4, 99, 97, 108, 108,
		100, 0, 6, 112, 104, 111, 116, 111, 120,
		100, 0, 8, 105, 109, 103, 95, 115, 105, 122, 101,
		108, 0, 0, 0, 1, 97, 99,
		106,
	})

	req, _ := UnmarshalRequest(buf)
	assertEqual(t, Atom("call"), req.Kind)
	assertEqual(t, Atom("photox"), req.Module)
	assertEqual(t, Atom("img_size"), req.Function)
	assertEqual(t, []Term{99}, req.Arguments)
}

func assertEqual(t *testing.T, expected interface{}, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %v, but was %v", expected, actual)
	}
}
