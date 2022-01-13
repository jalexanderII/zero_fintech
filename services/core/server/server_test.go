package server

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/jalexanderII/zero_fintech/services/core/config"
	"github.com/jalexanderII/zero_fintech/services/core/config/middleware"
	"github.com/jalexanderII/zero_fintech/services/core/database"
	_ "github.com/joho/godotenv/autoload"
)

var L = hclog.Default()

func GenServer() (*CoreServer, context.Context) {
	// ctx, cancel := config.NewDBContext(5 * time.Second)
	// defer cancel()

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
}
