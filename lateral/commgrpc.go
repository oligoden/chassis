package lateral

import (
	context "context"

	"github.com/oligoden/chassis"
	grpc "google.golang.org/grpc"
)

type CommGRPC struct {
	UnimplementedSyncQueueServeṛ
}

func NewCommGRPC() CommGRPC {
	d := CommGRPC{}
	grpcServer := grpc.NewServer()
	RegisterSyncQueueServer(grpcServer, d)
	return d
}

func (d *CommGRPC) Queue(ctx context.Context, in *QueueMessage) (*QueueMessage, error) {
	m := NewGRPCModel(in, d.Store)
	e := NewRecord()
	m.Data(e)

	m.Read()
	if m.Err() != nil {
		return nil, chassis.Mark("reading service provider", m.Err())
	}

	return &QueueMessage{
		QueueID: uint32(e.ID),
		Name:              e.Name,      
	}, nil
}
