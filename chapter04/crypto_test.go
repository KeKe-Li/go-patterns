package chapter04

import (
	"reflect"
	"testing"
)

func TestNewCrypto(t *testing.T) {
	type args struct {
		encodingAESKey string
		encryptMsg     string
	}
	tests := []struct {
		name string
		args args
		want ICrypto
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCrypto(tt.args.encodingAESKey, tt.args.encryptMsg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCrypto() = %v, want %v", got, tt.want)
			}
		})
	}
}
