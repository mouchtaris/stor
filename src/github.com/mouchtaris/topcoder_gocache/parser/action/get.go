package action

import (
    "github.com/mouchtaris/topcoder_gocache/parser/lex"
    "github.com/mouchtaris/topcoder_gocache/command"
)

//
type Get struct {
    consumer chan<- command.Command
}

func NewGet (consumer chan<- command.Command) *Get {
    return &Get {
        consumer: consumer,
    }
}

//
func (*Get) Name () string {
    return "get"
}

//
// Utility function for appending keys to a Get command
// object.
func appendKey (comm *command.Get, lexer* lex.Lexer) {
    comm.Keys = append(comm.Keys, string(lexer.Token()))
}

//
func (action *Get) Parse (lexer *lex.Lexer) error {
    comm := command.Get { Keys: make([]string, 0, 20) }

    err := lexer.ReadKey()
    if err != nil {
        return err
    }
    appendKey(&comm, lexer)

    for err = lexer.ReadEOC(); err == lex.ErrLexing; err = lexer.ReadEOC() {
        err = lexer.ReadKey()
        if err != nil {
            return err
        }
        appendKey(&comm, lexer)
    }
    if err != nil {
        return err
    }

    action.consumer <- &comm
    return nil
}
