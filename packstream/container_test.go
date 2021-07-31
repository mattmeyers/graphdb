package packstream

import (
	"bytes"
	"reflect"
	"testing"
)

func TestList_MarshalPackstream(t *testing.T) {
	tests := []struct {
		name    string
		l       List
		want    []byte
		wantErr bool
	}{
		{
			name:    "empty list",
			l:       List{},
			want:    []byte{0x90},
			wantErr: false,
		},
		{
			name:    "list with mixed elements",
			l:       List{true, "a"},
			want:    []byte{0x92, 0xC3, 0x81, 0x61},
			wantErr: false,
		},
		{
			name: "non short list",
			l: List{
				true, true, true, true,
				true, true, true, true,
				true, true, true, true,
				true, true, true, true,
				true, true,
			},
			want: []byte{
				0xD4, 0x12, 0xC3, 0xC3,
				0xC3, 0xC3, 0xC3, 0xC3,
				0xC3, 0xC3, 0xC3, 0xC3,
				0xC3, 0xC3, 0xC3, 0xC3,
				0xC3, 0xC3, 0xC3, 0xC3,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.l.MarshalPackstream()
			if (err != nil) != tt.wantErr {
				t.Errorf("List.MarshalPackstream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List.MarshalPackstream() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDictionary_MarshalPackstream(t *testing.T) {
	// Because Go's maps don't provide a guaranteed ordering, we check that each
	// key-value pair is contained in the result rather than checking the exact
	// return value.
	tests := []struct {
		name    string
		d       Dictionary
		want    [][]byte
		wantErr bool
	}{
		{
			name:    "empty dictionary",
			d:       Dictionary{},
			want:    [][]byte{{0xA0}},
			wantErr: false,
		},
		{
			name: "dictionary with mixed values",
			d:    Dictionary{"a": true, "b": "abc"},
			want: [][]byte{
				{0xA2},
				{0x81, 0x61, 0xC3},
				{0x81, 0x62, 0x83, 0x61, 0x62, 0x63},
			},
			wantErr: false,
		},
		{
			name: "Non short dictionary",
			d: Dictionary{
				"a": true,
				"b": true,
				"c": true,
				"d": true,
				"e": true,
				"f": true,
				"g": true,
				"h": true,
				"i": true,
				"j": true,
				"k": true,
				"l": true,
				"m": true,
				"n": true,
				"o": true,
				"p": true,
				"q": true,
			},
			want: [][]byte{
				{0xD8, 0x11},
				{0x81, 0x61, 0xC3},
				{0x81, 0x62, 0xC3},
				{0x81, 0x63, 0xC3},
				{0x81, 0x64, 0xC3},
				{0x81, 0x65, 0xC3},
				{0x81, 0x66, 0xC3},
				{0x81, 0x67, 0xC3},
				{0x81, 0x68, 0xC3},
				{0x81, 0x69, 0xC3},
				{0x81, 0x6A, 0xC3},
				{0x81, 0x6B, 0xC3},
				{0x81, 0x6C, 0xC3},
				{0x81, 0x6D, 0xC3},
				{0x81, 0x6E, 0xC3},
				{0x81, 0x6F, 0xC3},
				{0x81, 0x70, 0xC3},
				{0x81, 0x71, 0xC3},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.MarshalPackstream()
			if (err != nil) != tt.wantErr {
				t.Errorf("Dictionary.MarshalPackstream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, sub := range tt.want {
				if !bytes.Contains(got, sub) {
					t.Errorf("Dictionary.MarshalPackstream() = %x, want %x", got, tt.want)
				}
			}
		})
	}
}
