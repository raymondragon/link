package main

import (
    "log"
    "net/url"
    "os"
    "strings"
    "sync"
    "time"
)

func main() {
    if len(os.Args) < 2 {
        log.Fatalf("[ERRO] Usage: server/client/broker://linkAddr/targetAddr#http/https://authAddr/secretPath")
    }
    rawURL := os.Args[1]
    parsedURL, err := url.Parse(rawURL)
    if err != nil {
        log.Fatalf("[ERRO] URL Parsing: %v", err)
    }
    var ipStore sync.Map
    if parsedURL.Fragment != "" {
        parsedAuthURL, err := url.Parse(parsedURL.Fragment)
        if err != nil {
            log.Fatalf("[ERRO] URL Parsing: %v", err)
        }
        log.Printf("[INFO] Authorization: %v", parsedAuthURL)
        go func() {
            if err := handleAuthorization(parsedAuthURL, ipStore); err != nil {
                log.Fatalf("[ERRO] Authorization: %v", err)
            }
        }()
    }
    log.Printf("[INFO] Transmissions: %v", strings.Split(rawURL, "#")[0])
    for {
        switch parsedURL.Scheme {
        case "server":
            if err := runServer(parsedURL, ipStore); err != nil {
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
        case "broker":
            if err := runBroker(parsedURL, ipStore); err != nil {
                log.Printf("[ERRO] Broker: %v", err)
                time.Sleep(1 * time.Second)
                continue
            }
        default:
            log.Fatalf("[ERRO] Usage: server/client/broker://linkAddr/targetAddr#http/https://authAddr/secretPath")
        }
    }
}
