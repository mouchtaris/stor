package util

import (
    "errors"
)

var ErrBufferUnderflow = errors.New("buffer underflow")

//
// A fixed byte buffer wraps an existing
// allocation of a byte array.
//
// It offers convenience methods about manipulating
// such a fixed-size buffer.
//

type FixedByteBuffer struct {
    mem []byte
    pos, limit uint32
}

func NewFixedByteBuffer (mem []byte, start, end, pos, limit uint32) {
    return &FixedByteBuffer { mem[start:end], pos, limit }
}

func (buf *FixedByteBuffer) Bytes () []byte {
    return buf.mem
}

func (buf *FixedByteBuffer) ReadByte () (byte, error) {
    if buf.pos == buf.limit {
        return 0, ErrBufferUnderflow
    }
    result := buf.mem[buf.pos]
    buf.pos++
    return result, nil
}

func (buf *FixedByteBuffer) Pack () {
    copy(buf.mem[buf.pos:buf.limit], buf.mem)
    buf.limit -= buf.pos
    buf.pos = 0
}

func (buf *FixedByteBuffer) Flip () {
    buf.limit = buf.pos
    buf.pos = 0
}

func (buf *FixedByteBuffer) Clear () {
    buf.pos = 0
    buf.limit = len(buf.mem)
}
