package bencode

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	val := "abc"
	buf := new(bytes.Buffer)
	wLen := EncodeString(buf, val)
	assert.Equal(t, 5, wLen)
	str, _ := DecodeString(buf)
	assert.Equal(t, val, str)

	val = ""
	for i := 0; i < 20; i++ {
		val += string(byte('a' + i))
	}
	buf.Reset()
	wLen = EncodeString(buf, val)
	assert.Equal(t, 23, wLen)
	str, _ = DecodeString(buf)
	assert.Equal(t, val, str)
}

func TestInt(t *testing.T) {
	val := 999
	buf := new(bytes.Buffer)
	wLen := EncodeInt(buf, val)
	assert.Equal(t, 5, wLen)
	iv, _ := DecodeInt(buf)
	assert.Equal(t, val, iv)

	val = 0
	buf.Reset()
	wLen = EncodeInt(buf, val)
	assert.Equal(t, 3, wLen)
	iv, _ = DecodeInt(buf)
	assert.Equal(t, val, iv)

	val = -99
	buf.Reset()
	wLen = EncodeInt(buf, val)
	assert.Equal(t, 5, wLen)
	iv, _ = DecodeInt(buf)
	assert.Equal(t, val, iv)
}

func TestBencode(t *testing.T) {
	testCases := []struct {
		name      string
		input     *BObject
		wantError error
		wantLen   int
	}{
		{
			name: "empty string",
			input: &BObject{
				type_: BSTR,
				val_:  "",
			},
			wantError: nil,
			wantLen:   2, // "0:"
		},
		{
			name: "string",
			input: &BObject{
				type_: BSTR,
				val_:  "Hello, world!",
			},
			wantError: nil,
			wantLen:   16, // "13:Hello, world!"
		},
		{
			name: "empty list",
			input: &BObject{
				type_: BLIST,
				val_:  []*BObject{},
			},
			wantError: nil,
			wantLen:   2, // "le"
		},
		{
			name: "list",
			input: &BObject{
				type_: BLIST,
				val_: []*BObject{
					{type_: BSTR, val_: "hello"},
					{type_: BINT, val_: 123},
				},
			},
			wantError: nil,
			wantLen:   14, // "l5:helloi123ee"
		},
		{
			name: "empty dict",
			input: &BObject{
				type_: BDICT,
				val_:  map[string]*BObject{},
			},
			wantError: nil,
			wantLen:   2, // "de"
		},
		{
			name: "dict",
			input: &BObject{
				type_: BDICT,
				val_: map[string]*BObject{
					"hello": {type_: BSTR, val_: "world"},
					"num":   {type_: BINT, val_: 123},
				},
			},
			wantError: nil,
			wantLen:   26, // "d5:hello5:world3:numi123ee"
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bb := &bytes.Buffer{}
			gotLen := tc.input.Bencode(bb)
			if gotLen != tc.wantLen {
				t.Errorf("Bencode() got len = %d, want %d", bb.Len(), tc.wantLen)
			}
		})
	}
}
