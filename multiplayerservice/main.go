package main

import (
	"context"
	"github.com/horahoradev/nexus2/multiplayerservice/internal/grpc"
	log "github.com/sirupsen/logrus"
)

func main() {

	// Start GRPC Server
	err := grpc.NewGRPCServer(context.TODO(), 5555)
	if err != nil {
		log.Errorf("GRPC server returned with err: %s", err)
	}
}
