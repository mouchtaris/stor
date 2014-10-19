package lex

import (
    "github.com/mouchtaris/topcoder_gocache/util"
    "io"
    "fmt"
    "errors"
)

const MAX_KEY_SIZE    = 250
const MAX_VALUE_SIZE  = 8 << 10 // 8KiB
const LONGEST_COMMAND = 6 // for "delete"
const MAX_COMMAND_SIZE = LONGEST_COMMAND
const BUFFER_SIZE     = MAX_VALUE_SIZE * 2 // 16KiB

var ErrInputTooLarge = errors.New("input too large")
var ErrLexing = errors.New("lexing error")

// TODO remove
func log (format string, args ... interface { }) {
    var _ = fmt.Println
//  fmt.Printf(format, args...)
//  fmt.Println()
}

//
// The lexer is responsible for parsing raw input
// characters (bytes) into delimited tokens.
//
// The lexer is also resposible for ensuring
// input data size limitations.
//
// It is a very flexible lexer, as it operates
// on Readers.
type Lexer struct {
    r io.Reader
    buf *util.FixedByteBuffer
    length uint32
}

//
// A lexer -- operates simply on a reader.
func NewLexer (r io.Reader) *Lexer {
    return &Lexer {
        r,
        util.NewFixedByteBuffer(make([]byte, BUFFER_SIZE), 0, BUFFER_SIZE, 0, 0),
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
    return false
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
// Properly handle compacting the parser's buffer.
// This is required because the parser is using indefinite
// look-ahead so as to save memory. Thus, when compact()ing
// the buffer, the token being curretly read has to be
// saved from oblivion.
func (lex *Lexer) compact () {
    log("compacting: %s", lex.buf.Stats())
    if lex.length > 0 {
        log("stepbacking length(%d) bytes", lex.length)
        // We're in the middle of parsing a token.
        lex.buf.StepBack(lex.length)
        log("%s", lex.buf.Stats())
    }
    lex.buf.Compact()
    log("compacted: %s", lex.buf.Stats())
}

//
// Make sure that a buffer can always provide.
// Also, temporarily supress EOF when there are still
// buffered bytes to process.
func (lex *Lexer) fillBuffer () error {
    log("filling buffer")
    lex.compact()
    if (lex.buf.Available() == 0) {
        return util.ErrBufferOverflow
    }

    n, err := lex.buf.ReadFrom(lex.r)
    if err == nil || err == io.EOF {
        log("filled: %s", lex.buf.Stats())
        // bring back to "length" bytes read for current token
        lex.buf.Flip()
        log("flipped: %s", lex.buf.Stats())
        lex.buf.Read(make([]byte, lex.length))
        log("read(step forward): %s", lex.buf.Stats())
    }
    if err != nil && (err != io.EOF || n == 0) {
        // only propagate EOF errors if there are
        // no more buffered bytes left
        return err
    }
    return nil
}

//
// Read a byte from the buffer, and fill it if necessary.
// Will never return a meaningful byte and an error together.
func (lex *Lexer) readByte () (byte, error) {
    log("reading byte")
    var err error
    for ; lex.buf.Available() == 0 && err == nil; err = lex.fillBuffer() {
    }
    if err != nil {
        return 0, err
    }
    lex.length++
    return lex.buf.ReadByte()
}

//
// Unread the last read byte, so that it becomes available for
// reading again.
func (lex *Lexer) unreadByte (n uint32) error {
    log("unreading %d bytes... %s", n, lex.buf.Stats())
    err := lex.buf.StepBack(n)
    if err != nil {
        return err
    }
    lex.length -= n
    log("length=%d", lex.length)
    return nil
}

//
// Skips all characters for which the given predicate holds true.
// Does not keep track of anything, neither is any backtracking possible.
// Returns any error that possibly comes up during input reading.
func (lex *Lexer) skip (pred func (byte) bool) error {
    log("skipping, length=%d %s", lex.length, lex.buf.Stats())
    lex.length = 0
    c, err := lex.readByte()
    for pred(c) && err == nil {
        lex.length = 0
        c, err = lex.readByte()
    }
    if err == nil {
        err = lex.unreadByte(1)
    }
    log("skipped, length (should be 0) =%d %s", lex.length, lex.buf.Stats())
    return err
}

//
// Consumes all chars that are of no interest to anyone (like whitespace)
// without keeping track of anything.
func (lex *Lexer) consumeSpace () error {
    return lex.skip(func (c byte) bool { return !isWord(c); })
}

//
// Reads from the stream and marks the buffer so that
// it marks the last read token.
// Whitespace is ignored and the token is formulated
// according to the given predicate for characters.
func (lex *Lexer) readWhile (pred func(byte)bool, maxbytes uint32) error {
    if err := lex.consumeSpace(); err != nil {
        return err
    }

    i := uint32(0)
    c, err := lex.readByte()
    for ; pred(c) && i < maxbytes && err == nil; c, err = lex.readByte() {
        i++
    }
    if err == nil {
        if i == maxbytes {
            lex.unreadByte(i)
            return ErrInputTooLarge
        }
        // !pred(c) here
        lex.unreadByte(1)
        return nil
    }
    // err != nil here
    // pred(c) and i < maxbytes here (because of lex.readByte() semantics)

    // supress EOF "error" if bytes were read, it will reappear in the next call
    if err == io.EOF && i > 0 {
        return nil
    }
    lex.unreadByte(i)
    return err
}

//
// Read the next command token from the input stream.
// This returns any lexical token which could be
// a command according to the grammar.
// If no error is returned, the current
// token can be accessed by Token().
func (lex *Lexer) ReadCommand () error {
    return lex.readWhile(isCommand, MAX_COMMAND_SIZE)
}

//
// Read the next key from the input stream.
// This returns any lexical token which could
// be a key according to the grammar.
// If no error is returned, the current
// token can be accessed by Token().
func (lex *Lexer) ReadKey () error {
    return lex.readWhile(isWord, MAX_KEY_SIZE)
}

//
// Read the next value from the input stream.
// This returns any lexical token which could
// be a value according to the grammar.
// If no error is returned, the current
// token can be accessed by Token().
func (lex *Lexer) ReadValue () error {
    return lex.readWhile(isWord, MAX_VALUE_SIZE)
}

//
// Return the end-of-command sequence,
// whish is \r\n.
// If this is not the next sequence in the input stream,
// ErrLexing is return and the stream remains intact
// (except for consumed whitespace).
func (lex *Lexer) ReadEOC () error {
    rollback := uint32(0)
    defer func () {
        lex.unreadByte(rollback)
    }()

    err := lex.skip(func (c byte) bool {
        return !isWord(c) && c != '\r'
    })
    if err != nil {
        return err
    }

    b, err := lex.readByte()
    if err != nil {
        return err
    }
    rollback++

    if b != '\r' {
        return ErrLexing
    }

    b, err = lex.readByte()
    if err != nil {
        return err
    }
    rollback++

    if b != '\n' {
        return ErrLexing
    }

    return nil
}

//
// Return a slice view of the current
// token read, in the buffer memory.
func (lex *Lexer) Token () []byte {
    tok, err := lex.buf.Snapshot(lex.length)
    if err != nil {
        return nil
    }
    return tok
}
