// package main

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// // This is a user defined method to close resources.
// // This method closes mongoDB connection and cancel context.
// func close(client *mongo.Client, ctx context.Context,
// 	cancel context.CancelFunc) {

// 	defer cancel()
// 	defer func() {
// 		if err := client.Disconnect(ctx); err != nil {
// 			panic(err)
// 		}
// 	}()
// }

// // This is a user defined method that returns
// // mongo.Client, context.Context,
// // context.CancelFunc and error.
// // mongo.Client will be used for further
// // database operation. context.Context will be
// // used set deadlines for process.
// // context.CancelFunc will be used to cancel
// // context and resource associated with it.
// func connect(uri string) (*mongo.Client, context.Context,
// 	context.CancelFunc, error) {

// 	ctx, cancel := context.WithTimeout(context.Background(),
// 		30*time.Second)
// 	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
// 	return client, ctx, cancel, err
// }

// // deleteOne is a user defined function that delete,
// // a single document from the collection.
// // Returns DeleteResult and an  error if any.
// func deleteOne(client *mongo.Client, ctx context.Context,
// 	dataBase, col string, query interface{}) (result *mongo.DeleteResult, err error) {

// 	// select document and collection
// 	collection := client.Database(dataBase).Collection(col)

// 	// query is used to match a document  from the collection.
// 	result, err = collection.DeleteOne(ctx, query)
// 	return
// }

// // deleteMany is a user defined function that delete,
// // multiple documents from the collection.
// // Returns DeleteResult and an  error if any.
// func deleteMany(client *mongo.Client, ctx context.Context,
// 	dataBase, col string, query interface{}) (result *mongo.DeleteResult, err error) {

// 	// select document and collection
// 	collection := client.Database(dataBase).Collection(col)

// 	// query is used to match  documents  from the collection.
// 	result, err = collection.DeleteMany(ctx, query)
// 	return
// }

// func main() {

// 	// get Client, Context, CancelFunc and err from connect method.
// 	client, ctx, cancel, err := connect("mongodb://localhost:27017")

// 	if err != nil {
// 		panic(err)
// 	}

// 	//  free resource when main function is returned
// 	defer close(client, ctx, cancel)

// 	// This query delete document when the maths
// 	query := bson.D{{"commentid", "60111106"}}

// 	// Returns result of deletion and error
// 	result, err := deleteOne(client, ctx, "local", "comments", query)

// 	// print the count of affected documents
// 	fmt.Println("No.of rows affected by DeleteOne()")
// 	fmt.Println(result.DeletedCount)

// 	// // This query deletes  documents that has
// 	// // science field greater that 0
// 	// query = bson.D{
// 	//     {"science", bson.D{{"$gt", 0}}},
// 	// }

// 	// // Returns result of deletion and error
// 	// result, err = deleteMany(client, ctx, "gfg", "marks", query)

// 	// // print the count of affected documents
// 	// fmt.Println("No.of rows affected by DeleteMany()")
// 	// fmt.Println(result.DeletedCount)
// }
