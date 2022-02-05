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
	var uId = "61df93c0ac601d1be8e64613"

	userSelections := []*core.AccountInfo{
		{
			TransactionIds: []string{"61dfa20adebb9d4fb62b9703"}, // Pay Equinox charge
			AccountId:      "61df9b621d2c2b15a6e53ec9",           // Amex
			Amount:         325,
		},
		{
			TransactionIds: []string{},                 // Pay full account
			AccountId:      "61df9af7f18b94fc44d09fb9", // Chase
			Amount:         9000,
		},
	}
	metaData := &common.MetaData{
		PreferredPlanType:         common.PlanType_PLAN_TYPE_OPTIM_CREDIT_SCORE,
		PreferredTimelineInMonths: 3,
		PreferredPaymentFreq:      common.PaymentFrequency_PAYMENT_FREQUENCY_MONTHLY,
	}
	paymentPlans, err := server.GetPaymentPlan(ctx,
		&core.GetPaymentPlanRequest{
			AccountInfo: userSelections,
			UserId:      uId,
			MetaData:    metaData,
		})
	if err != nil {
		t.Errorf("1: Error creating new paymentTask: %v", err)
	}
	if len(paymentPlans.PaymentPlans) != 1 {
		t.Errorf("2: Error from Planning, should have only 1 payment plan, but have: %v", len(paymentPlans.PaymentPlans))
	}
	if paymentPlans.PaymentPlans[0].GetUserId() != uId {
		t.Errorf("3: Error from Planning, wrong user_id, expected %v, got %v", uId, paymentPlans.PaymentPlans[0].GetUserId())
	}
	expectedTotal := userSelections[0].Amount + userSelections[1].Amount
	expectedAmount := expectedTotal / metaData.PreferredTimelineInMonths
	if int(paymentPlans.PaymentPlans[0].GetAmountPerPayment()) != int(expectedAmount) {
		t.Errorf("3: Error from Planning, amount per payment is off, expected %v, got %v", expectedAmount, paymentPlans.PaymentPlans[0].GetAmountPerPayment())
	}
	var total = 0.0
	for _, action := range paymentPlans.PaymentPlans[0].GetPaymentAction() {
		total += action.GetAmount()
	}
	if total != expectedTotal {
		t.Errorf("4: Error from Planning, payment action total does not match, expected %v, got %v", expectedTotal, total)
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
