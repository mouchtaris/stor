package parser

import (
    "github.com/mouchtaris/topcoder_gocache/parser/lextoken"
    "io"
)

//
// The lexer is responsible for parsing raw input
// characters (bytes) into meaningful tokens.
//
// The lexer is also resposible for ensuring
// input data size limitations.
//
func NextToken (r io.Reader) lextoken.Token {
    return lextoken.SET
}

