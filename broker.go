package main

import (
    "net"
    "net/url"
    "strings"
)

func runBroker(parsedURL *url.URL) error {
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
    semTEMP := make(chan struct{}, 1024)
    for {
        linkConn, err := linkListen.AcceptTCP()
        if err != nil {
            continue
        }
        linkConn.SetNoDelay(true)
        semTEMP <- struct{}{}
        go func(linkConn net.Conn) {
            defer func() { <-semTEMP }()
            targetConn, err := net.DialTCP("tcp", nil, targetAddr)
            if err != nil {
                linkConn.Close()
                return
            }
            targetConn.SetNoDelay(true)
            handleConnections(linkConn, targetConn)
        }(linkConn)
    }
}
