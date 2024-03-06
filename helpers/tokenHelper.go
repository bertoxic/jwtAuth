package helpers

import (
	//"context"
	"log"
	"time"


	"github.com/dgrijalva/jwt-go"

	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	// //"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)


type SignedDetails struct {

    Email string
    Name string
    jwt.StandardClaims
}




func GenerateToken(email string, name string)(signedToken, signedRefreshToken string , err error ){

    claims:= &SignedDetails{
        Email: email,
        Name: name,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
        },
    }

    refreshClaims := &SignedDetails{
        Email: email,
        Name: name,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Local().Add(time.Hour*time.Duration(168)).Unix(),
        },
    }

    token , err:= jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("hope")) 
    if err != nil {
        log.Println("eror in generating xggkkk", err)
        return 
    }

    refreshTokens , err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte("hope"))
    if err != nil {
        log.Println("error in generating cvxvvxv", err)
        return
    }

    return token, refreshTokens , nil

    
}

// func UpdateTokenToDB(signedToken, signedRefreshToken, userId string, db database.DB){
//      ctx, cancel:=  context.WithTimeout(context.Background(), 30 *time.Second)
//      defer cancel()
//      var updateObj primitive.D
//      updateObj = append(updateObj, bson.E{"token",signedToken})
//      updateObj = append(updateObj, bson.E{"refresh_token",signedRefreshToken})

//      updated_at, _ := time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))

//      updateObj = append(updateObj, bson.E{"updated_at", updated_at})
//         upsrt:= true
//     userIdx, err:= primitive.ObjectIDFromHex(userId) 
//     if err != nil {
//         log.Println("error occured in updatedtoken function", err)
//         return
//     }
//     filter := bson.M{"_id":userIdx}

//     opt := options.UpdateOptions {
//         Upsert:&upsrt,

//     }
//     _, err = database.OpenCollection(Client,"user").InsertOne(ctx, )
    

// }