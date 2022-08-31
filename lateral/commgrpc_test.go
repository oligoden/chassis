package lateral_test

import (
	context "context"
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/oligoden/chassis"
	"github.com/oligoden/chassis/lateral"
	"github.com/stretchr/testify/assert"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func TestCommClient(t *testing.T) {
	comm := lateral.NewCommGRPC()
	comm.Open("")

	assert.Equal(t, "a", testStore.Data)
}

type mockServiceProviderServer struct {
	lateral.UnimplementedReadServiceProviderServer
}

func (*mockServiceProviderServer) ServiceProvider(ctx context.Context, req *lateral.ServiceProviderMessage) (*lateral.ServiceProviderMessage, error) {
	if req.GetServiceProviderUC() == "acd" {
		return &lateral.ServiceProviderMessage{
			ServiceProviderUC: "acd",
			ServiceProviderID: 1,
			Name:              "ABC Painters",
		}, nil
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "can't find %v", req.GetServiceProviderUC())
	}
}

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	lateral.RegisterReadServiceProviderServer(server, &mockServiceProviderServer{})
	go func() {
		if err := server.Serve(listener); err != nil {
			fmt.Println(chassis.Mark("opening mock GRPC server", err))
			os.Exit(1)
		}
	}()
	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}
