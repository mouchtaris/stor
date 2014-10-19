package action

import (
    "github.com/mouchtaris/topcoder_gocache/lex"
    "github.com/mouchtaris/topcoder_gocache/command"
    "errors"
)

var ErrQuit = errors.New("quiting")

//
type Quit struct {
    consumer chan<- command.Command
}

func NewQuit (consumer chan<- command.Command) *Quit {
    return &Quit {
        consumer: consumer,
    }
}

//
func (*Quit) Name () string {
    return "quit"
}

//
func (action *Quit) Parse (lex *lex.Lexer) error {
    err := lex.ReadEOC()
    if err != nil {
        return err
    }

    return ErrQuit
}
