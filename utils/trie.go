package utils

import "strings"

type Trie struct {
  Root *TrieNode
  Itr *TrieNode
}

type TrieNode struct {
  char *rune
  End bool
  parent *TrieNode
  children map[rune]*TrieNode
}

func NewTrie() Trie {
  root := NewTrieNode(nil, false, nil)
  return Trie{
    Root: root,
    Itr: root,
  }
}

func NewTrieNode(r *rune, e bool, p *TrieNode) *TrieNode {
  return &TrieNode {
    char: r,
    End: e,
    parent: p,
    children: make(map[rune]*TrieNode),
  }
}

// BuildTrie returns a trie populated with a list of strings
func BuildTrie(strings []string) Trie{
  trie := NewTrie()
  for _, s := range strings {
    trie.add(s) 
  }
  return trie
}

// add adds the provided string to the trie
func (t Trie) add(s string) { 
  // add to children if not there already
  itr := t.Root
  for _, char := range s {
    if node, exists := itr.children[char]; exists {
      itr = node
    } else {
      charCopy := char
      newNode := NewTrieNode(&charCopy, false, itr)
      itr.children[char] = newNode
      itr = newNode
    }
  }
  itr.End = true
}

// itrStep attempts to step down the trie to the provided rune.
// returns true & steps itr if possible, otherwise returns false and resets itr to root
func (t *Trie) ItrStep(r rune) bool {
  // if none of the children are the rune, return false & reset itr + path
  if node, exists := t.Itr.children[r]; exists {
    t.Itr = node
    return true;
  }
  t.Itr = t.Root 
  return false
}

func (t Trie) GetPath() string {
  itr := t.Itr
  var builder strings.Builder

  for ; itr.parent != nil; itr = itr.parent{
    builder.WriteRune(*itr.char)
  }

  return ReverseString(builder.String())
 }

