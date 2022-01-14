package server

import (
	"context"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
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
	DB, err := database.InitiateMongoClient()
	if err != nil {
		log.Fatal("MongoDB error: ", err)
	}
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

// CustomGenerator to handle enum fields
func CustomGenerator() {
	_ = faker.AddProvider("preferred_plan_type", func(v reflect.Value) (interface{}, error) {
		return core.PlanType_PLANTYPE_OPTIM_CREDIT_SCORE, nil
	})
	_ = faker.AddProvider("preferred_payment_freq", func(v reflect.Value) (interface{}, error) {
		return core.PaymentFrequency_PAYMENTFREQ_MONTHLY, nil
	})

	_ = faker.AddProvider("penalty_reason", func(v reflect.Value) (interface{}, error) {
		return core.PenaltyAPR_PENALTY_REASON_LATE_PAYMENT, nil
	})
}

func GenFakePaymentTask() (*database.PaymentTask, error) {
	CustomGenerator()
	var fake database.PaymentTask
	err := faker.FakeData(&fake)
	if err != nil {
		L.Error("[Error] Could not fake this object", "error", err)
		return nil, err
	}
	return &fake, nil
}
