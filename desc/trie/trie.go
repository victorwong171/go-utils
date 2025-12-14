package trie

import (
	"gopkg.in/errgo.v2/errors"
)

type trieNode struct {
	children  []*trieNode
	isWordEnd bool
}

type trie struct {
	root      *trieNode
	headChar  byte
	cnt       int
	charIndex func(byte) int
}

var (
	ErrInvalidChar = errors.New("invalid char")
)

func InitTrie(cnt int, headChar byte, charIndex func(byte) int) *trie {
	t := &trie{
		root: &trieNode{
			children: make([]*trieNode, cnt, cnt),
		},
		headChar:  headChar,
		cnt:       cnt,
		charIndex: charIndex,
	}
	if t.charIndex == nil {
		t.charIndex = t.defaultIndex
	}
	return t
}

func (t *trie) insert(word string) error {
	wordLength := len(word)
	current := t.root
	for i := 0; i < wordLength; i++ {
		if index := t.charIndex(word[i]); index != -1 {
			if current.children[index] == nil {
				current.children[index] = &trieNode{
					children: make([]*trieNode, t.cnt, t.cnt),
				}
			}
			current = current.children[index]
		} else {
			return ErrInvalidChar
		}
	}
	current.isWordEnd = true
	return nil
}

func (t *trie) defaultIndex(char byte) int {
	idx := int(char - t.headChar)
	// confirm char must be in range
	if idx >= 0 && idx < cap(t.root.children) {
		return idx
	}
	return -1
}

func (t *trie) find(word string) bool {
	wordLength := len(word)
	current := t.root
	for i := 0; i < wordLength; i++ {
		if index := t.charIndex(word[i]); index != -1 {
			if current.children[index] == nil {
				return false
			}
			current = current.children[index]
		} else {
			return false
		}
	}
	return current.isWordEnd
}
