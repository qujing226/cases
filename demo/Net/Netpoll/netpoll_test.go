package Netpoll

import (
	"context"
	"github.com/cloudwego/netpoll"
	"net"
	"testing"
	"time"
)

var eventLoop netpoll.EventLoop

func TestNetpool(t *testing.T) {
	listener1, _ := net.Listen("tcp", "127.0.0.1:8080")

	eventLoop, _ = netpoll.NewEventLoop(
		func(ctx context.Context, connection netpoll.Connection) error {
			return nil
		},
		netpoll.WithOnPrepare(prepare),
		netpoll.WithOnConnect(connect),
		netpoll.WithReadTimeout(time.Second),
	)
	err := eventLoop.Serve(listener1)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = eventLoop.Shutdown(ctx)

	dial, _ := netpoll.DialConnection("tcp", "127.0.0.1:8080", time.Second*17)
	dial, _ := netpoll.DialTCP(ctx, "tcp", et.TCPAddr{}, time.Second*17)
}

func prepare(ctx context.Context, connection netpoll.Connection) error {
	return nil
}

func TestNoCopyAPI(t *testing.T) {
	var conn netpoll.Connection
	var reader, writer = conn.Reader(), conn.Writer()
	buf, _ := reader.Next(1024)
	reader.Release()
	var write_data []byte
	alloc, _ := writer.Malloc(len(write_data))
	copy(alloc, write_data)
	writer.Flush()
}
