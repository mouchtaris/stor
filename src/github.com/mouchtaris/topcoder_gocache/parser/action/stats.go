package action

import (
    "github.com/mouchtaris/topcoder_gocache/lex"
    "github.com/mouchtaris/topcoder_gocache/command"
)

//
type Stats struct { }

//
func (Stats) Name () string {
    return "stats"
}

//
func (Stats) Parse (lex *lex.Lexer) (command.Command, error) {
    comm := command.Stats { }

    err := lex.ReadEOC()
    if err != nil {
        return nil, err
    }

    return &comm, nil
}
