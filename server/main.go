package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// Load server certificate and key
	serverCert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		fmt.Println("Error loading server certificate:", err)
		return
	}

	// Load CA certificate
	caCert, err := ioutil.ReadFile("ca.crt")
	if err != nil {
		fmt.Println("Error reading CA certificate:", err)
		return
	}

	// Create a certificate pool and add CA certificate to it
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create a HTTPS server with mutual TLS authentication
	server := &http.Server{
		Addr: ":8443", // Server listens on port 8443
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{serverCert}, // Set server certificate
			ClientAuth:   tls.RequireAndVerifyClientCert, // Require client certificate and verify it
			ClientCAs:    caCertPool, // Set trusted CA certificates for client verification
		},
	}

	// Define a handler for incoming requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, TLS client!\n") // Respond to client requests
	})

	// Start the HTTPS server
	fmt.Println("Server listening on https://localhost:8443")
	err = server.ListenAndServeTLS("", "") // Start server with provided TLS configuration
	if err != nil {
		fmt.Println("Error starting HTTPS server:", err)
	}
}
