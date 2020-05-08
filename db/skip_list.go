package db

import (
	"math/rand"
	"strings"
	"sync"
	"time"
)

const (
	kMaxHeight      int = 12
)

type SkipList struct {
	node Node
	m *sync.RWMutex
	arena *Arena
	head *Node
	maxHeight int
	rnd rand.Rand
}

type Iterator struct {
	list *SkipList
	node *Node
}


func newIterator(list *SkipList) *Iterator{
	iterator := new(Iterator)
	iterator.list = list
	iterator.node = nil
	return iterator
}

// Returns true iff the iterator is positioned at a valid node.
func (iterator *Iterator) Valid() bool {
	return iterator.node != nil
}

// Returns the key at the current position.
// REQUIRES: Valid()
func (iterator *Iterator) Key() string {
	return iterator.node.key
}

// Advances to the next position.
// REQUIRES: Valid()
func (iterator *Iterator) Next() {
	iterator.list.m.RLock()
	iterator.node = iterator.node.Next(0)
	iterator.list.m.RUnlock()
}

// Advances to the next position.
// REQUIRES: Valid()
func (iterator *Iterator) Prev() {
	iterator.node = iterator.list.FindLessThan(iterator.node.key)
	if iterator.node == iterator.list.head {
		iterator.node = nil
	}
}

// Advance to the first entry with a key >= target
func (iterator *Iterator) Seek(key string) {
	iterator.node = iterator.list.FindGreaterOrEqual(key, nil)
}

// Position at the first entry in list.
// Final state of iterator is Valid() iff list is not empty.
func (iterator *Iterator) SeekToFirst() {
	iterator.list.m.RLock()
	iterator.node = iterator.list.head.Next(0)
	iterator.list.m.RUnlock()
}

// Position at the last entry in list.
// Final state of iterator is Valid() iff list is not empty.
func (iterator *Iterator) SeekToLast() {
	iterator.node = iterator.list.FindLast()
	if iterator.node == iterator.list.head {
		iterator.node = nil
	}
}


func (s *SkipList) Insert(key string)  {
	// TODO(opt): We can use a barrier-free variant of FindGreaterOrEqual()
	// here since Insert() is externally synchronized.
	prev := make([]*Node, kMaxHeight)
	x := s.FindGreaterOrEqual(key, prev)

	height := s.RandomHeight();
	if height > s.GetMaxHeight() {
		for i := s.GetMaxHeight(); i < height; i++ {
			prev[i] = s.head;
		}
		// It is ok to mutate max_height_ without any synchronization
		// with concurrent readers.  A concurrent reader that observes
		// the new value of max_height_ will see either the old value of
		// new level pointers from head_ (nullptr), or a new value set in
		// the loop below.  In the former case the reader will
		// immediately drop to the next level since nullptr sorts after all
		// keys.  In the latter case the reader will use the new node.
		s.maxHeight = height
	}

	x = s.NewNode(key, height)
	for i := 0; i < height; i++ {
		// NoBarrier_SetNext() suffices since we will add a barrier when
		// we publish a pointer to "x" in prev[i].
		x.SetNext(i, prev[i].Next(i))
		prev[i].SetNext(i, x)
	}
}

func (s *SkipList) RandomHeight() int {
	// Increase height with probability 1 in kBranching
	const kBranching = 4
	height := 1
	for height < kMaxHeight && (s.rnd.Int() % kBranching) == 0 {
		height++
	}
	return height
}

func (s *SkipList) GetMaxHeight() int {
	return s.maxHeight
}

func (s *SkipList) FindLessThan(key string) * Node {
	x := s.head
	level := s.GetMaxHeight() - 1
	for true {
		s.m.RLock()
		next := x.Next(level)
		s.m.RUnlock()
		if next == nil || strings.Compare(next.key,key) >= 0{
			if level == 0 {
				return x
			} else {
				level--
			}
		} else {
			x = next
		}
	}
	return nil
}


func (s *SkipList) FindLast() * Node {
	x := s.head
	level := s.GetMaxHeight() - 1
	for true {
		s.m.RLock()
		next := x.Next(level)
		s.m.RUnlock()
		if next == nil {
			if level == 0 {
				return x
			} else {
				// Switch to next list
				level--
			}
		} else {
			x = next
		}
	}
	return nil
}


func (s *SkipList) KeyIsAfterNode(key string, n *Node) bool {
	// null n is considered infinite
	return (n != nil) && (strings.Compare(n.key, key) < 0)
}

func (s *SkipList) FindGreaterOrEqual(key string, prev []*Node) *Node {
	x := s.head
	level := s.GetMaxHeight() - 1
	for true {
		next := x.Next(level)
		if s.KeyIsAfterNode(key, next) {
			// Keep searching in this list
			x = next
		} else {
			if prev != nil {
				prev[level] = x
			}
		}
		if level == 0 {
			return next
		} else {
			// Switch to next list
			level--
		}
	}
	return nil
}

func newSkipList() *SkipList {
	s := new(SkipList)
	s.head = s.NewNode("" /* any key will do */, kMaxHeight)
	s.maxHeight = 1
	s.rnd = *rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < kMaxHeight; i++ {
		s.head.SetNext(i, nil)
	}

	return s
}

type Node struct {
	key string
	next []*Node
}

func (node *Node) Next(n int) *Node {
	return node.next[n]
}

func (node *Node) SetNext(n int, x *Node) {
	node.next[n] = x
}

func (s *SkipList) NewNode(key string, height int) *Node {
	newNode := new(Node)
	newNode.key = key
	newNode.next = make([]*Node , height)
	return newNode
}

func (s *SkipList) Contains(key string) bool {
	x := s.FindGreaterOrEqual(key, nil)
	if x != nil && strings.EqualFold(key, x.key) {
		return true
	} else {
		return false
	}
}

