package action

import (
    "github.com/mouchtaris/topcoder_gocache/lex"
    "github.com/mouchtaris/topcoder_gocache/command"
)

//
type Delete struct { }

//
func (Delete) Name () string {
    return "delete"
}

//
func (Delete) Parse (lex *lex.Lexer) (command.Command, error) {
    comm := command.Delete{ }

    err := lex.ReadKey()
    if err != nil {
        return nil, err
    }
    comm.Key = string(lex.Token())

    err = lex.ReadEOC()
    if err != nil {
        return nil, err
    }

    return &comm, nil
}
