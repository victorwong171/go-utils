package trie

import (
	"github.com/go-test/deep"
	"testing"
)

func TestInitTrie(t *testing.T) {
	type args struct {
		cnt       int
		headChar  byte
		charIndex func(byte) int
	}
	tests := []struct {
		name string
		args args
		want *trie
	}{
		{
			name: "all is ok",
			args: args{
				cnt:      26,
				headChar: 'a',
			},
			want: &trie{
				root: &trieNode{
					children: make([]*trieNode, 26),
				},
				headChar: 'a',
				cnt:      26,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InitTrie(tt.args.cnt, tt.args.headChar, tt.args.charIndex)
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("InitTrie() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_trie_defaultIndex(t1 *testing.T) {
	type fields struct {
		root      *trieNode
		headChar  byte
		cnt       int
		charIndex func(byte) int
	}
	type args struct {
		char byte
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "all is ok",
			args: args{
				char: 'b',
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := InitTrie(26, 'a', nil)
			if got := t.defaultIndex(tt.args.char); got != tt.want {
				t1.Errorf("defaultIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_trie_find(t1 *testing.T) {

	type args struct {
		word string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "all is ok",
			args: args{
				word: "hello",
			},
			want: true,
		},
		{
			name: "bad word",
			args: args{
				word: "1",
			},
			want: false,
		},
		{
			name: "not in",
			args: args{
				word: "world",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := InitTrie(26, 'a', nil)
			_ = t.insert("hello")
			if got := t.find(tt.args.word); got != tt.want {
				t1.Errorf("find() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_trie_insert(t1 *testing.T) {
	type args struct {
		word string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "all is ok",
			args: args{
				word: "hello",
			},
			wantErr: false,
		},
		{
			name: "bad word",
			args: args{
				word: "1",
			},
			wantErr: true,
		},
		{
			name: "already in",
			args: args{
				word: "hello",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := InitTrie(26, 'a', nil)
			if err := t.insert(tt.args.word); (err != nil) != tt.wantErr {
				t1.Errorf("insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
