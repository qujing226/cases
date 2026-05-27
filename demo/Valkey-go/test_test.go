package main

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/valkey-io/valkey-go"
	"testing"
	"time"
)

func TestSimpleTest(t *testing.T) {
	client, err := valkey.NewClient(valkey.ClientOption{
		DisableCache:      true,
		InitAddress:       []string{"192.168.110.68:26279"},
		ForceSingleClient: true,
	},
	)
	if err != nil {
		panic(err)
	}

	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 20000*time.Second)
	defer cancel()
	// SET key val NX
	err = client.Do(ctx, client.B().Set().Key("key").Value("val").Nx().Build()).Error()
	require.NoError(t, err)
	// HGETALL hm
	hm, err := client.Do(ctx, client.B().Hgetall().Key("hm").Build()).AsStrMap()
	require.NoError(t, err)
	for key, val := range hm {
		t.Log(key, ":", val)
	}
}
