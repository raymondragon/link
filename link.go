package main

import (
    "io"
    "log"
    "net"
    "net/url"
    "os"
    "time"
)

func main() {
    if len(os.Args) < 2 {
        log.Fatalf("[ERRO] Usage: server/client://linkAddr#targetAddr")
    }
    rawURL := os.Args[1]
    parsedURL, err := url.Parse(rawURL)
    if err != nil {
        log.Fatalf("[ERRO] URL Parsing: %v", err)
    }
    for {
        switch parsedURL.Scheme {
        case "server":
            log.Printf("[INFO] Server: %v <-- %v", parsedURL.Host, parsedURL.Fragment)
            if err := runServer(parsedURL); err != nil {
                log.Fatalf("[ERRO] Server: %v", err)
            }
        case "client":
            log.Printf("[INFO] Client: %v --> %v", parsedURL.Host, parsedURL.Fragment)
            if err := runClient(parsedURL); err != nil {
                log.Fatalf("[ERRO] Client: %v", err)
            }
            time.Sleep(1 * time.Second)
        default:
            log.Fatalf("[ERRO] Usage: server/client://linkAddr#targetAddr")
        }
    }
}

func runServer(parsedURL *url.URL) error {
    linkAddr := parsedURL.Host
    serverAddr := parsedURL.Fragment
    linkListen, err := net.Listen("tcp", linkAddr)
    if err != nil {
        return err
    }
    defer linkListen.Close()
    serverListen, err := net.Listen("tcp", serverAddr)
    if err != nil {
        return err
    }
    defer serverListen.Close()
    linkConn, err := linkListen.Accept()
    if err != nil {
        return err
    }
    timeChan := time.After(1 * time.Minute)
    doneChan := make(chan struct{})
    var serverConn net.Conn
    go func() {
        defer close(doneChan)
        var err error
        serverConn, err = serverListen.Accept()
        if err != nil {
            return err
        }
    }()
    select {
    case <-timeChan:
    case <-doneChan:
        handleConnections(linkConn, serverConn)
    }
    return nil
}

func runClient(parsedURL *url.URL) error {
    linkAddr := parsedURL.Host
    clientAddr := parsedURL.Fragment
    linkConn, err := net.Dial("tcp", linkAddr)
    if err != nil {
        return err
    }
    clientConn, err := net.Dial("tcp", clientAddr)
    if err != nil {
        return err
    }
    handleConnections(linkConn, clientConn)
    return nil
}

func handleConnections(conn1, conn2 net.Conn) {
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
