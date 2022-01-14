package database

import (
	"sync"
	"time"

	"github.com/jalexanderII/zero_fintech/services/core/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/* Used to create a singleton object of MongoDB client.
Initialized and exposed through  GetMongoClient().*/
var clientInstance *mongo.Client

// Used during creation of singleton client object in GetMongoClient().
var clientInstanceError error

// Used to execute client creation procedure only once.
var mongoOnce sync.Once

// InitiateMongoClient connects to MongoDB URI and binds a database
func InitiateMongoClient() (mongo.Database, error) {
	// clientOptions := options.Client().ApplyURI(config.GetEnv("MONGOURI"))
	// ctx, cancel := config.NewDBContext(10 * time.Second)
	// defer cancel()
	//
	// client, err := mongo.Connect(ctx, clientOptions)
	// if err != nil {
	// 	log.Fatalf("Error connecting to DB: %v", err)
	// }

	// Perform connection creation operation only once.
	mongoOnce.Do(func() {
		// Set client options
		clientOptions := options.Client().ApplyURI(config.GetEnv("MONGOURI"))
		ctx, cancel := config.NewDBContext(10 * time.Second)
		defer cancel()
		// Connect to MongoDB
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			clientInstanceError = err
		}
		// Check the connection
		err = client.Ping(ctx, nil)
		if err != nil {
			clientInstanceError = err
		}
		clientInstance = client
	})

	return *clientInstance.Database(config.GetEnv("CORE_DB_NAME")), clientInstanceError
}
