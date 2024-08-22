package tlsconfig

import (
    "crypto/tls"

    "github.com/caddyserver/certmagic"
)

func Application(username, hostname string) (*tls.Config, error) {
    certmagic.DefaultACME.CA = certmagic.LetsEncryptProductionCA
    certmagic.DefaultACME.Agreed = true
    mainAddr := username + "@" + hostname
    if username == "" {
        mainAddr = "no-reply@" + hostname
    }
    certmagic.DefaultACME.Email = mainAddr
    tlsConfig, err := certmagic.TLS([]string{hostname})
    return tlsConfig, err
}
