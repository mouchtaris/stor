package action

import (
    "github.com/mouchtaris/topcoder_gocache/parser/lex"
    "github.com/mouchtaris/topcoder_gocache/command"
)

//
type Get struct {
    consumer chan<- command.Get
}

func NewGet (consumer chan<- command.Get) *Get {
    return &Get {
        consumer: consumer,
    }
}

//
func (*Get) Name () string {
    return "set"
}

//
func (comm *Get) Parse (lex *lex.Lexer) error {
    panic("TODO implement me")
    return nil
}
