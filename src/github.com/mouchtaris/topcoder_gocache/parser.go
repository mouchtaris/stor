package topcoder_gocache

import (
    "github.com/mouchtaris/topcoder_gocache/util"
    "io"
    Log "log"
)

const MAX_KEY_SIZE    = 250
const MAX_VALUE_SIZE  = 8 << 10 // 8KiB
const LONGEST_COMMAND = 6 // for "delete"
const BUFFER_SIZE     = MAX_VALUE_SIZE * 2 // 16KiB

func log (fmt string, args ... interface { }) {
    Log.Printf(fmt, args...)
}

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
    buf *util.FixedByteBuffer
    length uint32
}

func NewParser (r io.Reader) *Parser {
    return &Parser {
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
func (lex *Parser) compact () {
    log("compacting: %s", lex.buf.Stats())
    if lex.length > 0 {
        log("stepbacking length bytes")
        // We're in the middle of parsing a token.
        lex.buf.StepBack(lex.length)
    }
    lex.buf.Compact()
    log("compacted: %s", lex.buf.Stats())
}

//
// Make sure that a buffer can always provide.
// Also, temporarily supress EOF when there are still
// buffered bytes to process.
func (lex *Parser) fillBuffer () error {
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
func (lex *Parser) readByte () (byte, error) {
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
func (lex *Parser) unreadByte (n uint32) error {
    log("unreading %d bytes...", n)
    err := lex.buf.StepBack(n)
    if err != nil {
        return err
    }
    lex.length -= n
    log("length=%d", lex.length)
    return nil
}

//
// Consumes all chars that are of no interest to anyone (like whitespace)
// without keeping track of anything.
func (lex *Parser) consumeSpace () error {
    log("consuming space, length=%d", lex.length)
    c, err := lex.readByte()
    for ; !isWord(c) && err == nil; c, err = lex.readByte() {
    }
    if err == nil {
        err = lex.unreadByte(1)
    }
    log("space consumed, length (before setting to 0) =%d", lex.length)
    lex.length = 0
    return err
}

//
// Reads from the stream and marks the buffer so that
// it marks the last read token.
// Whitespace is ignored and the token is formulated
// according to the given predicate for characters.
func (lex *Parser) readWhile (pred func(byte)bool) error {
    if err := lex.consumeSpace(); err != nil {
        return err
    }

    i := uint32(0)
    c, err := lex.readByte()
    for ; pred(c) && err == nil; c, err = lex.readByte() {
        i++
    }
    if err == nil {
        lex.unreadByte(1)
        return nil
    }
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
func (lex *Parser) readCommand () error {
    return lex.readWhile(isCommand)
}

//
// Read the next key from the input stream.
// This returns any lexical token which could
// be a key according to the grammer.
func (lex *Parser) readKey () error {
    return lex.readWhile(isWord)
}

//
// Return a slice view of the current
// token read, in the buffer memory.
func (lex *Parser) Token () []byte {
    return lex.buf.Snapshot(lex.length)
}

//
// Read the next token from the stream.
// If no error is returned, the token's textual
// value can be retrieved by calling Token().
func (lex *Parser) Next () error {
    err := lex.readCommand()
    if err != nil {
        return err
    }
    return nil
}

//
// Next() and Token() together
func (lex *Parser) NextToken() ([]byte, error) {
    err := lex.Next()
    return lex.Token(), err
}
