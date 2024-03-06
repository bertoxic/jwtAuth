package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)


var connectionString = "mongodb://localhost:27017/"

type DB struct {
    client *mongo.Client
}

func CreateDB()*DB{
     db:=  Connect()
    return db
}

 type mongoClientx *mongo.Client
func Connect() *DB{

    mongoContext, cancel := context.WithTimeout(context.Background(), 20 * time.Second)
    defer cancel()

    mongoClientx, err := mongo.Connect(mongoContext, options.Client().ApplyURI(connectionString))
    if err != nil {
        log.Fatal(err)
        return nil
    }

    err = mongoClientx.Ping(mongoContext, readpref.Primary())
    if err != nil {
        log.Fatal()
        return nil
    }
    log.Println("Connected to Database")
   
    return &DB{
       client:  mongoClientx,
    }
}

func DBinstance() *mongo.Client {
    // err := godotenv.Load(".env")
    // if err != nil {
    //  log.Fatal("error loading .env file")
    // }
     MongoDb := "mongodb://localhost:27017/"
 
     client , err :=mongo.NewClient(options.Client().ApplyURI(MongoDb))
     if err != nil {
         log.Fatal(err)
     }
     ctx,  cancel:= context.WithTimeout(context.Background(), 14* time.Second)
     defer cancel()
     err = client.Connect(ctx)
     if err != nil {
         log.Fatal(err)
     }
 
     fmt.Println("connected to mongodb!")
 
     return client
    }
    
    var Client *mongo.Client = DBinstance()
 
    func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection{
     var Collection *mongo.Collection = client.Database("jwtAuth").Collection(collectionName)
 
     return Collection
    }