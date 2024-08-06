package main

import (
    "io"
    "log"
    "net"
    "net/url"
    "os"
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
    switch parsedURL.Scheme {
    case "server":
        log.Printf("[INFO] Linking Server: %v <-- %v", parsedURL.Host, parsedURL.Fragment)
        if err := runServer(parsedURL); err != nil {
            log.Fatalf("[ERRO] Server: %v", err)
        }
    case "client":
        log.Printf("[INFO] Linking Client: %v --> %v", parsedURL.Host, parsedURL.Fragment)
        if err := runClient(parsedURL); err != nil {
            log.Fatalf("[ERRO] Client: %v", err)
        }
    default:
        log.Fatalf("[ERRO] Usage: server/client://linkAddr#targetAddr")
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
    for {
        linkConn, err := linkListen.Accept()
        if err != nil {
            continue
        }
        serverConn, err := serverListen.Accept()
        if err != nil {
            linkConn.Close()
            continue
        }
        go handleConnection(linkConn, serverConn)
    }
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
    handleConnection(linkConn, clientConn)
    select{}
}

func handleConnection(conn1, conn2 net.Conn) {
    defer conn1.Close()
    defer conn2.Close()
    go io.Copy(conn1, conn2)
    go io.Copy(conn2, conn1)
}