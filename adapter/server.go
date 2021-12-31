package adapter

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

func Serve(mux *http.ServeMux, addr string) {
	if addr == "" {
		addr = ":8080"
	}

	httpServer := &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	serverError := make(chan error)
	go func() {
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			serverError <- err
			return
		}
		fmt.Println("http server shutdown")
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	log.Println("running http server")
	select {
	case err := <-serverError:
		fmt.Println("http server error", err)
		time.Sleep(100 * time.Millisecond)
		os.Exit(1)
	case sig := <-quit:
		fmt.Println("\ngot signal", sig)
	}

	shutdown(httpServer)
	time.Sleep(100 * time.Millisecond)
}

func split(r rune) bool {
	return r == ',' || r == ' '
}

func ServeTLS(mux *http.ServeMux) {

	tlsConfig := tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.CurveP521,
			tls.CurveP384,
			tls.CurveP256,
			tls.X25519,
		},
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			// Best disabled, as they don't provide Forward Secrecy,
			// but might be necessary for some clients
			// tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			// tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	domains := strings.FieldsFunc(os.Getenv("ALLOW"), split)
	for _, domain := range domains {
		cert, err := tls.LoadX509KeyPair(domain+".crt", domain+".key")
		if err != nil {
			e := fmt.Errorf("loading certificate for domain %s -> %w", domain, err)
			fmt.Println(e)
			os.Exit(1)
		}

		tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
	}

	httpServer := &http.Server{
		Addr: ":80",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			host, _, _ := net.SplitHostPort(r.Host)
			u := r.URL
			u.Host = net.JoinHostPort(host, ":443")
			u.Scheme = "https"
			http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
		}),
	}

	httpsServer := &http.Server{
		Addr:           ":443",
		Handler:        mux,
		TLSConfig:      &tlsConfig,
		TLSNextProto:   make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	serverHTTPError := make(chan error)
	serverHTTPSError := make(chan error)
	go func() {
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			serverHTTPError <- err
			return
		}
		fmt.Println("http server shutdown")
	}()
	go func() {
		err := httpsServer.ListenAndServeTLS("", "")
		if err != nil && err != http.ErrServerClosed {
			serverHTTPSError <- err
			return
		}
		fmt.Println("https server shutdown")
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	log.Println("running https server")
	select {
	case err := <-serverHTTPError:
		fmt.Println("server error", err)
		shutdown(httpsServer)
		time.Sleep(100 * time.Millisecond)
		os.Exit(1)
	case err := <-serverHTTPSError:
		fmt.Println("server error", err)
		shutdown(httpServer)
		time.Sleep(100 * time.Millisecond)
		os.Exit(1)
	case sig := <-quit:
		fmt.Println("\ngot signal", sig)
	}

	shutdown(httpServer)
	shutdown(httpsServer)
	time.Sleep(100 * time.Millisecond)
}

func shutdown(s *http.Server) {
	ctxServer, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := s.Shutdown(ctxServer)
	if err != nil && err != http.ErrServerClosed {
		fmt.Println("server shutdown error", err)
	}
}
