package net

import (
    "net"
)

//
//
type NonBlockingListener interface {
    //
    //
    AcceptanceChannel () <-chan net.Conn

    //
    //
    Close () error
}


//
//
type nonBlockingListener struct {
    l           net.Listener
    acceptances chan net.Conn
}

//
//
func acceptForever (l net.Listener, acceptances chan<- net.Conn) {
    for {
        conn, err := l.Accept()
        if err != nil {
            acceptances <- conn
        }
    }
}

//
//
func NewNonBlockingListener (l net.Listener) NonBlockingListener {
    acceptances := make(chan net.Conn, 20)
    go acceptForever(l, acceptances)

    return &nonBlockingListener {
        l:              l,
        acceptances:    acceptances,
    }
}

//
//
func (nbl *nonBlockingListener) AcceptanceChannel () <-chan net.Conn {
    return nbl.acceptances
}

//
//
func (nbl *nonBlockingListener) Close () error {
    return nbl.l.Close()
}
