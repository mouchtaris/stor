package parser

import (
    "github.com/mouchtaris/topcoder_gocache/parser/action"
    "github.com/mouchtaris/topcoder_gocache/lex"
    "github.com/mouchtaris/topcoder_gocache/command"
    "errors"
)

type Parser struct {
    lex* lex.Lexer
    handlers map[string] action.Action
}

var ErrHandlerReregistered = errors.New("handler re-registered")
var ErrNoHandler           = errors.New("no handler registered for this command")

//
func NewParser (lex* lex.Lexer) *Parser {
    return &Parser {
        lex: lex,
        handlers: map[string] action.Action { },
    }
}

//
// Register a handler for handling handler.Name() commands.
// If a handler for that name/command already exists,
// ErrHandlerReregistered is returned.
func (yy *Parser) RegisterHandler (handler action.Action) error {
    _, ok := yy.handlers[handler.Name()]
    if ok {
        return ErrHandlerReregistered
    }
    yy.handlers[handler.Name()] = handler
    return nil
}

func (yy *Parser) Parse () (command.Command, error) {
    err := yy.lex.ReadCommand()
    if err != nil {
        return nil, err
    }

    commstr := string(yy.lex.Token())
    handler, ok := yy.handlers[commstr]
    if !ok {
        return nil, ErrNoHandler
    }

    return handler.Parse(yy.lex)
}
