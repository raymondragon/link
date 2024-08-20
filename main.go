package main

import (
    "log"
    "os"
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
    log.Printf("[INFO] %v", parsedURL)
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
        case "broker":
            if err := runBroker(parsedURL); err != nil {
                log.Printf("[ERRO] Broker: %v", err)
                time.Sleep(1 * time.Second)
                continue
            }
        default:
            log.Fatalf("[ERRO] Usage: server/client/broker://linkAddr/targetAddr#http/https://authAddr/secretPath")
        }
    }
}
