package server

import (
	"context"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/jalexanderII/zero_fintech/gen/Go/common"
	"github.com/jalexanderII/zero_fintech/gen/Go/core"
	"github.com/jalexanderII/zero_fintech/gen/Go/planning"
	"github.com/jalexanderII/zero_fintech/services/auth/config/middleware"
	"github.com/jalexanderII/zero_fintech/utils"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/jalexanderII/zero_fintech/services/core/database"

	_ "github.com/joho/godotenv/autoload"
)

var L = logrus.New()

func MockPlanningClient() planning.PlanningClient {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	planningConn, err := grpc.Dial("localhost:9092", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	return planning.NewPlanningClient(planningConn)
}

func GenServer() (*CoreServer, context.Context) {
	jwtManager := middleware.NewJWTManager(utils.GetEnv("JWTSecret"), 15*time.Minute)
	DB, err := database.InitiateMongoClient()
	if err != nil {
		log.Fatal("MongoDB error: ", err)
	}
	coreCollection := *DB.Collection(utils.GetEnv("CORE_COLLECTION"))
	accountCollection := *DB.Collection(utils.GetEnv("ACCOUNT_COLLECTION"))
	transactionCollection := *DB.Collection(utils.GetEnv("TRANSACTION_COLLECTION"))
	userCollection := *DB.Collection(utils.GetEnv("USER_COLLECTION"))

	server := NewCoreServer(coreCollection, accountCollection, transactionCollection, userCollection, jwtManager, MockPlanningClient(), L)
	return server, context.TODO()
}

func TestCoreServer_GetPaymentPlan(t *testing.T) {
	server, ctx := GenServer()

	user_selections := []*core.AccountInfo{
		{
			TransactionIds: []string{"61dfa20adebb9d4fb62b9703"},
			AccountId:      "61df9b621d2c2b15a6e53ec9",
			Amount:         325,
		},
		{
			TransactionIds: []string{},
			AccountId:      "61df9b621d2c2b15a6e53ec9",
			Amount:         10000,
		},
	}

	paymentPlans, err := server.GetPaymentPlan(ctx, &core.GetPaymentPlanRequest{AccountInfo: user_selections, UserId: "61ce2e19014fbb650838306c"})
	if err != nil {
		t.Errorf("1: Error creating new paymentTask: %v", err)
	}
	if paymentPlans.PaymentPlans[0].Timeline != 12 {
		t.Errorf("2: Error creating payment plan from DB response: %v", paymentPlans)
	}
	if len(paymentPlans.PaymentPlans[0].PaymentTaskId) != 2 {
		t.Errorf("3: Error payment plan doesnt have enough tasks ids: %v", paymentPlans.PaymentPlans[0].PaymentTaskId)
	}
}

// CustomGenerator to handle enum fields
func CustomGenerator() {
	_ = faker.AddProvider("preferred_plan_type", func(v reflect.Value) (interface{}, error) {
		return common.PlanType_PLAN_TYPE_OPTIM_CREDIT_SCORE, nil
	})
	_ = faker.AddProvider("preferred_payment_freq", func(v reflect.Value) (interface{}, error) {
		return common.PaymentFrequency_PAYMENT_FREQUENCY_MONTHLY, nil
	})

	_ = faker.AddProvider("penalty_reason", func(v reflect.Value) (interface{}, error) {
		return core.PenaltyAPR_PENALTY_REASON_LATE_PAYMENT, nil
	})
}

func GenFakePaymentTask() (*common.PaymentTask, error) {
	CustomGenerator()
	var fake common.PaymentTask
	err := faker.FakeData(&fake)
	if err != nil {
		L.Error("[Error] Could not fake this object", "error", err)
		return nil, err
	}
	return &fake, nil
}
