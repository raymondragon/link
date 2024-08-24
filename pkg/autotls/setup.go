package autotls

import (
    "crypto/tls"
)

func Setup(username, hostname string) (*tls.Config, error) {
    tlsConfig, err := Register(username, hostname)
    if err != nil {
        tlsConfig, err = Generate(hostname)
    }
    return tlsConfig, err
}
