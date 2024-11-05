package migrations

import (
	"context"
	_ "embed"

	"github.com/erodriguezg/go-mongodb-migrate/pkg/migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	personsCollection = "persons"
)

//go:embed 003_persons.go
var migration003 string

func init() {

	err := migrate.Register(

		// Embed source code for hashing
		&migration003,

		// Up function
		func(db *mongo.Database) error {

			_, err := db.Collection(personsCollection).InsertMany(
				context.TODO(),
				[]interface{}{
					bson.M{
						"email":       "erodriguez.cl@gmail.com",
						"firstName":   "Eduardo",
						"lastName":    "Rodriguez",
						"profileCode": 1,
						"active":      true},
					bson.M{
						"email":       "pongoway3@gmail.com",
						"firstName":   "Pongo",
						"lastName":    "Way3",
						"profileCode": 1,
						"active":      true},
					bson.M{
						"email":       "pongoway2@gmail.com",
						"firstName":   "Pongo",
						"lastName":    "Way2",
						"profileCode": 2,
						"active":      true},
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

			return db.Collection(personsCollection).Drop(context.TODO())
		})

	if err != nil {
		panic(err)
	}

}
