package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

const Performance = 100

var (
	// Used to create a singleton object of MongoDB client.
	// Initialized and exposed through  GetMongoClient()
	clientInstance *mongo.Client
	// Used during creation of singleton client object in GetMongoClient()
	clientInstanceError error
	// Used to execute client creation procedure only once
	mongoOnce sync.Once
)

// GetEnv func to get env values
func GetEnv(key string) string {
	_, b, _, _ := runtime.Caller(0)
	// Root folder of this project
	Root := filepath.Join(filepath.Dir(b), "../")
	environmentPath := filepath.Join(Root, ".env")
	err := godotenv.Load(environmentPath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv(key)
}

// InitiateMongoClient connects to MongoDB URI and binds a database
func InitiateMongoClient() (mongo.Database, error) {
	// Perform connection creation operation only once.
	mongoOnce.Do(func() {
		// Set client options
		clientOptions := options.Client().ApplyURI(GetEnv("MONGOURI"))
		ctx, cancel := NewDBContext(10 * time.Second)
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

	return *clientInstance.Database(GetEnv("CORE_DB_NAME")), clientInstanceError
}

// ListGRPCResources is a helper function that lists all URLs that are registered on gRPC server.
// This makes it easy to register all the relevant routes in your HTTP router of choice.
func ListGRPCResources(server *grpc.Server) []string {
	var ret []string
	for serviceName, serviceInfo := range server.GetServiceInfo() {
		for _, methodInfo := range serviceInfo.Methods {
			fullResource := fmt.Sprintf("/%s/%s", serviceName, methodInfo.Name)
			ret = append(ret, fullResource)
		}
	}
	return ret
}

// NewDBContext returns a new Context according to app performance
func NewDBContext(d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), d*Performance/100)
}
