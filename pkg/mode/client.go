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
    linkConn, err := net.DialTCP("tcp", nil, linkAddr)
    if err != nil {
        return err
    }
    linkConn.SetNoDelay(true)
    tempBuff := make([]byte, 1024)
    n, err := linkConn.Read(tempBuff)
    if err != nil {
        linkConn.Close()
        return err
    }
    if string(tempBuff[:n]) == "targetConn" {
        targetConn, err := net.DialTCP("tcp", nil, targetAddr)
        if err != nil {
            linkConn.Close()
            return err
        }
        targetConn.SetNoDelay(true)
        if _, err := linkConn.Write([]byte("targetReady")); err != nil {
            targetConn.Close()
            linkConn.Close()
            return err
        }
        handle.Conn(linkConn, targetConn)
    }
    return nil
}
