package scs

import (
    "topcoder.com/mouchtaris/scs/lex"
    "topcoder.com/mouchtaris/scs/parser/action"
    "topcoder.com/mouchtaris/scs/parser"
    "io"
    "fmt"
)

//
// Serves clients, by parsing their input, translating it
// into Request-s which are later sent into the request
// channel for handling.
type Server struct {
    requests chan<- Request
    errors chan<- error
}

//
//
func NewServer (requests chan<- Request, errors chan<- error) *Server {
    return &Server {
        requests: requests,
        errors: errors,
    }
}

//
// "Server" an input stream by parsing it and sending the parsed
// input as requests to be handled.
//
func (server *Server) Serve (input io.ReadWriteCloser) error {
    lexer := lex.NewLexer(input)
    parser := parser.NewParser(lexer)
    parser.RegisterHandler(action.Set    { })
    parser.RegisterHandler(action.Get    { })
    parser.RegisterHandler(action.Delete { })
    parser.RegisterHandler(action.Stats  { })
    parser.RegisterHandler(action.Quit   { })

    comm, err := parser.Parse()
    for ; err == nil; comm, err = parser.Parse() {
        server.requests <- Request { comm, input.Write, }
    }
    if err == action.ErrQuit && comm != nil {
        server.requests <- Request { comm, input.Write, }
    }

    if err != nil && err != action.ErrQuit {
        errmsg := fmt.Sprintf("ERROR %s\r\n", err)
        input.Write([]byte(errmsg))
    }

    input.Close()
    return err
}

//
// Close this server. This closes the underlying
// commands channel and thus serving anything else
// will result in panic.
func (server *Server) Close () {
    close(server.requests)
}
