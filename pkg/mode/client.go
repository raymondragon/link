package mode

import (
    "net"
    "net/url"
    "strings"

    "github.com/raymondragon/link/pkg/handle"
)

func Client(parsedURL *url.URL) error {
    linkAddr, err := net.ResolveTCPAddr("tcp", parsedURL.Host)
    if err != nil {
        return err
    }
    targetAddr, err := net.ResolveTCPAddr("tcp", strings.TrimPrefix(parsedURL.Path, "/"))
    if err != nil {
        return err
    }
    tempSlot := make(chan struct{}, 5)
    for {
        linkConn, err := net.DialTCP("tcp", nil, linkAddr)
        if err != nil {
            continue
        }
        linkConn.SetNoDelay(true)
        targetConn, err := net.DialTCP("tcp", nil, targetAddr)
        if err != nil {
            linkConn.Close()
            continue
        }
        targetConn.SetNoDelay(true)
        tempSlot <- struct{}{}
        go func(linkConn, targetConn *net.TCPConn) {
            defer func() { <-tempSlot }()
            handle.Conn(linkConn, targetConn)
        }(linkConn, targetConn)
    }
}
