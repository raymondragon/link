package main

import (
    "net"
    "net/url"
    "strings"
    "sync"
)

func runBroker(parsedURL *url.URL, ipStore sync.Map) error {
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
            if parsedURL.Fragment != "" {
                clientIP, _, err := net.SplitHostPort(linkConn.RemoteAddr().String())
                if err != nil {
                    return
                }
                if _, exists := ipStore.Load(clientIP); !exists {
                    linkConn.Close()
                    return
                }
            }
            targetConn, err := net.DialTCP("tcp", nil, targetAddr)
            if err != nil {
                linkConn.Close()
                return
            }
            targetConn.SetNoDelay(true)
            handleTransmissions(linkConn, targetConn)
        }(linkConn)
    }
}
