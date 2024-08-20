package main

import (
    "net"
    "net/http"
    "net/url"
)

func handleAuthorization(parsedURL *url.URL) error {
    http.HandleFunc(parsedURL.Path, func(w http.ResponseWriter, r *http.Request) {
        clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
        if err != nil {
            return
        }
        if _, err := w.Write([]byte(clientIP + "\n")); err != nil {
            return
        }
        authorizedIP.Store(clientIP, struct{}{})
    })
    if parsedURL.Scheme == "http" {
        if err := http.ListenAndServe(parsedURL.Host, nil); err != nil {
            return err
        }
    } else {
        tlsConfig, err := tlsConfigGeneration(parsedURL.Hostname())
        if err != nil {
            return err
        }
        authServer := &http.Server{
            Addr:      parsedURL.Host,
            TLSConfig: tlsConfig,
        }
        if err := authServer.ListenAndServeTLS("", ""); err != nil {
            return err
        }
    }
    return nil
}
