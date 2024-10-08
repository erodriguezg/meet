package config

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/erodriguezg/go-mongodb-migrate/pkg/migrate"
	_ "github.com/erodriguezg/meet/migrations"
)

var (
	mongoDB *mongo.Database
)

func configDatabases() {
	mongoDB = configMongoDB()
	executeMongoMigrations()
}

func configMongoDB() *mongo.Database {
	mongoUrl := propUtils.GetProp("MONGODB_URL")
	mongoDatabase := propUtils.GetProp("MONGODB_DATABASE")

	clientOptions := options.Client().ApplyURI(mongoUrl)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	return client.Database(mongoDatabase)
}

func executeMongoMigrations() {

	migrationEnabled := propUtils.GetBoolProp("MONGODB_MIGRATIONS_ENABLED")
	migrationAutoRepair := propUtils.GetBoolProp("MONGODB_MIGRATIONS_AUTOREPAIR")

	migrate.SetDatabase(mongoDB)
	migrate.SetMigrationsEnabled(migrationEnabled)
	migrate.SetMigrationsAutoRepair(migrationAutoRepair)
	if err := migrate.Up(migrate.AllAvailable); err != nil {
		panic(err)
	}
}
