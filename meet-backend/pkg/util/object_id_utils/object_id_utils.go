package objectidutils

import "go.mongodb.org/mongo-driver/bson/primitive"

func ObjectIDFromHexOrNil(hex *string) (*primitive.ObjectID, error) {
	if hex != nil {
		oid, err := primitive.ObjectIDFromHex(*hex)
		if err != nil {
			return nil, err
		}
		return &oid, nil
	} else {
		return nil, nil
	}
}

func HexFromObjectIDOrNil(oid *primitive.ObjectID) *string {
	if oid != nil {
		hexAux := oid.Hex()
		return &hexAux
	} else {
		return nil
	}
}
