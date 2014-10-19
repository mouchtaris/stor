package action

import (
    "topcoder.com/mouchtaris/scs/lex"
    "topcoder.com/mouchtaris/scs/command"
)

//
// A parsing action is an action that handles the parsing of
// a single command.
type Action interface {
    Name () string
    Parse (*lex.Lexer) (command.Command, error)
}
