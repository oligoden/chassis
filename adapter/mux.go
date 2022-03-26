package adapter

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
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

	// if mx.URL.Scheme == "https" {
	// 	mx.serveTLS()
	// 	return
	// }
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

	addr := fmt.Sprintf(":%d", port)
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

			// profile.RegisterReadProfileServer(grpcServer, dProfile)
			if err := grpcServer.Serve(lis); err != nil {
				e := chassis.Mark("failed to serve", err)
				grpcServerError <- grpcError{i, e}
				return
			}
		}()
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	log.Println("running http server")
	select {
	case err := <-serverError:
		fmt.Println(chassis.Mark("http server error", err))
		for _, grpcServer := range mx.GRPCServers {
			grpcServer.GracefulStop()
		}
		time.Sleep(100 * time.Millisecond)
		os.Exit(1)
	case gerr := <-grpcServerError:
		fmt.Println(chassis.Mark("grpc server error", gerr.e))
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
