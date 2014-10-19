package scs

import (
	"fmt"
	"io"
	"topcoder.com/mouchtaris/scs/lex"
	"topcoder.com/mouchtaris/scs/parser"
	"topcoder.com/mouchtaris/scs/parser/action"
)

type Server struct {
	sem      chan uint32
	requests chan<- Request
	errors   chan<- error
	done     chan uint32
	running  uint32
}

//
// Construct a new server, which will start at most "backlog"
// concurrent servers.
func NewServer(backlog uint32, requests chan<- Request, errors chan<- error) *Server {
	return &Server{
		sem:      make(chan uint32, backlog),
		done:     make(chan uint32, backlog),
		running:  0,
		requests: requests,
		errors:   errors,
	}
}

//
//
func serve(requests chan<- Request, input io.ReadWriteCloser) error {
	lexer := lex.NewLexer(input)
	parser := parser.NewParser(lexer)
	parser.RegisterHandler(action.Set{})
	parser.RegisterHandler(action.Get{})
	parser.RegisterHandler(action.Delete{})
	parser.RegisterHandler(action.Stats{})
	parser.RegisterHandler(action.Quit{})

	comm, err := parser.Parse()
	for ; err == nil; comm, err = parser.Parse() {
		requests <- Request{comm, input.Write, input.Close}
	}
	if err == action.ErrQuit && comm != nil {
		requests <- Request{comm, input.Write, input.Close}
	}

	if err != nil && err != action.ErrQuit {
		errmsg := fmt.Sprintf("ERROR %s\r\n", err)
		input.Write([]byte(errmsg))
	}
	input.Close()
	return err
}

//
//
func (server *Server) wrapServing(input io.ReadWriteCloser) {
	err := serve(server.requests, input)
	if err != nil {
		server.errors <- err
	}
	server.done <- 1
	<-server.sem
}

//
//
func (server *Server) drainDone() {
	for cont := true; cont; {
		select {
		case <-server.done:
			server.running--
		default:
			cont = false
		}
	}
}

//
// Server (asynchronously) parsing input from the given reader.
// Parsed commands are send to the server's command stream
// for further processing by some other consumer.
func (server *Server) GoServe(input io.ReadWriteCloser) {
	for cont := true; cont; {
		select {
		case server.sem <- 1:
			cont = false
		default:
			server.drainDone()
		}
	}
	server.running++
	go server.wrapServing(input)
}

//
// Wait for all running jobs to finish.
func (server *Server) Join() {
	for i := uint32(0); i < server.running; i++ {
		<-server.done
	}
	server.running = 0
}

//
// Close this server. This closes the underlying
// commands channel and thus serving anything else
// will result in panic.
func (server *Server) Close() {
	close(server.requests)
}
