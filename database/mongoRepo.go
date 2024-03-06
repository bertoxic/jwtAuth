package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/bertoxic/jwtAuth/graph/model"
	"github.com/bertoxic/jwtAuth/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (db *DB) CreateUser(user model.CreateUserInput) (*model.User, error) {

	userCollection := db.client.Database("jwtAuth").Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	addressBson := bson.M{
		"city":   user.Address.City,
		"state":  user.Address.State,
		"street": user.Address.Street,
		"zip":    user.Address.Zip,
	}

	insertedResult, err := userCollection.InsertOne(ctx, bson.M{
		"name":    user.Name,
		"email":   user.Email,
		"address": addressBson,
	})

	if err != nil {
		return nil, err
	}

	userx := &model.User{
		Name:    user.Name,
		Email:   user.Email,
		Address: (*model.Address)(user.Address),
	}
	insertedID := insertedResult.InsertedID.(primitive.ObjectID).Hex()
	userx.ID = insertedID

	return userx, nil

}

func (db *DB) CreatePost(userPost model.CreatePost) (*model.Post, error) {
	userCollection := db.client.Database("jwtAuth").Collection("post")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	newPost := &model.CreatePost{
		UserID: userPost.UserID,
		//ID: string(primitive.NewObjectID().Hex()),
		Title:   userPost.Title,
		Content: userPost.Content,
	}
	user, err := db.User(userPost.UserID)
	if err != nil {
		return nil, errors.New("unable to get user in createPost method")
	}
	useridx, err := primitive.ObjectIDFromHex(userPost.UserID)
	if err != nil {
		return nil, err
	}
	insertedIResult, err := userCollection.InsertOne(ctx, bson.M{
		"userId":  useridx,
		"title":   userPost.Title,
		"content": userPost.Content,
	})
	if err != nil {
		return nil, err
	}
	insertedID := insertedIResult.InsertedID.(primitive.ObjectID).Hex()

	return &model.Post{
		ID:       insertedID,
		UserID:   newPost.UserID,
		Title:    newPost.Title,
		Content:  newPost.Content,
		Author:   user,
		Comments: []*model.Comment{},
	}, nil
}

func (db *DB) User(id string) (*model.User, error) {

	jobCollects := db.client.Database("jwtAuth").Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}
	defer cancel()
	var user model.User
	err := jobCollects.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}
	return &user, nil
}

