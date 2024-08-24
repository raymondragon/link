package main

import (
    "log"
    "net/url"
    "os"
    "strings"
    "sync"
    "time"

    "github.com/raymondragon/link/pkg/handle"
    "github.com/raymondragon/link/pkg/mode"
)

func main() {
    if len(os.Args) < 2 {
        log.Fatalf("[ERRO] Usage: server/client/broker://linkAddr/targetAddr#http/https://authAddr/secretPath")
    }
    rawURL := os.Args[1]
    parsedURL, err := url.Parse(rawURL)
    if err != nil {
        log.Fatalf("[ERRO] URL: %v", err)
    }
    var whiteList sync.Map
    if parsedURL.Fragment != "" {
        parsedAuthURL, err := url.Parse(parsedURL.Fragment)
        if err != nil {
            log.Fatalf("[ERRO] URL: %v", err)
        }
        log.Printf("[INFO] Auth: %v", parsedAuthURL)
        go func() {
            for {
                if err := handle.Auth(parsedAuthURL, &whiteList); err != nil {
                    log.Printf("[ERRO] Auth: %v", err)
                    time.Sleep(1 * time.Second)
                    continue
                }
            }
        }()
    }
    for {
        switch parsedURL.Scheme {
        case "server":
            log.Printf("[INFO] Server: %v", strings.Split(rawURL, "#")[0])
            if err := mode.Server(parsedURL, &whiteList); err != nil {
                log.Printf("[ERRO] Server: %v", err)
                time.Sleep(1 * time.Second)
                continue
            }
        case "client":
            log.Printf("[INFO] Client: %v", strings.Split(rawURL, "#")[0])
            if err := mode.Client(parsedURL); err != nil {
                log.Printf("[ERRO] Client: %v", err)
                time.Sleep(1 * time.Second)
                continue
            }
        case "broker":
            log.Printf("[INFO] Broker: %v", strings.Split(rawURL, "#")[0])
            if err := mode.Broker(parsedURL, &whiteList); err != nil {
                log.Printf("[ERRO] Broker: %v", err)
                time.Sleep(1 * time.Second)
                continue
            }
        default:
            log.Fatalf("[ERRO] Usage: server/client/broker://linkAddr/targetAddr#http/https://authAddr/secretPath")
        }
    }
}
