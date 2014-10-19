package action

import (
    "github.com/mouchtaris/topcoder_gocache/parser/lex"
)

//
// A parsing action is an action that handles the parsing of
// a single command.
type Action interface {
    Name () string
    Parse (*lex.Lexer) error
}
