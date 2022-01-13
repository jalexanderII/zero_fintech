package server

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/jalexanderII/zero_fintech/services/core/config"
	"github.com/jalexanderII/zero_fintech/services/core/config/middleware"
	"github.com/jalexanderII/zero_fintech/services/core/database"
	"github.com/jalexanderII/zero_fintech/services/core/gen/core"
	_ "github.com/joho/godotenv/autoload"
)

var L = hclog.Default()

func GenServer() (*CoreServer, context.Context) {
	jwtManager := middleware.NewJWTManager(config.GetEnv("JWTSecret"), 15*time.Minute)
	DB := database.InitiateMongoClient()
	coreCollection := *DB.Collection(config.GetEnv("CORE_COLLECTION"))
	accountCollection := *DB.Collection(config.GetEnv("ACCOUNT_COLLECTION"))
	transactionCollection := *DB.Collection(config.GetEnv("TRANSACTION_COLLECTION"))
	userCollection := *DB.Collection(config.GetEnv("USER_COLLECTION"))

	server := NewCoreServer(coreCollection, accountCollection, transactionCollection, userCollection, jwtManager, L)
	return server, context.TODO()
}

func TestCoreServer_GetPaymentPlan(t *testing.T) {
	server, ctx := GenServer()

	ids := []string{"61dfa8296c734067e6726761", "61dfa8a087ac88bb1559099c"}

	paymentPlans, err := server.GetPaymentPlan(ctx, &core.GetPaymentPlanRequest{PaymentTasksIds: ids})
	if err != nil {
		t.Errorf("1: Error creating new paymentTask: %v", err)
	}
	if paymentPlans.PaymentPlans[0].Timeline != 12 {
		t.Errorf("2: Error creating payment plan from DB response: %v", paymentPlans)
	}
}
