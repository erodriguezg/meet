package migrations

import (
	"context"
	_ "embed"

	"github.com/erodriguezg/go-mongodb-migrate/pkg/migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	permissionsCollection = "permissions"
)

//go:embed 001_permissions.go
var migration001 string

func init() {

	err := migrate.Register(

		// Embed source code for hashing
		&migration001,

		// Up function
		func(db *mongo.Database) error {

			_, err := db.Collection(permissionsCollection).InsertMany(
				context.TODO(),
				[]interface{}{
					bson.M{"code": int(1), "name": "Manage System"},
					bson.M{"code": int(2), "name": "Edit Own Profile"},
					bson.M{"code": int(3), "name": "Create Room"},
				},
			)

			if err != nil {
				return err
			}

			return nil

		},

		// Down function
		func(db *mongo.Database) error {

			// ----------------------------------------------------------------------------
			// INFO:
			// you must delete the entire collection.
			// Because this is a type collection and must not contain generated use data.
			// ----------------------------------------------------------------------------

			return db.Collection(permissionsCollection).Drop(context.TODO())
		})

	if err != nil {
		panic(err)
	}

}
