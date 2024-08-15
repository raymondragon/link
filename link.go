package main

import (
    "log"
    "net"
    "net/url"
    "os"
    "strings"
    "time"
)

func main() {
    if len(os.Args) < 2 {
        log.Fatalf("[ERRO] Usage: server/client://linkAddr/targetAddr")
    }
    rawURL := os.Args[1]
    parsedURL, err := url.Parse(rawURL)
    if err != nil {
        log.Fatalf("[ERRO] URL Parsing: %v", err)
    }
    log.Printf("[INFO] %v", parsedURL)
    if parsedURL.Scheme == "broker" {
        if err := runBroker(parsedURL); err != nil {
            log.Printf("[ERRO] Broker: %v", err)
        }
    }
    for {
        switch parsedURL.Scheme {
        case "server":
            if err := runServer(parsedURL); err != nil {
                log.Printf("[ERRO] Server: %v", err)
                time.Sleep(1 * time.Second)
                continue
            }
        case "client":
            if err := runClient(parsedURL); err != nil {
                log.Printf("[ERRO] Client: %v", err)
                time.Sleep(1 * time.Second)
                continue
            }
        default:
            log.Fatalf("[ERRO] Usage: server/client://linkAddr#targetAddr")
        }
    }
}

func runBroker(parsedURL *url.URL) error {
    linkAddr, err := net.ResolveTCPAddr("tcp", parsedURL.Host)
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
        linkConn, err := linkListen.Accept()
        if err != nil {
            continue
        }
        linkConn.SetNoDelay(true)
        semTEMP <- struct{}{}
        go func(linkConn net.Conn) {
            defer func() { <-semTEMP }()
            targetConn, err := net.Dial("tcp", strings.TrimPrefix(parsedURL.Path, "/"))
            if err != nil {
                linkConn.Close()
                return
            }
            targetConn.SetNoDelay(true)
            handleConnections(linkConn, targetConn)
        }(linkConn)
    }
}

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

func runClient(parsedURL *url.URL) error {
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
    handleConnections(linkConn, targetConn)
    return nil
}
