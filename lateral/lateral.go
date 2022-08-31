// The Lateral package's responsibility is to take storage interaction queries
// and perform it in such a way the all the storage engins (e.g. databases)
// are all synced.
//
// Typical procedure:
// incoming query --> [r|w] -r-> [check cache] -available-> done
//                          -w-> (queue) --> (sync) --> (write)
package lateral

import (
	context "context"

	"github.com/oligoden/chassis"
	grpc "google.golang.org/grpc"
)

type Service struct {
	UnimplementedSyncQueueServeṛ
}

func New() Service {
	d := Service{}
	return d
}

func (srv Service) Query() {
	
}