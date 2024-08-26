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
    var linkConn *net.TCPConn
    defer func() {
        if linkConn != nil {
            linkConn.Close()
        }
    }()
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
        }
    }()
    targetConn, err := targetListen.AcceptTCP()
    if err != nil {
        return err
    }
    defer targetConn.Close()
    targetConn.SetNoDelay(true)
    if parsedURL.Fragment != "" {
        clientIP, _, err := net.SplitHostPort(targetConn.RemoteAddr().String())
        if err != nil {
            return err
        }
        if _, exists := whiteList.Load(clientIP); !exists && linkConn != nil {
            return nil
        }
    }
    if _, err := linkConn.Write([]byte("targetConn")); err != nil {
        return nil
    }
    return handle.Conn(linkConn, targetConn)
}
