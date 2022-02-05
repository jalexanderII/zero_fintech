package client

import (
	"log"

	"github.com/jalexanderII/zero_fintech/gen/Go/payments"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SetUpPlaidClient() payments.PlaidClient {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	plaidConn, err := grpc.Dial("localhost:9095", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	return payments.NewPlaidClient(plaidConn)
}
