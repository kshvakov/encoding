package encoding

import (
	"io"
	"sync"
)

// Recycle column buffers, preallocate column buffers
var chunkPool = sync.Pool{}

func newBuffer(initSize int) *buffer {
	b := &buffer{}
	b.addChunk(0, initSize)
	return b
}

type buffer struct{ chunks [][]byte }

func (b *buffer) Write(data []byte) (int, error) {
	var (
		chunkIdx = len(b.chunks) - 1
		dataSize = len(data)
	)
	for {
		freeSize := cap(b.chunks[chunkIdx]) - len(b.chunks[chunkIdx])
		if freeSize >= len(data) {
			b.chunks[chunkIdx] = append(b.chunks[chunkIdx], data...)
			return dataSize, nil
		}
		b.chunks[chunkIdx] = append(b.chunks[chunkIdx], data[:freeSize]...)
		data = data[freeSize:]
		b.addChunk(0, b.calcCap(len(data)))
		chunkIdx++
	}
}

func (b *buffer) alloc(size int) []byte {
	var (
		chunkIdx = len(b.chunks) - 1
		chunkLen = len(b.chunks[chunkIdx])
	)
	if (cap(b.chunks[chunkIdx]) - chunkLen) < size {
		b.addChunk(size, b.calcCap(size))
		return b.chunks[chunkIdx+1]
	}
	b.chunks[chunkIdx] = b.chunks[chunkIdx][:chunkLen+size]
	return b.chunks[chunkIdx][chunkLen : chunkLen+size]
}

func (b *buffer) addChunk(size, capacity int) {
	var chunk []byte
	if c, ok := chunkPool.Get().([]byte); ok && cap(c) >= size {
		chunk = c[:size]
	} else {
		chunk = make([]byte, size, capacity)
	}
	b.chunks = append(b.chunks, chunk)
}

func (b *buffer) writeTo(w io.Writer) error {
	for _, chunk := range b.chunks {
		if _, err := w.Write(chunk); err != nil {
			b.free()
			return err
		}
	}
	b.free()
	return nil
}

func (b *buffer) bytes() []byte {
	if len(b.chunks) == 1 {
		return b.chunks[0]
	}
	bytes := make([]byte, 0, b.len())
	for _, chunk := range b.chunks {
		bytes = append(bytes, chunk...)
	}
	return bytes
}

func (b *buffer) len() int {
	var v int
	for _, chunk := range b.chunks {
		v += len(chunk)
	}
	return v
}

func (b *buffer) calcCap(dataSize int) int {
	dataSize = max(dataSize, 64)
	if len(b.chunks) == 0 {
		return dataSize
	}
	// Always double the size of the last chunk
	return max(dataSize, cap(b.chunks[len(b.chunks)-1])*2)
}

func (b *buffer) free() {
	if len(b.chunks) == 0 {
		return
	}
	// Recycle all chunks except the last one
	chunkSizeThreshold := cap(b.chunks[0])
	for _, chunk := range b.chunks[:len(b.chunks)-1] {
		// Drain chunks smaller than the initial size
		if cap(chunk) >= chunkSizeThreshold {
			chunkPool.Put(chunk[:0])
		} else {
			chunkSizeThreshold = cap(chunk)
		}
	}
	// Keep the largest chunk
	b.chunks[0] = b.chunks[len(b.chunks)-1][:0]
	b.chunks = b.chunks[:1]
}

func max(a, b int) int {
	if b > a {
		return b
	}
	return a
}
