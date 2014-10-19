package action

import (
    "github.com/mouchtaris/topcoder_gocache/parser/lex"
    "github.com/mouchtaris/topcoder_gocache/command"
)

//
type Set struct {
    consumer chan<- command.Set
}

func NewSet (consumer chan<- command.Set) *Set {
    return &Set {
        consumer: consumer,
    }
}

//
func (*Set) Name () string {
    return "set"
}

//
func (action *Set) Parse (lex *lex.Lexer) error {
    comm := command.Set { }

    err := lex.ReadKey()
    if err != nil {
        return err
    }
    comm.Key = string(lex.Token())

    err = lex.ReadEOC()
    if err != nil {
        return err
    }

    err = lex.ReadValue()
    if err != nil {
        return err
    }
    comm.Data = string(lex.Token())

    err = lex.ReadEOC()
    if err != nil {
        return err
    }

    action.consumer <- comm
    return nil
}
