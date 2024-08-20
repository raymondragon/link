package main

import (
    "net"
    "net/url"
    "strings"
    "time"
)

func runServer(parsedURL *url.URL) error {
    linkAddr, err := net.ResolveTCPAddr("tcp", parsedURL.Host)
    if err != nil {
        return err
    }
    targetAddr, err := net.ResolveTCPAddr("tcp", strings.TrimPrefix(parsedURL.Path, "/"))
    if err != nil {
        return err
    }
    linkListen, err := net.ListenTCP("tcp", linkAddr)
    if err != nil {
        return err
    }
    defer linkListen.Close()
    targetListen, err := net.ListenTCP("tcp", targetAddr)
    if err != nil {
        return err
    }
    defer targetListen.Close()
    var linkConn *net.TCPConn
    go func() {
        for {
            tempConn, err := linkListen.AcceptTCP()
            if err != nil {
                time.Sleep(1 * time.Second)
                continue
            }
            if linkConn != nil {
                linkConn.Close()
            }
            linkConn = tempConn
            linkConn.SetNoDelay(true)
            time.Sleep(1 * time.Second)
        }
    }()
    targetConn, err := targetListen.AcceptTCP()
    if err != nil {
        return err
    }
    targetConn.SetNoDelay(true)
    if linkConn == nil {
        targetConn.Close()
        return nil
    }
    handleConnections(linkConn, targetConn)
    return nil
}
