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
	"github.com/jalexanderII/zero_fintech/services/auth/config/middleware"
	"github.com/jalexanderII/zero_fintech/utils"
	"github.com/sirupsen/logrus"

	"github.com/jalexanderII/zero_fintech/services/core/database"

	_ "github.com/joho/godotenv/autoload"
)

var L = logrus.New()

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

func GenFakePaymentTask() (*core.PaymentTask, error) {
	CustomGenerator()
	var fake core.PaymentTask
	err := faker.FakeData(&fake)
	if err != nil {
		L.Error("[Error] Could not fake this object", "error", err)
		return nil, err
	}
	return &fake, nil
}
