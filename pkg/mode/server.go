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
    linkConn, err := linkListen.AcceptTCP()
    if err != nil {
        return err
    }
    linkConn.SetNoDelay(true)
    targetConn, err := targetListen.AcceptTCP()
    if err != nil {
        return err
    }
    targetConn.SetNoDelay(true)
    if parsedURL.Fragment != "" {
        clientIP, _, err := net.SplitHostPort(targetConn.RemoteAddr().String())
        if err != nil {
            targetConn.Close()
            return err
        }
        if _, exists := whiteList.Load(clientIP); !exists {
            targetConn.Close()
            return nil
        }
    }
    if linkConn == nil {
        targetConn.Close()
        return nil
    }
    if _, err := linkConn.Write([]byte("targetConn")); err != nil {
        targetConn.Close()
        linkConn.Close()
        return err
    }
    tempBuff := make([]byte, 1024)
    n, err := linkConn.Read(tempBuff)
    if err != nil {
        targetConn.Close()
        linkConn.Close()
        return err
    }
    if string(tempBuff[:n]) == "targetReady" {
        handle.Conn(linkConn, targetConn)
    }
    return nil
}
