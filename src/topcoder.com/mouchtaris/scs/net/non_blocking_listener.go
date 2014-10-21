package net

import (
    "net"
)

//
//
type NonBlockingListener interface {
    // Start accepting connections, which are send to
    // the AcceptanceChannel().
    StartAccepting ()
    // Stop accepting connections. This does not close
    // the underlying listener, so connections will still
    // be accepted by the OS on a lower level (TCP server
    // socket, for instance, if it is a TCPListener).
    // This will, however, close the current AcceptanceChannel(),
    // so that any consumer routines will stop processing.
    StopAccepting ()
    //
    AcceptanceChannel () <-chan net.Conn

    //
    // Implies StopAccepting().
    Close () error
}

//
//
type nonBlockingListener struct {
    l               net.Listener
    acceptances     chan net.Conn
    stop, stopped   chan byte
}

//
//
func (nbl nonBlockingListener) newAcceptancesChannel () {
    nbl.acceptances = make(chan net.Conn, 20) // TODO magic number 20
}

//
//
// TODO: the last connection accepted could stall the server
// shutdown forever. There should be more graceful handling
// for this case (the server should just serve-by-bye-bye it).
func (nbl nonBlockingListener) acceptForever () {
    for {
        select {
        case <-nbl.stop:
            nbl.stopped <- 1
            return
        default:
            conn, err := nbl.l.Accept()
            if err != nil {
                nbl.acceptances <- conn
                // TODO where does err go?
            }
        }
    }
}

//
//
func NewNonBlockingListener (l net.Listener) NonBlockingListener {
    return &nonBlockingListener {
        l:              l,
        stop:           make(chan byte, 1),
        stopped:        make(chan byte, 1),
    }
}

//
//
func (nbl nonBlockingListener) AcceptanceChannel () <-chan net.Conn {
    return nbl.acceptances
}

//
//
func (nbl nonBlockingListener) StartAccepting () {
    nbl.newAcceptancesChannel()
    go nbl.acceptForever()
}

//
//
func (nbl nonBlockingListener) StopAccepting () {
    nbl.stop <- 1
    <-nbl.stopped
    close(nbl.acceptances)
}

//
//
func (nbl nonBlockingListener) Close () error {
    err := nbl.l.Close()
    nbl.StopAccepting()
    return err
}
