package livemonitor

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/networkservicemesh/networkservicemesh/pkg/livemonitor/api"

	"google.golang.org/grpc"
)

// EventStream is an unified interface for blocking receiver
type EventStream interface {
	Recv() (interface{}, error)
}

// EventStreamConstructor is a type for EventStream constructor
type EventStreamConstructor func(ctx context.Context, cc *grpc.ClientConn) (EventStream, error)

// Client is an unified interface for GRPC monitoring API client
type Client interface {
	ErrorChannel() <-chan error

	Context() context.Context
	Close()
}

type client struct {
	errorCh <-chan error

	ctx    context.Context
	cancel context.CancelFunc
}

// NewClient creates a new Client on given GRPC connection
func NewClient(cc *grpc.ClientConn) (Client, error) {
	ctx, cancel := context.WithCancel(context.Background())

	monitorClient := api.NewLivenessMonitorClient(cc)
	stream, err := monitorClient.MonitorLiveness(ctx, &empty.Empty{})

	if err != nil {
		cancel()
		return nil, err
	}

	errorChannel := make(chan error, 1)
	go func() {
		defer close(errorChannel)

		for {
			_, err := stream.Recv()
			if err != nil {
				errorChannel <- err
				return
			}
		}
	}()

	return &client{
		errorCh: errorChannel,
		ctx:     ctx,
		cancel:  cancel,
	}, nil
}

// ErrorChannel returns client errorChannel
func (c *client) ErrorChannel() <-chan error {
	return c.errorCh
}

// Context returns client context
func (c *client) Context() context.Context {
	return c.ctx
}

// Close cancels client context
func (c *client) Close() {
	c.cancel()
}

