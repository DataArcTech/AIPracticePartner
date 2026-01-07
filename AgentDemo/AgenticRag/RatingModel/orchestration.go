package RatingModel

import (
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func BuildRatingModel(ctx context.Context) (r compose.Runnable[map[string]any, *schema.Message], err error) {
	const (
		RatingChatTemplate = "RatingChatTemplate"
		RatingLambda       = "RatingLambda"
	)
	g := compose.NewGraph[map[string]any, *schema.Message]()
	ratingChatTemplateKeyOfChatTemplate, err := newChatTemplate(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddChatTemplateNode(RatingChatTemplate, ratingChatTemplateKeyOfChatTemplate)
	ratingLambdaKeyOfLambda, err := newLambda(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddLambdaNode(RatingLambda, ratingLambdaKeyOfLambda)
	_ = g.AddEdge(compose.START, RatingChatTemplate)
	_ = g.AddEdge(RatingLambda, compose.END)
	_ = g.AddEdge(RatingChatTemplate, RatingLambda)
	r, err = g.Compile(ctx, compose.WithGraphName("RatingModel"))
	if err != nil {
		return nil, err
	}
	return r, err
}
