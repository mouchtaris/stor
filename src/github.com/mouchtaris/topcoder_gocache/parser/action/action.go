package action

import (
    "github.com/mouchtaris/topcoder_gocache/lex"
    "github.com/mouchtaris/topcoder_gocache/command"
)

//
// A parsing action is an action that handles the parsing of
// a single command.
type Action interface {
    Name () string
    Parse (*lex.Lexer) (command.Command, error)
}
