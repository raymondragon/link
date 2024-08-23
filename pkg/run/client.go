package run

import (
    "net"
    "net/url"
    "strings"
    "time"

    "github.com/raymondragon/link/pkg/handle"
)

func NewClient(parsedURL *url.URL) error {
    linkAddr, err := net.ResolveTCPAddr("tcp", parsedURL.Host)
    if err != nil {
        return err
    }
    targetAddr, err := net.ResolveTCPAddr("tcp", strings.TrimPrefix(parsedURL.Path, "/"))
    if err != nil {
        return err
    }
    linkConn, err := net.DialTCP("tcp", nil, linkAddr)
    if err != nil {
        return err
    }
    linkConn.SetNoDelay(true)
    targetConn, err := net.DialTCP("tcp", nil, targetAddr)
    if err != nil {
        linkConn.Close()
        return err
    }
    targetConn.SetNoDelay(true)
    go func() {
        buffer := make([]byte, 1024)
        for {
            targetConn.SetReadDeadline(time.Now().Add(10 * time.Second))
            if _, err := targetConn.Read(buffer); err != nil {
                if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
                    targetConn.Close()
                    linkConn.Close()
                    return
                }
                continue
            }
        }
    }()
    handle.Transmissions(linkConn, targetConn)
    return nil
}
