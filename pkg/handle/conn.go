package handle

import (
    "io"
    "net"
)

func Conn(conn1, conn2 net.Conn) {
    done := make(chan struct{}, 2)
    go func() {
        defer conn1.Close()
        io.Copy(conn1, conn2)
        done <- struct{}{}
    }()
    go func() {
        defer conn2.Close()
        io.Copy(conn2, conn1)
        done <- struct{}{}
    }()
    <-done
    <-done
}
