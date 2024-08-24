package autotls

import (
    "crypto/tls"
)

func Setup(username, hostname string) (*tls.Config, error) {
    tlsConfig, err := Application(username, hostname)
    if err != nil {
        tlsConfig, err = Generation(hostname)
    }
    return tlsConfig, err
}
