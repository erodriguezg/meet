package migrations

import (
	"context"
	_ "embed"

	"github.com/erodriguezg/go-mongodb-migrate/pkg/migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	profilesCollection = "profiles"
)

//go:embed 002_profiles.go
var migration002 string

func init() {

	err := migrate.Register(

		// Embed source code for hashing
		&migration002,

		// Up function
		func(db *mongo.Database) error {

			_, err := db.Collection(profilesCollection).InsertMany(
				context.TODO(),
				[]interface{}{
					bson.M{"code": int(1), "name": "Administrator", "permissionsCodes": []int{1, 3}},
					bson.M{"code": int(2), "name": "Publisher", "permissionsCodes": []int{2, 4, 5}},
					bson.M{"code": int(3), "name": "User", "permissionsCodes": []int{2, 5}},
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

			return db.Collection(profilesCollection).Drop(context.TODO())
		})

	if err != nil {
		panic(err)
	}

}
