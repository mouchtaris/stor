package action

import (
    "topcoder.com/mouchtaris/scs/lex"
    "topcoder.com/mouchtaris/scs/command"
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
