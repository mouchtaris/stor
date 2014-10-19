package action

import (
    "github.com/mouchtaris/topcoder_gocache/lex"
    "github.com/mouchtaris/topcoder_gocache/command"
)

//
type Quit struct { }

//
func (Quit) Name () string {
    return "quit"
}

//
func (Quit) Parse (lex *lex.Lexer) (command.Command, error) {
    comm := command.Quit { }

    err := lex.ReadEOC()
    if err != nil {
        return nil, err
    }

    return &comm, nil
}
