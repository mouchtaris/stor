package action

import (
    "topcoder.com/mouchtaris/scs/lex"
    "topcoder.com/mouchtaris/scs/command"
)

//
type Get struct { }

//
func (Get) Name () string {
    return "get"
}

//
// Utility function for appending keys to a Get command
// object.
func appendKey (comm *command.Get, lexer* lex.Lexer) {
    comm.Keys = append(comm.Keys, string(lexer.Token()))
}

//
func (Get) Parse (lexer *lex.Lexer) (command.Command, error) {
    comm := command.Get { Keys: make([]string, 0, 20) }

    err := lexer.ReadKey()
    if err != nil {
        return nil, err
    }
    appendKey(&comm, lexer)

    for err = lexer.ReadEOC(); err == lex.ErrLexing; err = lexer.ReadEOC() {
        err = lexer.ReadKey()
        if err != nil {
            return nil, err
        }
        appendKey(&comm, lexer)
    }
    if err != nil {
        return nil, err
    }

    return &comm, nil
}
