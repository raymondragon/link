package main

import (
    "log"
    "net/url"
    "os"
    "strings"
    "sync"
    "time"

    "github.com/raymondragon/link/pkg/handle"
    "github.com/raymondragon/link/pkg/run"
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
    var authorizedIP sync.Map
    if parsedURL.Fragment != "" {
        parsedAuthURL, err := url.Parse(parsedURL.Fragment)
        if err != nil {
            log.Fatalf("[ERRO] URL Parsing: %v", err)
        }
        log.Printf("[INFO] Authorization: %v", parsedAuthURL)
        go func() {
            for {
                if err := handle.Authorization(parsedAuthURL, &authorizedIP); err != nil {
                    log.Printf("[ERRO] Authorization: %v", err)
                    time.Sleep(1 * time.Second)
                    continue
                }
            }
        }()
    }
    log.Printf("[INFO] Transmissions: %v", strings.Split(rawURL, "#")[0])
    for {
        switch parsedURL.Scheme {
        case "server":
            if err := run.NewServer(parsedURL, &authorizedIP); err != nil {
                log.Printf("[ERRO] Server: %v", err)
                time.Sleep(1 * time.Second)
                continue
            }
        case "client":
            if err := run.NewClient(parsedURL); err != nil {
                log.Printf("[ERRO] Client: %v", err)
                time.Sleep(1 * time.Second)
                continue
            }
        case "broker":
            if err := run.NewBroker(parsedURL, &authorizedIP); err != nil {
                log.Printf("[ERRO] Broker: %v", err)
                time.Sleep(1 * time.Second)
                continue
            }
        default:
            log.Fatalf("[ERRO] Usage: server/client/broker://linkAddr/targetAddr#http/https://authAddr/secretPath")
        }
    }
}
