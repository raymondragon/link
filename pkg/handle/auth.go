package handle

import (
    "net"
    "net/http"
    "net/url"
    "sync"

    "github.com/raymondragon/link/pkg/autotls"
)

func Auth(parsedURL *url.URL, whiteList *sync.Map) error {
    http.HandleFunc(parsedURL.Path, func(w http.ResponseWriter, r *http.Request) {
        clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
        if err != nil {
            return
        }
        if _, err := w.Write([]byte(clientIP + "\n")); err != nil {
            return
        }
        whiteList.Store(clientIP, struct{}{})
    })
    if parsedURL.Scheme == "http" {
        if err := http.ListenAndServe(parsedURL.Host, nil); err != nil {
            return err
        }
    } else {
        tlsConfig, err := tlsconfig.Setup(parsedURL.User.Username(), parsedURL.Hostname())
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
