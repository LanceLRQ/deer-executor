package server

import (
	"context"
	"github.com/LanceLRQ/deer-executor/v2/server/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"testing"
	"time"
)

func TestRpcConnection(t *testing.T) {
	conn, err := grpc.Dial("127.0.0.1:7150", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := rpc.NewJudgementServiceClient(conn)
	timeoutContext, _ := context.WithTimeout(context.Background(), 10*time.Second)
	resp, err := client.Ping(timeoutContext, &rpc.PingRequest{})
	if err != nil {
		t.Fatalf("cannot ping: %v", err)
	}
	log.Println("ok, server timestamp:", resp.GetTime())
}
