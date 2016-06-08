package bert

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/big"
	"reflect"
	"unicode/utf8"
)

func write1(w io.Writer, ui8 uint8) { w.Write([]byte{ui8}) }

func write2(w io.Writer, ui16 uint16) {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, ui16)
	w.Write(b)
}

func write4(w io.Writer, ui32 uint32) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, ui32)
	w.Write(b)
}

func writeSmallInt(w io.Writer, n uint8) {
	write1(w, SmallIntTag)
	write1(w, n)
}

func writeInt(w io.Writer, n uint32) {
	write1(w, IntTag)
	write4(w, n)
}

func writeFloat(w io.Writer, f float32) {
	write1(w, FloatTag)

	s := fmt.Sprintf("%.20e", float32(f))
	w.Write([]byte(s))

	pad := make([]byte, 31-len(s))
	w.Write(pad)
}

func writeNewFloat(w io.Writer, f float64) {
	write1(w, NewFloatTag)

	ui64 := math.Float64bits(f)
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, ui64)
	w.Write(b)
}

func writeAtom(w io.Writer, a string) {
	write1(w, AtomTag)
	write2(w, uint16(len(a)))
	w.Write([]byte(a))
}

func writeAtomUtf8(w io.Writer, a string) {
	write1(w, AtomUtf8Tag)
	write2(w, uint16(len(a)))
	w.Write([]byte(a))
}

func writeSmallTuple(w io.Writer, t reflect.Value, minorVersion int) {
	write1(w, SmallTupleTag)
	size := t.Len()
	write1(w, uint8(size))

	for i := 0; i < size; i++ {
		writeTag(w, t.Index(i), minorVersion)
	}
}

func writeNil(w io.Writer) { write1(w, NilTag) }

func writeString(w io.Writer, s string) {
	write1(w, StringTag)
	write2(w, uint16(len(s)))
	w.Write([]byte(s))
}

func writeList(w io.Writer, l reflect.Value, minorVersion int) {
	write1(w, ListTag)
	size := l.Len()
	write4(w, uint32(size))

	for i := 0; i < size; i++ {
		writeTag(w, l.Index(i), minorVersion)
	}

	writeNil(w)
}

func writeMap(w io.Writer, l reflect.Value, minorVersion int) {
	write1(w, MapTag)
	keys := l.MapKeys()
	len := uint32(len(keys))
	write4(w, len)

	for i := uint32(0); i < len; i++ {
		writeTag(w, keys[i], minorVersion)
		writeTag(w, l.MapIndex(keys[i]), minorVersion)
	}
}

func writeNode(w io.Writer, node Atom) {

	// If the number of runes equal the number of bytes then UTF-8 is not required
	if utf8.RuneCount([]byte(node)) == len(node) {
		writeAtom(w, string(node))
	} else {
		writeAtomUtf8(w, string(node))
	}
}

func writeReference(w io.Writer, l Reference) {
	write1(w, ReferenceTag)
	writeNode(w, l.Node)
	write4(w, l.ID)
	write1(w, l.Creation)
}

func writeNewReference(w io.Writer, l NewReference) {
	write1(w, NewReferenceTag)
	write2(w, uint16(len(l.ID)))
	writeNode(w, l.Node)
	write1(w, l.Creation)
	for i := 0; i < len(l.ID); i++ {
		write4(w, l.ID[i])
	}
}

func writeLargeBigNum(w io.Writer, l big.Int) {
	write1(w, LargeBignumTag)
	write4(w, uint32(len(l.Bytes())))
	writeBigNum(w, l)
}

func writeSmallBigNum(w io.Writer, l big.Int) {
	write1(w, SmallBignumTag)
	write1(w, uint8(len(l.Bytes())))
	writeBigNum(w, l)
}

func writeBigNum(w io.Writer, l big.Int) {
	if l.Sign() < 0 {
		write1(w, 1)
	} else {
		write1(w, 0)
	}

	// The big.Int uses BigEndian for the bytes and the format is
	// expecting LittleEndian. Reverse the bytes.
	bEndBits := l.Bytes()
	var lEndBits []byte
	for i := len(bEndBits) - 1; i >= 0; i-- {
		lEndBits = append(lEndBits, bEndBits[i])
	}
	w.Write(lEndBits)
}

func writePort(w io.Writer, l Port) {
	write1(w, PortTag)
	writeNode(w, l.Node)
	write4(w, l.ID)
	write1(w, l.Creation)
}

