package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"

	"github.com/fuku01/test-v2-api/pkg/graph/generated/model"
)

// PostMessage is the resolver for the postMessage field.
func (r *mutationResolver) PostMessage(ctx context.Context, input model.PostMessageInput) (*model.PostMessagePayload, error) {
	return r.Handler.MessageHandler.PostMessage(ctx, input)
}

// Messages is the resolver for the messages field.
func (r *queryResolver) Messages(ctx context.Context) ([]*model.Message, error) {
	return r.Handler.MessageHandler.ListMessages(ctx)
}
