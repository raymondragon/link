package mode

import (
    "net"
    "net/url"
    "strings"
    "sync"
    "time"

    "github.com/raymondragon/link/pkg/handle"
)

func Server(parsedURL *url.URL, whiteList *sync.Map) error {
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
    linkChan := make(chan *net.TCPConn)
    targetChan := make(chan *net.TCPConn)
    go func() {
        for {
            tempConn, err := linkListen.AcceptTCP()
            if err != nil {
                time.Sleep(1 * time.Second)
                continue
            }
            tempConn.SetNoDelay(true)
            linkChan <- tempConn
        }
    }()
    go func() {
        for {
            tempConn, err := targetListen.AcceptTCP()
            if err != nil {
                time.Sleep(1 * time.Second)
                continue
            }
            tempConn.SetNoDelay(true)
            targetChan <- tempConn
        }
    }()
    tempSlot := make(chan struct{}, 1024)
    for {
        linkConn := <-linkChan
        tempSlot <- struct{}{}
        go func(linkConn *net.TCPConn) {
            defer func() { <-tempSlot }()
            targetConn := <-targetChan
            if parsedURL.Fragment != "" {
                clientIP, _, err := net.SplitHostPort(targetConn.RemoteAddr().String())
                if err != nil {
                    return
                }
                if _, exists := whiteList.Load(clientIP); !exists {
                    return
                }
            }
            handle.Conn(linkConn, targetConn)
        }(linkConn)
    }
}
