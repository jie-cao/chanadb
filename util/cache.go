package util

import "sync"

type LRUHandle struct {
	value     interface{}
	deleter   *func(slice *Slice, )
	nextHash  *LRUHandle
	next      *LRUHandle
	prev      *LRUHandle
	charge    uint
	keyLength uint
	inCache   bool
	refs      uint32
	hash      uint32
	keyData   []byte
}

func (lruHandle *LRUHandle) Key() *Slice {
	return NewSlice(lruHandle.keyData, lruHandle.keyLength)
}

type HandleTable struct {
	length   uint32
	elements uint32
	list     []*LRUHandle
}

func (handleTable *HandleTable) Insert(h *LRUHandle) *LRUHandle {
	ptr := handleTable.FindPointer(h.Key(), h.hash)
	old := *ptr
	if old == nil {
		h.nextHash = nil
	} else {
		h.nextHash = old.nextHash
	}

	*ptr = h
	if old == nil {
		handleTable.elements++
		if handleTable.elements > handleTable.length {
			handleTable.Resize()
		}
	}

	return old
}

func (handleTable *HandleTable) Remove(key *Slice, hash uint32) *LRUHandle {
	ptr := handleTable.FindPointer(key, hash)
	result := *ptr
	if result != nil {
		*ptr = result.nextHash
		handleTable.elements--
	}

	return result
}

func (handleTable *HandleTable) Resize() {
	var newLength int = 4
	for newLength < int(handleTable.elements) {
		newLength *= 2
	}

	newList := make([]*LRUHandle, newLength)
	count := 0
	for i := 0; i < int(handleTable.length); i++ {
		h := handleTable.list[i]
		for h != nil {
			next := h.nextHash
			hash := h.hash
			ptr := &newList[int(hash)&(newLength-1)]
			h.nextHash = *ptr
			*ptr = h
			h = next
			count++
		}
	}
	handleTable.list = newList
	handleTable.length = uint32(newLength)
}

func (handleTable *HandleTable) FindPointer(key *Slice, hash uint32) **LRUHandle {
	ptr := &(handleTable.list[hash&(handleTable.length)])
	for *ptr != nil && ((*ptr).hash != hash || key != (*ptr).Key()) {
		ptr = &((*ptr).nextHash)
	}
	return ptr
}

type LRUCache struct {
	capacity int
	mutex    sync.Mutex
	usage    int
	lru      LRUHandle
	inUse    LRUHandle
	table    HandleTable
}

func newLRUCache() *LRUCache {
	var lruCache = new(LRUCache)
	lruCache.lru.next = &(lruCache.lru)
	lruCache.lru.prev = &(lruCache.lru)
	lruCache.inUse.next = &(lruCache.inUse)
	lruCache.inUse.prev = &(lruCache.inUse)

	return lruCache
}

func (lruCache *LRUCache) Ref(e *LRUHandle) {
	if e.refs == 1 && e.inCache {
		lruCache.LRURemove(e)
		lruCache.LRUAppend(&(lruCache.inUse), e)
	}
	e.refs++
}

func (lruCache *LRUCache) LRURemove(e *LRUHandle) {

}

func (lruCache *LRUCache) LRUAppend(list *LRUHandle, e *LRUHandle) {

}

func (lruCache *LRUCache) LookUp() *LRUHandle {

}

func (lruCache *LRUCache) Release(cache *LRUCache) {

}
func (lruCache *LRUCache) Insert(key *Slice, hash uint32, value interface{}, charge int, deleter *func(key *Slice, value interface{})) {

}


