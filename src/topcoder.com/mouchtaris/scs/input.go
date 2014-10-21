package scs

import (
    net0 "topcoder.com/mouchtaris/scs/net"
    "net"
    "topcoder.com/mouchtaris/scs/util"
)

func ServeIncomingTCPClients (s *Server, l net.Listener, em ExecutionManager, stop <-chan byte) {
    nbl     := net0.NewNonBlockingListener(l)
    errors  := util.NewErrorHandler("ServeIncomnigTCPClients").ErrorsChannel()

    for {
        select {
        case <-stop:
            return
        case stream := <-nbl.AcceptanceChannel():
            em.Execute(func () { errors <- s.Serve(stream) })
        }
    }
}
