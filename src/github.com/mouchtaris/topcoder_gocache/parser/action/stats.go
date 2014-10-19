package action

import (
    "github.com/mouchtaris/topcoder_gocache/parser/lex"
    "github.com/mouchtaris/topcoder_gocache/command"
)

//
type Stats struct {
    consumer chan<- command.Stats
}

func NewStats (consumer chan<- command.Stats) *Stats {
    return &Stats {
        consumer: consumer,
    }
}

//
func (*Stats) Name () string {
    return "stats"
}

//
func (action *Stats) Parse (lex *lex.Lexer) error {
    comm := command.Stats { }

    err := lex.ReadEOC()
    if err != nil {
        return err
    }

    action.consumer <- comm
    return nil
}
