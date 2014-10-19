package scs

import (
    "topcoder.com/mouchtaris/scs/lex"
    "topcoder.com/mouchtaris/scs/parser/action"
    "topcoder.com/mouchtaris/scs/parser"
    "io"
    "fmt"
)

type Server struct {
    sem chan uint32
    requests chan<- Request
    errors chan<- error
    done chan uint32
    running uint32
}

//
// Construct a new server, which will start at most "backlog"
// concurrent servers.
func NewServer (backlog uint32, requests chan<- Request, errors chan<- error) *Server {
    return &Server {
        sem: make(chan uint32, backlog),
        done: make(chan uint32, backlog),
        running: 0,
        requests: requests,
        errors: errors,
    }
}

//
// Drain the done queue.
//
// This function runs on the main routine.
func (server *Server) drainDone () {
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
// Try to claim a job pass.
// If there are too many jobs running currently
// this will block until some other job is done.
//
// This function runs on the main routine.
func (server *Server) jobSemaphoreDown () {
    // If we can append to the job semaphore, it's
    // ok to start another routine.
    //
    // If not, then we need to wait for other
    // jobs to complete.
    //
    // Some (or many) of them may be waiting
    // to signal their completion on the
    // "done" channel, so if there is no
    // space on the semaphore queue, we should
    // also try to drain the done-queue.
    for cont := true; cont; {
        select {
        case server.sem <- 1:
            cont = false
        default:
            server.drainDone()
        }
    }
}

//
// Mark a job opening.
// Send a job done signal down the done queue
// and releave the jobs semaphore queue by one.
//
// This function runs on the job subroutine.
func (server* Server) jobSemaphoreUp () {
    server.done <- 1
    <-server.sem
}

//
// Perform the actual per-subroutine work.
//
// This function runs on the job subroutine.
func serve (requests chan<- Request, input io.ReadWriteCloser) error {
    lexer := lex.NewLexer(input)
    parser := parser.NewParser(lexer)
    parser.RegisterHandler(action.Set    { })
    parser.RegisterHandler(action.Get    { })
    parser.RegisterHandler(action.Delete { })
    parser.RegisterHandler(action.Stats  { })
    parser.RegisterHandler(action.Quit   { })

    comm, err := parser.Parse()
    for ; err == nil; comm, err = parser.Parse() {
        requests <- Request { comm, input.Write, }
    }
    if err == action.ErrQuit && comm != nil {
        requests <- Request { comm, input.Write, }
    }

    if err != nil && err != action.ErrQuit {
        errmsg := fmt.Sprintf("ERROR %s\r\n", err)
        input.Write([]byte(errmsg))
    }

    input.Close()
    return err
}

//
// Wrap the actual work function in concurrency
// markers.
//
// This function runs on the job subroutine.
func (server* Server) wrapServing (input io.ReadWriteCloser) {
    err := serve(server.requests, input)
    if err != nil {
        server.errors<- err
    }
    server.jobSemaphoreUp()
}

//
// Server (asynchronously) parsing input from the given reader.
// Parsed commands are send to the server's command stream
// for further processing by some other consumer.
func (server *Server) GoServe (input io.ReadWriteCloser) {
    server.jobSemaphoreDown()
    server.running++
    go server.wrapServing(input)
}

//
// Wait for all running jobs to finish.
func (server *Server) Join () {
    for i := uint32(0); i < server.running; i++ {
        <-server.done
    }
    server.running = 0
}

//
// Close this server. This closes the underlying
// commands channel and thus serving anything else
// will result in panic.
func (server *Server) Close () {
    close(server.requests)
}
