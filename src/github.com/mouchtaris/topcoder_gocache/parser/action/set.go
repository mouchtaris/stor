package action

import (
    "github.com/mouchtaris/topcoder_gocache/lex"
    "github.com/mouchtaris/topcoder_gocache/command"
)

//
type Set struct { }

//
func (Set) Name () string {
    return "set"
}

//
func (action Set) Parse (lex *lex.Lexer) (command.Command, error) {
    comm := command.Set { }

    err := lex.ReadKey()
    if err != nil {
        return nil, err
    }
    comm.Key = string(lex.Token())

    err = lex.ReadEOC()
    if err != nil {
        return nil, err
    }

    err = lex.ReadValue()
    if err != nil {
        return nil, err
    }
    comm.Data = string(lex.Token())

    err = lex.ReadEOC()
    if err != nil {
        return nil, err
    }

    return &comm, nil
}
