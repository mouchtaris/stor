package action

import (
    "github.com/mouchtaris/topcoder_gocache/lex"
    "github.com/mouchtaris/topcoder_gocache/command"
)

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
    comm := command.Quit { }

    err := lex.ReadEOC()
    if err != nil {
        return err
    }

    action.consumer <- &comm
    return nil
}
