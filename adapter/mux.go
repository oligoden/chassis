package adapter

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/oligoden/chassis"
	"github.com/oligoden/chassis/storage/gosql"
	"google.golang.org/grpc"
)

type Mux struct {
	Mux         *http.ServeMux
	Domain      string
	URL         *url.URL
	Stores      map[string]*gosql.Store
	RPDs        []string
	GRPCServers []*grpc.Server
	GRPCPorts   []uint
	Err         error
}

func NewMux() *Mux {
	return &Mux{
		Mux:         http.NewServeMux(),
		Stores:      map[string]*gosql.Store{},
		RPDs:        []string{},
		GRPCServers: []*grpc.Server{},
		GRPCPorts:   []uint{},
	}
}

func (mx *Mux) SetDomain(domain string) *Mux {
	mx.Domain = domain
	return mx
}

func (mx *Mux) SetURL(us string) *Mux {
	u, err := url.Parse(us)
	if err != nil {
		mx.Err = err
		return mx
	}

	mx.URL = u
	return mx
}

func (mx *Mux) SetStore(key string, store *gosql.Store) *Mux {
	mx.Stores[key] = store
	return mx
}

func (mx *Mux) AddRPD(dest string) *Mux {
	mx.RPDs = append(mx.RPDs, dest)
	return mx
}

func (mx *Mux) AddGRPC(srv *grpc.Server, port uint) *Mux {
	mx.GRPCServers = append(mx.GRPCServers, srv)
	mx.GRPCPorts = append(mx.GRPCPorts, port)
	return mx
}

func (mx *Mux) Compile(hs func(*Mux)) *http.ServeMux {
	hs(mx)
	return mx.Mux
}

func (mx *Mux) Register(hs func(*Mux)) *Mux {
	hs(mx)
	return mx
}

func (mx *Mux) ServeMux() *http.ServeMux {
	return mx.Mux
}

func (mx *Mux) Adapter() Adapter {
	return Adapter{
		Host: mx.URL.Hostname(),
		mx:   mx,
	}
}

func (mx *Mux) Handle(pattern string) Adapter {
	return Adapter{
		Host:    mx.URL.Hostname(),
		pattern: pattern,
		mx:      mx,
	}
}

func (mx *Mux) Serve() {
	if mx.URL.Scheme == "http" {
		mx.serve()
		return
	}

	if mx.URL.Scheme == "https" {
		mx.serveTLS()
		return
	}
}

type grpcError struct {
	i int
	e error
}

func (mx *Mux) serve() {
	_, port, err := net.SplitHostPort(mx.URL.Host)
	if err != nil {
		fmt.Println(chassis.Mark("getting port to serve", err))
		os.Exit(1)
	}

	addr := fmt.Sprintf(":%s", port)
	if addr == "" {
		addr = ":8080"
	}

	httpServer := &http.Server{
		Addr:           addr,
		Handler:        mx.Mux,
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

	grpcServerError := make(chan grpcError)
	for i, grpcServer := range mx.GRPCServers {
		go func() {
			lis, err := net.Listen("tcp", fmt.Sprintf(":%d", mx.GRPCPorts[i]))
			if err != nil {
				e := chassis.Mark("failed to listen", err)
				grpcServerError <- grpcError{i, e}
				return
			}

			if err := grpcServer.Serve(lis); err != nil {
				e := chassis.Mark("failed to serve", err)
				grpcServerError <- grpcError{i, e}
				return
			}
		}()
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	log.Println("started servers")
	select {
	case err := <-serverError:
		e := chassis.Mark("http server error", err)
		fmt.Println(chassis.ErrorTrace(e))
		for _, grpcServer := range mx.GRPCServers {
			grpcServer.GracefulStop()
		}
		time.Sleep(100 * time.Millisecond)
		os.Exit(1)
	case gerr := <-grpcServerError:
		e := chassis.Mark("grpc server error", gerr.e)
		fmt.Println(chassis.ErrorTrace(e))
		shutdown(httpServer)
		for i, grpcServer := range mx.GRPCServers {
			if i != gerr.i {
				grpcServer.GracefulStop()
			}
		}
		time.Sleep(100 * time.Millisecond)
		os.Exit(1)
	case sig := <-quit:
		fmt.Println("\ngot signal", sig)
	}

	shutdown(httpServer)
	for _, grpcServer := range mx.GRPCServers {
		grpcServer.GracefulStop()
	}
	time.Sleep(100 * time.Millisecond)
}

func (mx *Mux) serveTLS() {
	_, port, err := net.SplitHostPort(mx.URL.Host)
	if err != nil {
		fmt.Println(chassis.Mark("getting port to serve", err))
		os.Exit(1)
	}

	addr := fmt.Sprintf(":%s", port)
	if addr == "" {
		addr = ":8080"
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
		cert, err := tls.LoadX509KeyPair("certs/certificates/"+domain+".crt", "certs/certificates/"+domain+".key")
		if err != nil {
			e := fmt.Errorf("loading certificate for domain %s -> %w", domain, err)
			fmt.Println(e)
			os.Exit(1)
		}

		tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
	}

	httpsServer := &http.Server{
		Addr:           ":443",
		Handler:        mx.Mux,
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

	grpcServerError := make(chan grpcError)
	for i, grpcServer := range mx.GRPCServers {
		go func() {
			lis, err := net.Listen("tcp", fmt.Sprintf(":%d", mx.GRPCPorts[i]))
			if err != nil {
				e := chassis.Mark("failed to listen", err)
				grpcServerError <- grpcError{i, e}
				return
			}

			if err := grpcServer.Serve(lis); err != nil {
				e := chassis.Mark("failed to serve", err)
				grpcServerError <- grpcError{i, e}
				return
			}
		}()
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	log.Println("started servers")
	select {
	case err := <-serverHTTPError:
		e := chassis.Mark("http server error", err)
		fmt.Println(chassis.ErrorTrace(e))
		shutdown(httpsServer)
		for _, grpcServer := range mx.GRPCServers {
			grpcServer.GracefulStop()
		}
		time.Sleep(100 * time.Millisecond)
		os.Exit(1)
	case err := <-serverHTTPSError:
		e := chassis.Mark("https server error", err)
		fmt.Println(chassis.ErrorTrace(e))
		shutdown(httpServer)
		for _, grpcServer := range mx.GRPCServers {
			grpcServer.GracefulStop()
		}
		time.Sleep(100 * time.Millisecond)
		os.Exit(1)
	case gerr := <-grpcServerError:
		e := chassis.Mark("grpc server error", gerr.e)
		fmt.Println(chassis.ErrorTrace(e))
		shutdown(httpServer)
		shutdown(httpsServer)
		for i, grpcServer := range mx.GRPCServers {
			if i != gerr.i {
				grpcServer.GracefulStop()
			}
		}
		time.Sleep(100 * time.Millisecond)
		os.Exit(1)
	case sig := <-quit:
		fmt.Println("\ngot signal", sig)
	}

	shutdown(httpServer)
	shutdown(httpsServer)
	for _, grpcServer := range mx.GRPCServers {
		grpcServer.GracefulStop()
	}
	time.Sleep(100 * time.Millisecond)
}
