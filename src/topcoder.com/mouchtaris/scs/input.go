package scs

import (
    net0 "topcoder.com/mouchtaris/scs/net"
    "net"
    "topcoder.com/mouchtaris/scs/util"
)

//
// Indefinitely serve incoming connections, by starting
// a Server routine for each one.
type ConnServer struct {
    server      *Server
    nbl         net0.NonBlockingListener
    em          ExecutionManager
    eh          util.ErrorHandler
}

//
//
func NewConnServer (s *Server, l net.Listener, em ExecutionManager) *ConnServer {
    return &ConnServer {
        server:     s,
        nbl:        net0.NewNonBlockingListener(l),
        eh:         util.NewErrorHandler("ConnServer"),
        em:         em,
    }
}

//
// Handle an incoming connection.
func (conserv *ConnServer) handleConnection (conn net.Conn) {
    conserv.em.Execute(func () {
        conserv.eh.ErrorsChannel() <- conserv.server.Serve(conn)
    })
}

//
// Server incoming connections indefinitely.
func (conserv *ConnServer) ServeAll () {
    for stream := range conserv.nbl.AcceptanceChannel() {
        conserv.handleConnection(stream)
    }
}

//
// Close this server.
//
// Stop accepting any new connections.
//
// Then handle all pending connections.
//
// Then wait for all running serving routines to finish.
//
// Finally close the underlying listener.
func (conserv *ConnServer) Close () error {
    return conserv.nbl.Close()
}
