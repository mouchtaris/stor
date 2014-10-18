package topcoder_gocache

import (
    "github.com/mouchtaris/topcoder_gocache/util"
    "io"
    "bytes"
)

const MAX_KEY_SIZE    = 250
const MAX_VALUE_SIZE  = 8 << 10 // 8KiB
const LONGEST_COMMAND = 6 // for "delete"
const BUFFER_SIZE     = MAX_VALUE_SIZE * 2 // 16KiB
//
// The parser is responsible for parsing raw input
// characters (bytes) into meaningful tokens and
// commands.
//
// The parser is also resposible for ensuring
// input data size limitations.
//
type Parser struct {
    r io.Reader
    buf util.FixedByteBuffer
    err error
    length uint32
}

func NewParser (r io.Reader) *Parser {
    return &Parser {
        r,
        util.NewFixedByteBuffer([BUFFER_SIZE]byte, 0, BUFFER_SIZE, 0, 0),
        nil,
        uint32(0),
    }
}

// Check whether an input character is a command character.
// NOTE: allowing both upper and lower case letters in commands.
func isCommand (c byte) bool {
    switch {
    case 'a' <= c && c <= 'z',
         'A' <= c && c <= 'Z':
        return true
    }
}

//
// Check whether an input character is a "word" character.
func isWord (c byte) bool {
    // Assuming that switch is faster than bytes.Contains()
    switch {
    case isCommand(c),
         '0' <= c && c <= '9':
        return true
    }
    switch c {
    case '!', '#', '$', '%', '&', '\'', '"', '*', '+', '-',
        '/', '=', '?', '^', '_', '`', '{', '|', '}', '~',
        '(', ')', '<', '>', '[', ']', ':', ';', '@', ',', '.':
        return true
    }
    return false
}

//
// Make sure that a buffer can always provide.
// Also, temporarily supress EOF when there are still
// buffered bytes to process.
func (lex *Parser) fillBuffer {
    if lex.err != nil {
        return
    }

    lex.buf.Compact()
    if (lex.buf.Available() == 0) {
        lex.err = util.ErrBufferOverflow
        return
    }

    n, err := lex.buf.ReadFrom(lex.r)
    if err != nil && (err != io.EOF || n == 0) {
        // only propagate EOF errors if there are
        // no more buffered bytes left
        lex.err = err
    }
    lex.buf.Flip()
}

//
// Read a byte from the buffer, and fill it if necessary.
func (lex *Parser) readByte () byte {
    for lex.buf.Available() == 0 && lex.err == nil {
        lex.fillBuffer()
    }
    if lex.err != nil {
        return 0
    }
    return lex.buf.ReadByte()
}

//
// Reads from the stream and marks the buffer so that
// it marks the last read token.
// Whitespace is ignored and the token is formulated
// according to the given predicate for characters.
func (lex *Parser) readWhile (pred func(byte)bool) {
    if lex.err != nil {
        return
    }

    i := uint32(0)
    c := lex.ReadByte()
    for ; pred(c) && lex.err == nil; c = lex.ReadByte() {
        ++i
    }
    // supress EOF "error" if bytes were read, it will reappear in the next call
    if lex.err == nil {
        lex.buf.StepBack(1)
        lex.length = i - 1
        return
    }
    if lex.err == io.EOF && i > 0
        lex.length = i
        return
    }
    lex.buf.StepBack(i)
}

//
// Read a command token.
func (lex *Parser) readCommand () {
}

//    for lex.err == nil && lex.buf.Available() < LONGEST_COMMAND {
//        lex.fillBuffer()
//    }
//    if (lex.err != nil) {
//        return nil, lex.err
//    }
//    i := uint(0)
//    var c byte
//    for (c, err := lex.buf.ReadByte(); err == nil && isCommand(c); c, err = lex.buf.ReadByte()) {
//        ++i
//    }
//    if (!isCommand(c)) {
//        lex.buf.StepBack(1)
//    }
//    return lex.buf.StringSnapshot(i), nil
//}

//
// Read the next command from the input stream.
// This returns any lexical token which could be
// a command according to the grammar.
func (lex *Parser) readCommand (string, error) {
    for lex.err == nil && lex.buf.Available() < LONGEST_COMMAND {
        lex.fillBuffer()
    }
    if (lex.err != nil) {
        return nil, lex.err
    }
    i := uint32(0)
    var c byte
    for (c, err := lex.buf.ReadByte(); err == nil && isCommand(c); c, err = lex.buf.ReadByte()) {
        ++i
    }
    if (!isCommand(c)) {
        lex.buf.StepBack(1)
    }
    return lex.buf.StringSnapshot(i), nil
}

// Read the next key from the input stream.
// This returns any lexical token which could
// be a key according to the grammer.

//
// Read the next 
func (lex *Parser) NextToken () (lextoken.Token, error) {
    if lex.err != nil {
        return 0, lex.err
    }

    lex.fillBuffer()
    if lex.err != nil {
        lex.buf.Flip()
        word, lex.err := lex.readWord()
        if lex.err != nil {
            return 0, lex.err
        }
        token, lex.err = parseWord(word)
    }
}