func writePid(w io.Writer, l Pid) {
	write1(w, PidTag)
	writeNode(w, l.Node)
	write4(w, l.ID)
	write4(w, l.Serial)
	write1(w, l.Creation)
}

func writeFunc(w io.Writer, l Func, minorVersion int) {
	write1(w, FunTag)
	write4(w, uint32(len(l.FreeVars)))
	writePid(w, l.Pid)
	writeNode(w, l.Module)
	writeInt(w, l.Index)
	writeInt(w, l.Uniq)
	for i := 0; i < len(l.FreeVars); i++ {
		writeTag(w, reflect.ValueOf(l.FreeVars[i]), minorVersion)
	}
}

func writeNewFunc(w io.Writer, l NewFunc, minorVersion int) {
	write1(w, NewFunTag)

	// Create a new buffer to write the bytes to so the
	// total size can be written
	buf := bytes.NewBuffer([]byte{})
	write1(buf, l.Arity)
	buf.Write(l.Uniq)
	write4(buf, l.Index)
	write4(buf, uint32(len(l.FreeVars)))
	writeNode(buf, l.Module)
	writeInt(buf, l.OldIndex)
	writeInt(buf, l.OldUnique)
	writePid(buf, l.Pid)
	for i := 0; i < len(l.FreeVars); i++ {
		writeTag(buf, reflect.ValueOf(l.FreeVars[i]), minorVersion)
	}

	write4(w, uint32(buf.Len()+4))
	w.Write(buf.Bytes())
}

func writeExport(w io.Writer, l Export) {
	write1(w, ExportTag)
	writeNode(w, l.Module)
	writeNode(w, l.Function)
	writeSmallInt(w, l.Arity)
}

func writeTag(w io.Writer, val reflect.Value, minorVersion int) (err error) {
	switch v := val; v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n := v.Int()
		if n >= 0 && n < 256 {
			writeSmallInt(w, uint8(n))
		} else {
			writeInt(w, uint32(n))
		}
	case reflect.Float32, reflect.Float64:
		if minorVersion == MinorVersion0 {
			writeFloat(w, float32(v.Float()))
		} else {
			writeNewFloat(w, v.Float())
		}
	case reflect.String:
		if v.Type().Name() == "Atom" {

			// If the number of runes equal the number of bytes then UTF-8 is not required
			if utf8.RuneCount([]byte(v.String())) == len(v.String()) {
				writeAtom(w, v.String())
			} else {
				writeAtomUtf8(w, v.String())
			}
		} else {
			writeString(w, v.String())
		}
	case reflect.Slice:
		writeSmallTuple(w, v, minorVersion)
	case reflect.Array:
		writeList(w, v, minorVersion)
	case reflect.Interface, reflect.Ptr:
		writeTag(w, v.Elem(), minorVersion)
	case reflect.Map:
		writeMap(w, v, minorVersion)
	case reflect.Struct:
		vali := v.Interface()
		switch rVal := vali.(type) {
		case DistributionHeader:
			err = ErrUnknownType
		case Reference:
			writeReference(w, rVal)
		case NewReference:
			writeNewReference(w, rVal)
		case big.Int:
			if len(rVal.Bytes()) >= 1<<8 {
				writeLargeBigNum(w, rVal)
			} else {
				writeSmallBigNum(w, rVal)
			}
		case Port:
			writePort(w, rVal)
		case Pid:
			writePid(w, rVal)
		case Func:
			writeFunc(w, rVal, minorVersion)
		case NewFunc:
			writeNewFunc(w, rVal, minorVersion)
		case Export:
			writeExport(w, rVal)
		}
	default:
		if !reflect.Indirect(val).IsValid() {
			writeNil(w)
		} else {
			err = ErrUnknownType
		}
	}

	return
}

// EncodeTo encodes val and writes it to w, returning any error.
func EncodeTo(w io.Writer, val interface{}) (err error) {
	return EncodeToUsingMinorVersion(w, val, MinorVersion1)
}

// Encode encodes val and returns it or an error.
func Encode(val interface{}) ([]byte, error) {
	return EncodeUsingMinorVersion(val, MinorVersion1)
}

// Marshal is an alias for EncodeTo.
func Marshal(w io.Writer, val interface{}) error {
	return MarshalUsingMinorVersion(w, val, MinorVersion1)
}

