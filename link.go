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
            log.Printf("[INFO] Link: %v <-- %v", parsedURL.Host, parsedURL.Fragment)
            if err := runServer(parsedURL); err != nil {
                log.Printf("[ERRO] Server: %v", err)
                time.Sleep(1 * time.Second)
                continue
            }
        case "client":
            log.Printf("[INFO] Link: %v --> %v", parsedURL.Host, parsedURL.Fragment)
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

func runServer(parsedURL *url.URL) error {
    linkAddr, err := net.ResolveTCPAddr("tcp", parsedURL.Host)
    if err != nil {
        return err
    }
    serverAddr, err := net.ResolveTCPAddr("tcp", parsedURL.Fragment)
    if err != nil {
        return err
    }
    linkListen, err := net.ListenTCP("tcp", linkAddr)
    if err != nil {
        return err
    }
    defer linkListen.Close()
    serverListen, err := net.ListenTCP("tcp", serverAddr)
    if err != nil {
        return err
    }
    defer serverListen.Close()
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
    serverConn, err := serverListen.AcceptTCP()
    if err != nil {
        return err
    }
    serverConn.SetNoDelay(true)
    if linkConn == nil {
        serverConn.Close()
        return nil
    }
    handleConnections(linkConn, serverConn)
    return nil
}

func runClient(parsedURL *url.URL) error {
    linkAddr, err := net.ResolveTCPAddr("tcp", parsedURL.Host)
    if err != nil {
        return err
    }
    clientAddr, err := net.ResolveTCPAddr("tcp", parsedURL.Fragment)
    if err != nil {
        return err
    }
    linkConn, err := net.DialTCP("tcp", nil, linkAddr)
    if err != nil {
        return err
    }
    linkConn.SetNoDelay(true)
    clientConn, err := net.DialTCP("tcp", nil, clientAddr)
    if err != nil {
        linkConn.Close()
        return err
    }
    clientConn.SetNoDelay(true)
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