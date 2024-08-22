package run

import (
    "net"
    "net/url"
    "strings"
)

func newClient(parsedURL *url.URL) error {
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
    handleTransmissions(linkConn, targetConn)
    return nil
}
