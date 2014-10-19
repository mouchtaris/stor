package action

import (
    "github.com/mouchtaris/topcoder_gocache/parser/lex"
    "github.com/mouchtaris/topcoder_gocache/command"
)

//
type Delete struct {
    consumer chan<- command.Command
}

func NewDelete (consumer chan<- command.Command) *Delete {
    return &Delete {
        consumer: consumer,
    }
}

//
func (*Delete) Name () string {
    return "delete"
}

//
func (action *Delete) Parse (lex *lex.Lexer) error {
    comm := command.Delete{ }

    err := lex.ReadKey()
    if err != nil {
        return err
    }
    comm.Key = string(lex.Token())

    err = lex.ReadEOC()
    if err != nil {
        return err
    }

    action.consumer <- &comm
    return nil
}
