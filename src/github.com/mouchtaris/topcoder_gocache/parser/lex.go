package parser

import (
    "github.com/mouchtaris/topcoder_gocache/parser/lextoken"
    "io"
    "bytes"
)

const MAX_KEY_SIZE = 250
const MAX_VALUE_SIZE = 8 << 10 // 8KiB

//
// The lexer is responsible for parsing raw input
// characters (bytes) into meaningful tokens.
//
// The lexer is also resposible for ensuring
// input data size limitations.
//
type Lexer struct {
    r io.Reader
    buf [MAX_VALUE_SIZE]byte
}

//
// Check whether an input character is a "word" character.
func isWord (c byte) bool {
    switch {
    case 'a' <= c && c <= 'z',
         'A' <= c && c <= 'Z',
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
func (lex *Lexer) fillBuffer {
    lex.r.Read(lex.buf)
}

func NextToken (r io.Reader) lextoken.Token {
    c := [1]byte
    bytenum, err := 0, nil
    keepLooping := true
    for keepLooping {
    }
    return lextoken.SET
}

