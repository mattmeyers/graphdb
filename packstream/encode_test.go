package packstream

import (
	"reflect"
	"testing"
)

func TestMarshal(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "null",
			args:    args{v: nil},
			want:    []byte{0xC0},
			wantErr: false,
		},
		{
			name:    "false",
			args:    args{v: false},
			want:    []byte{0xC2},
			wantErr: false,
		},
		{
			name:    "true",
			args:    args{v: true},
			want:    []byte{0xC3},
			wantErr: false,
		},
		{
			name:    "bool ponter",
			args:    args{v: boolPtr(true)},
			want:    []byte{0xC3},
			wantErr: false,
		},
		{
			name:    "empty byte slice",
			args:    args{v: []byte{}},
			want:    []byte{0xCC, 0x00},
			wantErr: false,
		},
		{
			name:    "empty byte slice",
			args:    args{v: []byte{}},
			want:    []byte{0xCC, 0x00},
			wantErr: false,
		},
		{
			name:    "8 bit byte slice",
			args:    args{v: []byte{0x61, 0x62, 0x63}},
			want:    []byte{0xCC, 0x03, 0x61, 0x62, 0x63},
			wantErr: false,
		},
		{
			name: "16 bit byte slice",
			args: args{
				v: []byte{
					0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50,
					0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49,
					0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A, 0x41, 0x42,
					0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55,
					0x56, 0x57, 0x58, 0x59, 0x5A, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E,
					0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47,
					0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A,
					0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50, 0x51, 0x52, 0x53,
					0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C,
					0x4D, 0x4E, 0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A, 0x41, 0x42, 0x43, 0x44, 0x45,
					0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58,
					0x59, 0x5A, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50, 0x51,
					0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A,
					0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A,
				},
			},
			want: []byte{
				0xCD, 0x01, 0x04, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50,
				0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49,
				0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A, 0x41, 0x42,
				0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55,
				0x56, 0x57, 0x58, 0x59, 0x5A, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E,
				0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47,
				0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A,
				0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50, 0x51, 0x52, 0x53,
				0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C,
				0x4D, 0x4E, 0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A, 0x41, 0x42, 0x43, 0x44, 0x45,
				0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58,
				0x59, 0x5A, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50, 0x51,
				0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A,
				0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A,
			},
			wantErr: false,
		},
		{
			name:    "empty string",
			args:    args{v: ""},
			want:    []byte{0x80},
			wantErr: false,
		},
		{
			name:    "short string",
			args:    args{v: "abc"},
			want:    []byte{0x83, 0x61, 0x62, 0x63},
			wantErr: false,
		},
		{
			name:    "max length short string",
			args:    args{v: "abcabcabcabcabc"},
			want:    []byte{0x8F, 0x61, 0x62, 0x63, 0x61, 0x62, 0x63, 0x61, 0x62, 0x63, 0x61, 0x62, 0x63, 0x61, 0x62, 0x63},
			wantErr: false,
		},
		{
			name: "8 bit string",
			args: args{v: "ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
			want: []byte{
				0xD0, 0x1A, 0x41, 0x42, 0x43, 0x44,
				0x45, 0x46, 0x47, 0x48, 0x49, 0x4A,
				0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50,
				0x51, 0x52, 0x53, 0x54, 0x55, 0x56,
				0x57, 0x58, 0x59, 0x5A,
			},
			wantErr: false,
		},
		{
			name: "16 bit string",
			args: args{v: "ABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZ"},
			want: []byte{
				0xD1, 0x01, 0x04, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50,
				0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49,
				0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A, 0x41, 0x42,
				0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55,
				0x56, 0x57, 0x58, 0x59, 0x5A, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E,
				0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47,
				0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A,
				0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50, 0x51, 0x52, 0x53,
				0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C,
				0x4D, 0x4E, 0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A, 0x41, 0x42, 0x43, 0x44, 0x45,
				0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58,
				0x59, 0x5A, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50, 0x51,
				0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A,
				0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A,
			},
			wantErr: false,
		},
		{
			name:    "string pointer",
			args:    args{v: strPtr("")},
			want:    []byte{0x80},
			wantErr: false,
		},
		{
			name:    "tiny int min",
			args:    args{v: -16},
			want:    []byte{0xF0},
			wantErr: false,
		},
		{
			name:    "tiny int max",
			args:    args{v: 127},
			want:    []byte{0x7F},
			wantErr: false,
		},
		{
			name:    "negative 8 bit int min",
			args:    args{v: -128},
			want:    []byte{0xC8, 0x80},
			wantErr: false,
		},
		{
			name:    "negative 8 bit int max",
			args:    args{v: -17},
			want:    []byte{0xC8, 0xEF},
			wantErr: false,
		},
		{
			name:    "positive 16 bit int min",
			args:    args{v: 128},
			want:    []byte{0xC9, 0x00, 0x80},
			wantErr: false,
		},
		{
			name:    "positive 16 bit int max",
			args:    args{v: 32_767},
			want:    []byte{0xC9, 0x7F, 0xFF},
			wantErr: false,
		},
		{
			name:    "negative 16 bit int min",
			args:    args{v: -32_768},
			want:    []byte{0xC9, 0x80, 0x00},
			wantErr: false,
		},
		{
			name:    "negative 16 bit int max",
			args:    args{v: -129},
			want:    []byte{0xC9, 0xff, 0x7f},
			wantErr: false,
		},
		{
			name:    "positive 32 bit int min",
			args:    args{v: 32_768},
			want:    []byte{0xCA, 0x00, 0x00, 0x80, 0x00},
			wantErr: false,
		},
		{
			name:    "positive 32 bit int max",
			args:    args{v: 2_147_483_647},
			want:    []byte{0xCA, 0x7F, 0xFF, 0xFF, 0xFF},
			wantErr: false,
		},
		{
			name:    "negative 32 bit int min",
			args:    args{v: -2_147_483_648},
			want:    []byte{0xCA, 0x80, 0x00, 0x00, 0x00},
			wantErr: false,
		},
		{
			name:    "negative 32 bit int max",
			args:    args{v: -32_769},
			want:    []byte{0xCA, 0xFF, 0xFF, 0x7F, 0xFF},
			wantErr: false,
		},
		{
			name:    "positive 64 bit int min",
			args:    args{v: 2_147_483_648},
			want:    []byte{0xCB, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00},
			wantErr: false,
		},
		{
			name:    "positive 64 bit int max",
			args:    args{v: 9_223_372_036_854_775_807},
			want:    []byte{0xCB, 0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
			wantErr: false,
		},
		{
			name:    "negative 64 bit int min",
			args:    args{v: -9_223_372_036_854_775_808},
			want:    []byte{0xCB, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			wantErr: false,
		},
		{
			name:    "negative 64 bit int max",
			args:    args{v: -2_147_483_649},
			want:    []byte{0xCB, 0xFF, 0xFF, 0xFF, 0xFF, 0x7F, 0xFF, 0xFF, 0xFF},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Marshal(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Marshal() = %x, want %x", got, tt.want)
			}
		})
	}
}

func strPtr(v string) *string { return &v }
func boolPtr(v bool) *bool    { return &v }