// MarshalResponse encodes val into a BURP Response struct and writes it to w,
// returning any error.
func MarshalResponse(w io.Writer, val interface{}) (err error) {
	return MarshalResponseUsingMinorVersion(w, val, MinorVersion1)
}

// EncodeToAndCompress encodes val and writes it to w, returning any error.
// If compress is true the body will be compressed
func EncodeToAndCompress(w io.Writer, val interface{}, compress bool) (err error) {
	return EncodeToAndCompressUsingMinorVersion(w, val, compress, MinorVersion1)
}

// EncodeTo encodes val and writes it to w, returning any error.
func EncodeToUsingMinorVersion(w io.Writer, val interface{}, minorVersion int) (err error) {
	return EncodeToAndCompressUsingMinorVersion(w, val, false, minorVersion)
}

// Encode encodes val and returns it or an error.
func EncodeUsingMinorVersion(val interface{}, minorVersion int) ([]byte, error) {
	return EncodeAndCompressUsingMinorVersion(val, false, minorVersion)
}

// Marshal is an alias for EncodeTo.
func MarshalUsingMinorVersion(w io.Writer, val interface{}, minorVersion int) error {
	return EncodeToUsingMinorVersion(w, val, minorVersion)
}

// MarshalResponse encodes val into a BURP Response struct and writes it to w,
// returning any error.
func MarshalResponseUsingMinorVersion(w io.Writer, val interface{}, minorVersion int) (err error) {
	return MarshalResponseAndCompressUsingMinorVersion(w, val, false, minorVersion)
}

// EncodeToAndCompress encodes val and writes it to w, returning any error.
// If compress is true the body will be compressed
func EncodeToAndCompressUsingMinorVersion(w io.Writer, val interface{}, compress bool, minorVersion int) (err error) {
	write1(w, VersionTag)
	if compress {

		// Write the bytes to a buffer (the original length is required)
		buf := bytes.NewBuffer([]byte{})
		err = writeTag(buf, reflect.ValueOf(val), minorVersion)
		if err == nil {
			write1(w, CompressedTag)
			write4(w, uint32(buf.Len()))
			zw := zlib.NewWriter(w)
			_, err = zw.Write(buf.Bytes())
			zw.Close()
		}
	} else {
		// Write directly to the writer
		err = writeTag(w, reflect.ValueOf(val), minorVersion)
	}
	return
}

// EncodeAndCompress encodes val and returns it or an error.
// If compress is true the body will be compressed
func EncodeAndCompress(val interface{}, compress bool) ([]byte, error) {
	return EncodeAndCompressUsingMinorVersion(val, compress, MinorVersion1)
}

// MarshalAndCompress is an alias for EncodeTo.
// If compress is true the body will be compressed
func MarshalAndCompress(w io.Writer, val interface{}, compress bool) error {
	return MarshalAndCompressUsingMinorVersion(w, val, compress, MinorVersion1)
}

// MarshalResponseAndCompress encodes val into a BURP Response struct and writes it to w,
// returning any error.
// If compress is true the body will be compressed
func MarshalResponseAndCompress(w io.Writer, val interface{}, compress bool) (err error) {
	return MarshalResponseAndCompressUsingMinorVersion(w, val, compress, MinorVersion1)
}

// EncodeAndCompressUsingMinorVersion encodes val and returns it or an error.
// If compress is true the body will be compressed
// It will use the minor version for the encoding
func EncodeAndCompressUsingMinorVersion(val interface{}, compress bool, minorVersion int) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	err := EncodeToAndCompressUsingMinorVersion(buf, val, compress, minorVersion)
	return buf.Bytes(), err
}

// MarshalAndCompress is an alias for EncodeTo.
// If compress is true the body will be compressed
// It will use the minor version for the encoding
func MarshalAndCompressUsingMinorVersion(w io.Writer, val interface{}, compress bool, minorVersion int) error {
	return EncodeToAndCompressUsingMinorVersion(w, val, compress, minorVersion)
}

// MarshalResponseAndCompress encodes val into a BURP Response struct and writes it to w,
// returning any error.
// If compress is true the body will be compressed
// It will use the minor version for the encoding
func MarshalResponseAndCompressUsingMinorVersion(w io.Writer, val interface{}, compress bool, minorVersion int) (err error) {
	resp, err := EncodeAndCompressUsingMinorVersion(val, compress, minorVersion)

	write4(w, uint32(len(resp)))
	w.Write(resp)

	return
}