func (db *DB) CreateComment(comment model.CreateCommentInput) (*model.Comment, error) {
	commentCollection := db.client.Database("jwtAuth").Collection("comment")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var err error
	commentCreated, err := commentCollection.InsertOne(ctx, bson.M{
		"postid":          comment.PostID,
		"userid":          comment.UserID,
		"content":         comment.Content,
		"parentCommentId": comment.ParentCommentID,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var commentnew = model.Comment{}

	commentList, err := db.AllComments(comment.UserID)
	if err != nil {
		return nil, err
	}

	commentnew.ID = commentCreated.InsertedID.(primitive.ObjectID).Hex()
	commentnew.Content = comment.Content
	commentnew.UserID = comment.UserID
	commentnew.PostID = comment.PostID
	commentnew.Replies = commentList
	log.Println(comment.ParentCommentID, "xxxxxxxxxxxxxxxmm")

	if comment.ParentCommentID != nil {
		commentnew.ParentComment, err = db.Comment(*comment.ParentCommentID)
		log.Println(commentnew.ParentComment, "oooooooooooooooooo")
		if err != nil {
			println("wahala be like bicycle >>>>>>>>>>>>>>>>>>>>", err.Error())
			return nil, err
		}
	}
	return &commentnew, nil
}

func (db *DB) Posts(userID string) ([]*model.Post, error) {

	postCollection := db.client.Database("jwtAuth").Collection("post")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_id, _ := primitive.ObjectIDFromHex(userID)
	filter := bson.M{"_id": _id}
	var postList []*model.Post
	cursor, err := postCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &postList); err != nil {
		return nil, fmt.Errorf("error occured in post method %s", err)
	}

	return postList, nil
}

func (db *DB) Post(postID string) (*model.Post, error) {
	postCollection := db.client.Database("jwtAuth").Collection("post")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	postIDx, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": postIDx}
	var foundPost model.Post
	err = postCollection.FindOne(ctx, filter).Decode(&foundPost)
	if err != nil {
		return nil, err
	}

	comments, err := db.AllComments(postID)
	log.Println(comments, "tttttttttttttttt")
	if err != nil {
		return nil, err
	}
	foundPost.Comments = comments
	return &foundPost, nil

}

func (db *DB) Comment(commentID string) (*model.Comment, error) {
	commentCollection := db.client.Database("jwtAuth").Collection("comment")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	commentIDx, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": commentIDx}

	var founCommentx = model.CreateCommentInput{}
	var founComment = model.Comment{}
	err = commentCollection.FindOne(ctx, filter).Decode(&founCommentx)
	if err != nil {
		log.Printf("eroor in findone foudcommentx: %+v", err.Error())
		return nil, err
	}
	founComment.UserID = founCommentx.UserID
	founComment.Content = founCommentx.Content
	founComment.PostID = founCommentx.PostID
	founComment.ID = commentID

	var author = &model.User{}
	author, _ = db.User(founCommentx.UserID)
	founComment.Author = author

	if founCommentx.ParentCommentID != nil {
		var pComment = &model.CreateCommentInput{}
		parentIDx, err := primitive.ObjectIDFromHex(*founCommentx.ParentCommentID)
		if err != nil {
			log.Printf("eroor in findone parentIDXXXXXXXX: %+v", err.Error())
			return nil, err
		}

		filterpare := bson.M{"_id": parentIDx}
		err = commentCollection.FindOne(ctx, filterpare).Decode(&pComment)
		if err != nil {
			log.Printf("eroor in findone pcomment: %+v", err.Error())
			return nil, err
		}

		var parentcomment = &model.Comment{}
		parentcomment.ID = *founCommentx.ParentCommentID
		parentcomment.Content = pComment.Content
		parentcomment.PostID = pComment.PostID

		founComment.ParentComment = &model.Comment{}

		founComment.ParentComment = parentcomment
		parentcomment.ParentComment = &model.Comment{}
		parentcomment.ParentComment, _ = db.Comment(pComment.PostID)
		log.Printf("passed here-----------4.555: %+v", *founCommentx.ParentCommentID)
		var repliesList = []*model.Comment{}
		repliesList, _ = db.AllComments(founComment.PostID)
		founComment.Replies = repliesList

	}
	//log.Printf("found comment issssss: %+v", &founComment)

	return &founComment, nil

}

func (db *DB) AllComments(userID string) ([]*model.Comment, error) {
    commentCollection := db.client.Database("jwtAuth").Collection("comment") // Use the correct collection name
    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()
	//userIDx, _ := primitive.ObjectIDFromHex(userID)
    filter := bson.M{"userid": userID}

    var commentList []*model.Comment
    cursor, err := commentCollection.Find(ctx, filter)
    if err != nil {
        return nil, err
    }

    if err = cursor.All(ctx, &commentList); err != nil { 
        return nil, fmt.Errorf("error occurred in AllComments method: %w", err)
    }
	log.Println(commentList)
    return commentList, nil
}


func (db *DB) AllUsern(userID string) ([]*model.User, error) {
	user := &model.User{}
	userCollection := db.client.Database("jwtAuth").Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	id, _ := primitive.ObjectIDFromHex(userID)
	filter := bson.M{"_id": id}

	err := userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return []*model.User{}, nil
}

func (db *DB) Signup(createUser *model.CreateUserInput) (*model.User,error){
	userCollection := db.client.Database("jwtAuth").Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
		log.Println(createUser)	
	addressBson := bson.M{
		"city":   createUser.Address.City,
		"state":  createUser.Address.State,
		"street": createUser.Address.Street,
		"zip":    createUser.Address.Zip,
	}

	insertedResult, err := userCollection.InsertOne(ctx, bson.M{
		"name":    createUser.Name,
		"email":   createUser.Email,
		"address": addressBson,
	})
	if err != nil {
	return nil, err
	}
	var insertedID = insertedResult.InsertedID.(primitive.ObjectID).Hex()
	_id, _ := primitive.ObjectIDFromHex(insertedID)
	user := &model.User{}
	filter := bson.M{"_id":_id }
	userCollection.FindOne(ctx, filter).Decode(&user)
	user.ID = insertedID
	token, refreshToken , err := helpers.GenerateToken(user.Email, user.Name)
	if err != nil {
		log.Println("unable to generate token in singnup",err)
		return nil, err
	}

	log.Println("generated tokennnnnn", token)
	db.UpdateTokenToDB(token,refreshToken,insertedID)
	userCollection.FindOne(ctx, filter).Decode(&user)

	return user, nil
}

func (db *DB) UpdateTokenToDB(signedToken, signedRefreshToken, userId string) {
    log.Print("userid is", userId)
    userCollection := db.client.Database("jwtAuth").Collection("user")
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    objectId, err := primitive.ObjectIDFromHex(userId)
	log.Println("primitive id converted in uppdatetoken is", objectId)
    if err != nil {
        log.Println(err,objectId)
        return
    }

    updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

    filter := bson.M{"_id": &objectId}
    update := bson.D{
        {"$set", bson.D{
            {"token", signedToken},
            {"refresh_token", signedRefreshToken},
            {"updated_at", updatedAt},
        }},
    }

    upsert := false
    opt := options.UpdateOptions{
        Upsert: &upsert,
    }

    result, err := userCollection.UpdateOne(ctx, filter, update, &opt)
    if err != nil {
        log.Println(err)
        return
    }

    log.Printf("Updated %v documents", result.ModifiedCount)
}
type UserX struct {
	ID            primitive.ObjectID `bson:"_id"`
	Name     *string            `json:"name"`
	Password      *string            `json:"password"`
	Email      *string            `json:"email"`
	}    

func (db *DB) Login(email, password string) (*model.User, error) {
    userCollection := db.client.Database("jwtAuth").Collection("user")
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    var foundUser model.User
    err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&foundUser)
	filter := bson.M{"email": email}
	var foundUserID UserX
	err = userCollection.FindOne(ctx, filter).Decode(&foundUserID)
    if err != nil {
        log.Println(err)
        return nil, err
    }

    fmt.Println("xxxxxxxxxxxxxxxxxxxxx", foundUser)
    fmt.Println("xxxxxxxxxxxxxxxxxxxxx", foundUserID)
	foundUser.ID = foundUserID.ID.Hex()

    token, refreshToken, err := helpers.GenerateToken(email, foundUser.Name)
    if err != nil {
        log.Println(err)
        return nil, err
    }
	
    db.UpdateTokenToDB(token, refreshToken, foundUser.ID)
    return &foundUser, nil
}