package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.44

import (
	"context"
	"fmt"

	"github.com/bertoxic/jwtAuth/database"
	"github.com/bertoxic/jwtAuth/graph/model"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	return db.CreateUser(input)
}

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input model.CreatePost) (*model.Post, error) {
	return db.CreatePost(input)
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, input model.CreateCommentInput) (*model.Comment, error) {
	return db.CreateComment(input)
}

// FollowUser is the resolver for the followUser field.
func (r *mutationResolver) FollowUser(ctx context.Context, input model.FollowUserInput) (*model.Follow, error) {
	panic(fmt.Errorf("not implemented: FollowUser - followUser"))
}

// UnfollowUser is the resolver for the unfollowUser field.
func (r *mutationResolver) UnfollowUser(ctx context.Context, input *model.FollowUserInput) (bool, error) {
	panic(fmt.Errorf("not implemented: UnfollowUser - unfollowUser"))
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*model.User, error) {
	return db.Login(email, password)
}

// Signup is the resolver for the signup field.
func (r *mutationResolver) Signup(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	return db.Signup(&input)
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	return db.User(id)
}

// AllUsers is the resolver for the allUsers field.
func (r *queryResolver) AllUsers(ctx context.Context) ([]*model.User, error) {
	panic(fmt.Errorf("not implemented: AllUsers - allUsers"))
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context, userID string) ([]*model.Post, error) {
	return db.Posts(userID)
}

// IsFollowing is the resolver for the isFollowing field.
func (r *queryResolver) IsFollowing(ctx context.Context, userID string, targetUserID string) (bool, error) {
	panic(fmt.Errorf("not implemented: IsFollowing - isFollowing"))
}

// Followers is the resolver for the followers field.
func (r *queryResolver) Followers(ctx context.Context, id string) ([]*model.User, error) {
	panic(fmt.Errorf("not implemented: Followers - followers"))
}

// Comment is the resolver for the comment field.
func (r *queryResolver) Comment(ctx context.Context, commetID string) (*model.Comment, error) {
	return db.Comment(commetID)
}

// AllComments is the resolver for the allComments field.
func (r *queryResolver) AllComments(ctx context.Context, id string) ([]*model.Comment, error) {
	return db.AllComments(id)
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, postID string) (*model.Post, error) {
	return db.Post(postID)
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
var db = database.Connect()
