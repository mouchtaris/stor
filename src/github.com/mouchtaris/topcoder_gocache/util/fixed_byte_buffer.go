package util

import (
    "errors"
    "io"
)

var ErrBufferUnderflow = errors.New("buffer underflow")

func Min (a, b int) int {
    return a < b? a : b
}

//
// A fixed byte buffer wraps an existing
// allocation of a byte array.
//
// It offers convenience methods about manipulating
// such a fixed-size buffer.
//
// This class copies the concept of a java.nio.ByteBuffer,
// so the documentation of that class is a better
// explaination for what is going on here
// (especially the meaning of start, end, limit and pos
// members/variables/values).
//

type FixedByteBuffer struct {
    mem []byte
    pos, limit uint32
}

func NewFixedByteBuffer (mem []byte, start, end, pos, limit uint32) *FixedByteBuffer {
    return &FixedByteBuffer { mem[start:end], pos, limit }
}

func NewFixedByteBuffer (mem []byte) *FixedByteBuffer {
    return NewFixedByteBuffer(mem, 0, len(mem), 0, len(mem))
}

//
// freeBytes() returns a slice view which includes only
// the available space of the buffer, ie, buf[pos:limit].
func (buf *FixedByteBuffer) freeBytes () {
    return buf.mem[buf.pos:buf.limit]
}

//
// Copies bytes from the Reader stream into the underlying buffer.
// Returns the number of bytes actually written and any error
// the reader provides after reading.
func (buf *FixedByteBuffer) ReadFrom (r io.Reader) (bytesWritten int, err error) {
    bytesWritten, err = r.Read(buf.freeBytes())
    if (bytesWritten > 0) {
        buf.pos += bytesWritten
    }
    return
}

//
// Read a single byte. Returns io.EOF error if no more
// bytes are available.
func (buf *FixedByteBuffer) ReadByte () (byte, error) {
    if buf.Available() == 0 {
        return 0, io.EOF
    }
    result := buf.mem[buf.pos]
    ++buf.pos
    return result, nil
}

//
// StepBack moves position n steps back. Unless there
// has been a modifying operation, this method will
// cause ReadByte() to return the same values as the
// n previous times it was called.
// If there is no room to step back, this method returns
// ErrBufferOverflow as an error.
func (buf *FixedByteBuffer) StepBack (n uint32) error {
    if buf.pos < n {
        return ErrBufferOverflow
    }
    buf.pos -= n
    return nil
}

//
// Construct and return a string which includes bytes
// [pos - n, pos). If n is too crazy, the underlying
// memory slice will panic.
func (buf *FixedByteBuffer) StringSnapshot (n uint32) string {
    return string(buf.mem[buf.pos - n, buf.pos])
}

//
// Transfers all "used" buffer space to the beginning
// of the buffer, so that all free space is continuous.
// Also, mark available space as consumed and unavailable
// space as available (resume previous operation).
func (buf *FixedByteBuffer) Compact () {
    if buf.pos > 0 {
        copy(buf.mem[buf.pos:buf.limit], buf.mem)
    }
    buf.pos = buf.limit - buf.pos
    buf.limit = len(buf.mem)
}

//
// Mark used space as available, so that after a write
// written bytes become available for reading.
func (buf *FixedByteBuffer) Flip () {
    buf.limit = buf.pos
    buf.pos = 0
}

//
// Reset to the initial state (the whole space becomes
// available).
func (buf *FixedByteBuffer) Clear () {
    buf.pos = 0
    buf.limit = len(buf.mem)
}

//
// Return the number of available bytes (for reading
// or writing).
func (buf *FixedByteBuffer) Available () uint32 {
    return buf.limit - buf.pos
}
