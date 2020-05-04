// Copyright (c) 2020 The LevelDB Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file. See the AUTHORS file for names of contributors.

package db

const kBlockSize int = 4096

//内存池管理
type Arena struct {
	// 当前内存块未分配内存的起始地址, 用于下次请求内存的时候找到位置
	allocPtr int
	// 当前写入的内存块（block）的index
	curBlock int
	// 当前内存块剩余的内存
	allocBytesRemaining int

	// 总的allocate的内存
	blocks [][]byte

	// 内存统计
	memoryUsage int
}

// 关闭内存池
func (arena *Arena) Close() {
	for i := 0; i < len(arena.blocks); i++ {
		arena.blocks[i] = nil
	}
}

// 申请内存主函数
// 返回对应的可以写入的byte的slice
func (arena *Arena) Allocate(bytes int) []byte {

	// 如果现在的block剩余的bytes可以存放申请的空间
	if bytes <= arena.allocBytesRemaining {

		//存放现在的位置
		current := arena.allocPtr

		//向后移动到下次的位置
		arena.allocPtr += bytes

		//扣除申请的内存
		arena.allocBytesRemaining -= bytes
		return arena.blocks[arena.curBlock][current:]
	}
	return arena.AllocateFallback(bytes)
}

// 申请内存兜底函数
// 返回对应的可以写入的byte的slice
func (arena *Arena) AllocateFallback(bytes int) []byte {
	// 如果所需的空间大于块的1/4，直接申请一个bytes大小的空间
	if bytes > kBlockSize/4 {
		var result = arena.AllocateNewBlock(bytes)
		return result
	}

	// 否则申请一个4K大小的内存空间
	var result = arena.AllocateNewBlock(kBlockSize)
	arena.allocPtr += bytes
	arena.allocBytesRemaining = kBlockSize - bytes
	return result
}

//
func (arena *Arena) AllocateAligned(bytes int) []byte {
	var align = 8
	currentMod := arena.allocPtr & (align - 1)
	slop := align - currentMod
	if currentMod == 0 {
		slop = 0
	}

	needed := bytes + slop
	if needed <= arena.allocBytesRemaining {
		var current = arena.allocPtr + slop
		arena.allocPtr += slop
		arena.allocPtr += needed
		arena.allocBytesRemaining -= needed
		return arena.blocks[arena.curBlock][current:]
	} else {
		return arena.AllocateFallback(bytes)
	}
}

// 申请一块新的内存块，一般出现的情况是现在的内存块不足以放下申请的空间
// 返回新申请的内存空间
func (arena *Arena) AllocateNewBlock(blockBytes int) []byte {
	result := make([]byte, blockBytes, blockBytes)
	arena.blocks = append(arena.blocks, result)
	arena.curBlock += 1
	arena.memoryUsage += blockBytes
	return result
}
