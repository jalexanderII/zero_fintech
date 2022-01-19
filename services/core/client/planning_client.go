package client

import (
	"log"

	"github.com/jalexanderII/zero_fintech/gen/Go/planning"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SetUpPlanningClient() planning.PlanningClient {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithBlock())

	planningConn, err := grpc.Dial("localhost:50051", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	return planning.NewPlanningClient(planningConn)
}
