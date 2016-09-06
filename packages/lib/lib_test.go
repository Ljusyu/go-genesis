package lib

import (
	"bytes"
	"github.com/DayLightProject/go-daylight/packages/test"
	"math/rand"
	"testing"
	"encoding/hex"
)

type ByteTest struct {
	src  []byte
	want []byte
}

type EncodeType struct {
	value int64
	data  []byte
}

var testList = []EncodeType{
	{0, []byte{0}},
	{1, []byte{1, 1}},
	{127, []byte{1, 0x7f}},
	{65000, []byte{2, 0xe8, 0xfd}},
	{156507890, []byte{4, 0xf2, 0x1e, 0x54, 0x09}},
	{1565073467890890, []byte{7, 0xca, 0xdc, 0x19, 0x10, 0x6d, 0x8f, 0x05}},
}

func TestEncodeLenInt64(t *testing.T) {
	var off int
	buf := make([]byte, 0)
	for _, val := range testList {
		off = len(buf)
		EncodeLenInt64(&buf, val.value)
		if bytes.Compare(buf[off:len(buf)], val.data) != 0 {
			t.Errorf("different slice %d", val.value)
		}
	}
}

func TestDecodeLenInt64(t *testing.T) {
	for _, val := range testList {
		buf := val.data
		x, err := DecodeLenInt64(&buf)
		if err != nil {
			t.Error(err.Error)
		}
		if x != val.value {
			t.Errorf("different int64 %d != %d", x, val.value)
		}
	}
}

func TestAddress(t *testing.T) {
	for i := 0; i < 50; i++ {
		key, seed := test.RandBytes(64)
		address := KeyToAddress(key)
		if (i % 10) == 0 {
			if IsValidAddress(address[:len(address)-1]) {
				t.Errorf("valid address %s for %x seed: %d", address[:len(address)-1], key, seed)
			}
		} else if !IsValidAddress(address) {
			t.Errorf("not valid address %s for %x seed: %d", address, key, seed)
		}
	}
}

func TestEncodeDecodeLength(t *testing.T) {
	vals := []int64{1, 67, 127, 128, 256, 1024, 2000, 10000, 65000, 1000000, 0xffeeffff,
		8123498762, 25000060000, 400000000035, -10000000044546, -1}
	for _, i := range vals {
		result := EncodeLength(i)
		got, _ := DecodeLength(&result)
		if got != i {
			t.Errorf("wrong length encoding %d != %d", i, got)
		}
	}
	if length, _ := DecodeLength(&[]byte{}); length != 0 {
		t.Errorf("wrong decoding empty slice")
	}

}

func TestFill(t *testing.T) {
	for i := 0; i < 50; i++ {
		size := rand.Intn(33)
		input, _ := test.RandBytes(size)

		out := FillLeft(input)
		if bytes.Compare(out[:32-size], make([]byte, 32-size)) != 0 ||
			bytes.Compare(out[32-size:], input) != 0 {
			t.Errorf(`different slices %x %x`, input, out)
		}
	}
}

func TestEncodeBinary(t *testing.T) {
	var (
		out []byte
		off int
	)
	check := func( format string, cmp []byte, args ...interface{}) {
		if err := EncodeBinary(&out, format, args...); err!=nil {
			t.Errorf(err.Error())
		} else if bytes.Compare(out[off:], cmp) != 0 {
			t.Errorf(`different output binary data %x`, out )
		}
		off = len(out)
	}
	check( `1`, []byte{255}, 255)
	check( `414`, []byte{0x01,0x01,0,0, 0x7e, 0xa1,0x86,1,0}, 257, 126, 100001 )
	check( `ii4i`, []byte{0x01,0x43, 0x3,0x9a,0x31,1, 0xff,0xff,0,0, 0x3,0x2c,0xdd,0x15}, 
	               67, 78234, 0xffff, int64(1432876))
	cmp,_ := hex.DecodeString(`0474657374c8057b0001ff86`)
	check( `s1s`, cmp, `test`, 200, []byte{ 123, 0, 1, 255, 134})
//	fmt.Printf( "\r\n%x", out[24:])
}
