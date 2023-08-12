package main

import (
	"crypto/tls"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"

	"golang.org/x/crypto/acme/autocert"
)

func makeHttpsListener(domain string, port int) (net.Listener, error) {
	// Certificate manager
	m := &autocert.Manager{
		Prompt: autocert.AcceptTOS,
		// Replace with your domain
		HostPolicy: autocert.HostWhitelist(domain),
		// Folder to store the certificates
		Cache: autocert.DirCache("./certs"),
	}

	// TLS Config
	cfg := &tls.Config{
		// Get Certificate from Let's Encrypt
		GetCertificate: m.GetCertificate,
		// By default NextProtos contains the "h2"
		// This has to be removed since Fasthttp does not support HTTP/2
		// Or it will cause a flood of PRI method logs
		// http://webconcepts.info/concepts/http-method/PRI
		NextProtos: []string{
			"http/1.1", "acme-tls/1",
		},
	}
	log.Infof("Listening on port %v for HTTPS requests to %v", port, domain)
	return tls.Listen("tcp", fmt.Sprintf(":%v", port), cfg)
}
