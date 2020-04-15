package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var collsToCopy = []string{
	"acronym",
	"chunk",
	"file",
	"group",
	"lines",
	"reset",
	"store",
	"user",
	"view",
}

func copyColls(ctx context.Context, client *mongo.Client, src, dest string) error {
	srcDB := client.Database(src)
	destDB := client.Database(dest)

	for _, coll := range collsToCopy {
		srcColl := srcDB.Collection(coll)
		destColl := destDB.Collection(coll)
		cursor, err := srcColl.Find(ctx, bson.D{})
		if err != nil {
			return err
		}
		var result []interface{}
		if err := cursor.All(ctx, &result); err == nil {
			_, err = destColl.InsertMany(ctx, result)
			if err != nil {
				return err
			}
		} else if err != mongo.ErrNoDocuments {
			return err
		}
	}
	return nil
}
