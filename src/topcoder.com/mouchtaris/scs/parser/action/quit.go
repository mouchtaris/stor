package action

import (
	"errors"
	"topcoder.com/mouchtaris/scs/command"
	"topcoder.com/mouchtaris/scs/lex"
)

var ErrQuit = errors.New("quiting")

//
type Quit struct{}

//
func (Quit) Name() string {
	return "quit"
}

//
func (Quit) Parse(lex *lex.Lexer) (command.Command, error) {
	comm := command.Quit{}

	err := lex.ReadEOC()
	if err != nil {
		return nil, err
	}

	return &comm, ErrQuit
}
